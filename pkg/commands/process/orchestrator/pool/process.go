package pool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	gopsutilprocess "github.com/shirou/gopsutil/v3/process"

	"github.com/moonkit02/dearer/pkg/commands/process/orchestrator/work"
	"github.com/moonkit02/dearer/pkg/commands/process/orchestrator/worker"
	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/util/output"
)

var (
	ErrorCrashed     = errors.New("exited unexpectedly")
	ErrorNotSpawned  = errors.New("didn't start within expected time")
	ErrorOutOfMemory = errors.New("exceeded memory limit")
)

type Process struct {
	id            string
	command       *exec.Cmd
	context       context.Context
	cancelContext context.CancelFunc
	errorChannel  chan error
	exitChannel   chan struct{}
	client        *http.Client
	baseURL       string
	memoryUsage   uint64
}

type ProcessOptions struct {
	executable    string
	baseArguments []string
	config        *settings.Config
}

func newProcess(options *ProcessOptions, id string) (*Process, error) {
	port, err := allocatePort()
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("%s spawning on port %d", id, port)

	arguments := append(
		options.baseArguments,
		"--parent-process-id", strconv.Itoa(syscall.Getpid()),
		"--port", strconv.Itoa(port),
		"--worker-id", id,
	)
	command := exec.Command(options.executable, arguments...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	context, cancelContext := context.WithCancel(context.Background())

	process := &Process{
		id:            id,
		command:       command,
		context:       context,
		cancelContext: cancelContext,
		errorChannel:  make(chan error, 1),
		exitChannel:   make(chan struct{}),
		client:        &http.Client{Timeout: 0},
		baseURL:       fmt.Sprintf("http://localhost:%d", port),
	}

	if err := process.start(options.config); err != nil {
		process.Close()
		return nil, fmt.Errorf("failed to start process: %w", err)
	}

	return process, nil
}

func (process *Process) start(config *settings.Config) error {
	if err := process.command.Start(); err != nil {
		close(process.exitChannel)
		return err
	}

	go process.monitorCommand()
	go process.monitorMemory()

	if err := process.initialize(config); err != nil {
		var result = strings.Split(err.Error(), "failed to create detector customDetector:")
		if len(result) > 1 {
			// custom detector issue ; assume custom rule parse issue
			var ruleName = strings.TrimSpace(strings.Split(result[1], ":")[0])
			log.Debug().Msgf(err.Error())
			output.Fatal(fmt.Sprintf("could not parse rule %s. Is this a custom rule? See documentation on rule patterns and format https://docs.bearer.com/guides/custom-rule/", ruleName))
		} else {
			output.Fatal(fmt.Sprintf("failed to start bearer, error with your configuration %s", err))
		}

		return err
	}

	return nil
}

func (process *Process) monitorCommand() {
	go func() {
		select {
		case <-process.context.Done():
			log.Debug().Msgf("shutting down %s", process.id)

			if err := process.command.Process.Signal(os.Interrupt); err != nil {
				log.Debug().Msgf("killing %s due to error sending interrupt: %s", process.id, err)
				process.kill()
				return
			}

			timeout := time.NewTimer(settings.TimeoutWorkerShutdown)
			select {
			case <-timeout.C:
				log.Debug().Msgf("killing %s after timeout", process.id)
				process.kill()
			case <-process.exitChannel:
				log.Debug().Msgf("%s stopped", process.id)
			}
		case <-process.exitChannel:
			process.errorChannel <- ErrorCrashed
			return
		}
	}()

	process.command.Wait() //nolint:errcheck
	close(process.exitChannel)
}

func (process *Process) kill() {
	if err := process.command.Process.Kill(); err != nil {
		log.Debug().Msgf("%s failed killing process: %s", process.id, err)
	}
}

func (process *Process) monitorMemory() {
	tick := time.NewTicker(500 * time.Millisecond)
	monitor, err := gopsutilprocess.NewProcessWithContext(process.context, int32(process.command.Process.Pid))
	if err != nil {
		log.Debug().Msgf("failed to start memory monitor: %s", err)
		return
	}

	for {
		select {
		case <-process.context.Done():
			log.Debug().Msgf("%s memory monitor shutting down", process.id)
			return
		case <-tick.C:
			stats, err := monitor.MemoryInfo()
			if err != nil {
				log.Debug().Msgf("failed to get memory usage %s", err)
				continue
			}

			if stats.RSS > settings.MemoryMaximum {
				process.memoryUsage = stats.RSS
				process.errorChannel <- ErrorOutOfMemory
				return
			}

			if stats.RSS > settings.MemorySoftMaximum {
				process.reduceMemoryUsage()
			}
		}
	}
}

func (process *Process) reduceMemoryUsage() {
	request, err := process.buildRequest(work.RouteReduceMemory, nil)
	if err != nil {
		log.Debug().Msgf("%s failed to build memory reduction request %s", process.id, err)
	}

	response, err := process.client.Do(request)
	if err != nil {
		log.Debug().Msgf("%s failed to request memory reduction: %s", process.id, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Debug().Msgf("%s memory reduction not ok response", process.id)
	}
}

func (process *Process) initialize(config *settings.Config) error {
	log.Debug().Msgf("%s initializing", process.id)
	start := time.Now()
	killTime := time.Now().Add(config.Worker.TimeoutWorkerOnline)

	marshalledConfig, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("couldn't marshal config %w", err)
	}

	request, err := process.buildRequest(work.RouteInitialize, bytes.NewBuffer(marshalledConfig))
	if err != nil {
		return fmt.Errorf("%s failed to build initialization request %s", process.id, err)
	}

	for {
		if process.context.Err() != nil {
			return process.context.Err()
		}
		if time.Now().After(killTime) {
			return ErrorNotSpawned
		}

		response, err := process.client.Do(request)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		online, err := process.initializationResponse(response)
		if !online {
			continue
		}

		log.Debug().Msgf("%s spawned after %.2f seconds", process.id, time.Since(start).Seconds())
		if err != nil {
			return err
		}

		log.Debug().Msgf("%s is online", process.id)
		return nil
	}
}

func (process *Process) Scan(scanRequest work.ProcessRequest) (*work.ProcessResponse, error) {
	scanComplete := make(chan *work.ProcessResponse)

	go func() {
		taskBytes, err := json.Marshal(scanRequest)
		if err != nil {
			log.Debug().Msgf("%s failed to marshall task: %s", process.id, err)
			return
		}

		request, err := process.buildRequest(work.RouteProcess, bytes.NewBuffer(taskBytes))
		if err != nil {
			log.Debug().Msgf("%s failed to build scan request: %s", process.id, err)
			return
		}

		response, err := process.client.Do(request)
		if err != nil {
			log.Debug().Msgf("%s failed to scan: %s", process.id, err)
			return
		}

		defer response.Body.Close()

		var scanResponse work.ProcessResponse
		if err := json.NewDecoder(response.Body).Decode(&scanResponse); err != nil {
			log.Debug().Msgf("%s failed to decode scan: %s", process.id, err)
		}

		scanComplete <- &scanResponse
	}()

	timeout := time.NewTimer(scanRequest.File.Timeout + settings.TimeoutWorkerFileGrace)
	select {
	case response := <-scanComplete:
		return response, nil
	case err := <-process.errorChannel:
		process.Close()
		return nil, err
	case <-timeout.C:
		process.Close()
		return nil, worker.ErrorTimeoutReached
	}
}

func (process *Process) buildRequest(route string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(
		process.context,
		http.MethodPost,
		process.baseURL+route,
		body,
	)
}

func (process *Process) initializationResponse(response *http.Response) (bool, error) {
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, nil
	}

	var result work.InitializeResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return true, fmt.Errorf("error decoding status response: %w", err)
	}

	if result.Error != "" {
		return true, errors.New(result.Error)
	}

	return true, nil
}

func (process *Process) Close() {
	process.cancelContext()
	<-process.exitChannel
}

func allocatePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to resolve localhost %w", err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to resolve address %w", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

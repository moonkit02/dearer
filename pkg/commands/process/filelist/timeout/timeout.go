package timeout

import (
	"io/fs"
	"time"

	"github.com/moonkit02/dearer/pkg/commands/process/settings"
)

func Assign(fileInfo fs.FileInfo, config settings.Config) time.Duration {
	var timeout time.Duration

	timeout = config.Worker.TimeoutFileMinimum

	timeoutFileSize := time.Duration(fileInfo.Size() / int64(config.Worker.TimeoutFileBytesPerSecond) * int64(time.Second))
	if timeoutFileSize > timeout {
		if timeoutFileSize > config.Worker.TimeoutFileMaximum {
			timeout = config.Worker.TimeoutFileMaximum
		} else {
			timeout = timeoutFileSize
		}
	}

	return timeout
}

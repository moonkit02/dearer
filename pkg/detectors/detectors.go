package detectors

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime/debug"
	"slices"

	"github.com/rs/zerolog/log"

	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/detectors/beego"
	"github.com/moonkit02/dearer/pkg/detectors/csharp"
	"github.com/moonkit02/dearer/pkg/detectors/custom"
	"github.com/moonkit02/dearer/pkg/detectors/dependencies"
	"github.com/moonkit02/dearer/pkg/detectors/django"
	"github.com/moonkit02/dearer/pkg/detectors/dotnet"
	"github.com/moonkit02/dearer/pkg/detectors/envfile"
	"github.com/moonkit02/dearer/pkg/detectors/gitleaks"
	"github.com/moonkit02/dearer/pkg/detectors/golang"
	"github.com/moonkit02/dearer/pkg/detectors/graphql"
	"github.com/moonkit02/dearer/pkg/detectors/html"
	"github.com/moonkit02/dearer/pkg/detectors/ipynb"
	"github.com/moonkit02/dearer/pkg/detectors/java"
	"github.com/moonkit02/dearer/pkg/detectors/javascript"
	"github.com/moonkit02/dearer/pkg/detectors/openapi"
	"github.com/moonkit02/dearer/pkg/detectors/php"
	"github.com/moonkit02/dearer/pkg/detectors/proto"
	"github.com/moonkit02/dearer/pkg/detectors/python"
	"github.com/moonkit02/dearer/pkg/detectors/rails"
	"github.com/moonkit02/dearer/pkg/detectors/ruby"
	"github.com/moonkit02/dearer/pkg/detectors/simple"
	"github.com/moonkit02/dearer/pkg/detectors/spring"
	"github.com/moonkit02/dearer/pkg/detectors/sql"
	"github.com/moonkit02/dearer/pkg/detectors/symfony"
	"github.com/moonkit02/dearer/pkg/detectors/tsx"
	"github.com/moonkit02/dearer/pkg/detectors/typescript"
	"github.com/moonkit02/dearer/pkg/detectors/yamlconfig"
	"github.com/moonkit02/dearer/pkg/scanner"
	"github.com/moonkit02/dearer/pkg/scanner/stats"
	"github.com/moonkit02/dearer/pkg/util/file"

	"github.com/moonkit02/dearer/pkg/detectors/types"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"

	reporttypes "github.com/moonkit02/dearer/pkg/report"
	reportdetectors "github.com/moonkit02/dearer/pkg/report/detectors"
)

type InitializedDetector struct {
	Type reportdetectors.Type
	types.Detector
}

type activeDetector struct {
	*file.Path
	reporttypes.Report
}

var customDetector = InitializedDetector{reportdetectors.DetectorCustom, custom.New(&nodeid.UUIDGenerator{})}

func SetupLegacyDetector(config map[string]*settings.Rule) error {
	detector := customDetector.Detector.(*custom.Detector)

	return detector.CompileRules(config)
}

func Registrations(scanners []string) []InitializedDetector {
	// The order of these is important, the first one to claim a file will win
	detectors := []InitializedDetector{}

	if slices.Contains(scanners, "secrets") {
		// Enable GitLeaks
		detectors = append(
			detectors,
			InitializedDetector{
				reportdetectors.DetectorGitleaks, gitleaks.New(&nodeid.UUIDGenerator{}),
			},
		)
	}

	if slices.Contains(scanners, "sast") {
		detectors = append(
			detectors,
			[]InitializedDetector{
				{reportdetectors.DetectorCustom, customDetector},

				{reportdetectors.DetectorDependencies, dependencies.New()},

				{reportdetectors.DetectorBeego, beego.New()},
				{reportdetectors.DetectorGo, golang.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorDjango, django.New()},
				{reportdetectors.DetectorPython, python.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorDotnet, dotnet.New()},
				{reportdetectors.DetectorCSharp, csharp.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorOpenAPI, openapi.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorEnvFile, envfile.New()},

				{reportdetectors.DetectorJavascript, javascript.New(&nodeid.UUIDGenerator{})},
				{reportdetectors.DetectorTsx, tsx.New(&nodeid.UUIDGenerator{})},
				{reportdetectors.DetectorTypescript, typescript.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorRails, rails.New(&nodeid.UUIDGenerator{})},
				{reportdetectors.DetectorRuby, ruby.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorSpring, spring.New()},
				{reportdetectors.DetectorJava, java.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorSymfony, symfony.New()},
				{reportdetectors.DetectorPHP, php.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorYamlConfig, yamlconfig.New()},

				{reportdetectors.DetectorSQL, sql.New(&nodeid.UUIDGenerator{})},
				{reportdetectors.DetectorProto, proto.New(&nodeid.UUIDGenerator{})},
				{reportdetectors.DetectorGraphQL, graphql.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorHTML, html.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorIPYNB, ipynb.New(&nodeid.UUIDGenerator{})},

				{reportdetectors.DetectorSimple, simple.New()},
			}...,
		)
	}

	return detectors
}

func Extract(
	ctx context.Context,
	path string,
	filename string,
	report reporttypes.Report,
	fileStats *stats.FileStats,
	enabledScanners []string,
	sastScanner *scanner.Scanner,
	skipTest bool,
) error {
	return ExtractWithDetectors(
		ctx,
		path,
		filename,
		report,
		fileStats,
		Registrations(enabledScanners),
		sastScanner,
		skipTest,
	)
}

func ExtractWithDetectors(
	ctx context.Context,
	rootDir string,
	filename string,
	report reporttypes.Report,
	fileStats *stats.FileStats,
	allDetectors []InitializedDetector,
	sastScanner *scanner.Scanner,
	skipTest bool,
) error {

	activeDetectors := make(map[InitializedDetector]activeDetector)

	if err := file.IterateFilesList(
		rootDir,
		[]string{filename},
		skipTest,
		func(dir *file.Path) (bool, error) {
			for _, detector := range allDetectors {
				active, isActive := activeDetectors[detector]
				if isActive && !isParentedBy(active.AbsolutePath, dir.AbsolutePath) {
					delete(activeDetectors, detector)
					isActive = false
				}

				if !isActive {
					activate, err := detector.Detector.AcceptDir(dir)
					if err != nil {
						report.AddError(dir.RelativePath, fmt.Errorf("accept dir failed for detector %s: %s", detector.Type, err))
						continue
					}

					if activate {
						activeDetectors[detector] = activeDetector{Path: dir, Report: report}
					}
				}
			}

			return true, nil
		},
		func(file *file.FileInfo) error {
			recovery := func() {
				if r := recover(); r != nil {
					log.Printf("file %s -> error recovered %s", file.AbsolutePath, r)
					log.Print(string(debug.Stack()))
					report.AddError(file.Path.RelativePath, fmt.Errorf("skipping file: due to panic %s", r))
				}
			}
			defer recovery()

			if err := sastScanner.Scan(ctx, report, fileStats, file); err != nil {
				log.Debug().Msgf("failed to scan file %s: %s", file.RelativePath, err)
				report.AddError(file.RelativePath, fmt.Errorf("failed to scan file: %s", err))
			}

			for _, detector := range allDetectors {
				if ctx.Err() != nil {
					return ctx.Err()
				}

				active, isActive := activeDetectors[detector]
				if !isActive {
					continue
				}

				if file.Size() == 0 {
					continue
				}

				wasConsumed, err := detector.Detector.ProcessFile(file, active.Path, active.Report)
				if err != nil {
					log.Debug().Msgf("failed to process file %s for detector %s: %s", file.RelativePath, detector.Type, err)
					report.AddError(file.RelativePath, fmt.Errorf("failed to process file for detector %s: %s", detector.Type, err))
					continue
				}

				if wasConsumed {
					break
				}
			}

			return nil
		},
	); err != nil {
		return err
	}

	return nil
}

func isParentedBy(rootPath, path string) bool {
	relativePath, err := filepath.Rel(rootPath, path)
	if err != nil {
		return false
	}

	// directory can have one letter names like "R"
	if len(relativePath) == 1 {
		return true
	} else {
		return relativePath[:2] != ".."
	}
}

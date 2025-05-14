package implementation

import (
	"github.com/moonkit02/dearer/pkg/commands/process/filelist/files"
	"github.com/moonkit02/dearer/pkg/commands/process/orchestrator"
	"github.com/moonkit02/dearer/pkg/commands/process/orchestrator/work"
	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/engine"
	"github.com/moonkit02/dearer/pkg/scanner/language"
	"github.com/moonkit02/dearer/pkg/scanner/stats"
)

type implementation struct {
	languages    []language.Language
	orchestrator *orchestrator.Orchestrator
}

func New(languages []language.Language) engine.Engine {
	return &implementation{languages: languages}
}

func (engine *implementation) GetLanguages() []language.Language {
	return engine.languages
}

func (engine *implementation) GetLanguageById(id string) language.Language {
	for _, language := range engine.languages {
		if language.ID() == id {
			return language
		}
	}

	return nil
}

func (engine *implementation) Initialize(logLevel string) error {
	return nil
}

func (engine *implementation) LoadRule(yamlDefinition string) error {
	return nil
}

func (engine *implementation) Scan(
	config *settings.Config,
	stats *stats.Stats,
	reportPath,
	targetPath string,
	files []files.File,
) error {
	if engine.orchestrator == nil {
		var err error
		engine.orchestrator, err = orchestrator.New(work.Repository{Dir: targetPath}, config, stats, len(files))
		if err != nil {
			return err
		}
	}

	return engine.orchestrator.Scan(reportPath, files)
}

func (engine *implementation) Close() {
	if engine.orchestrator != nil {
		engine.orchestrator.Close()
	}
}

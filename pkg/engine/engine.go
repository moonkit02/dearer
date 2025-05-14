package engine

import (
	"github.com/moonkit02/dearer/pkg/commands/process/filelist/files"
	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/scanner/language"
	"github.com/moonkit02/dearer/pkg/scanner/stats"
)

type Engine interface {
	GetLanguages() []language.Language
	GetLanguageById(id string) language.Language
	Initialize(logLevel string) error
	LoadRule(yamlDefinition string) error
	Scan(config *settings.Config, stats *stats.Stats, reportPath, targetPath string, files []files.File) error
	Close()
}

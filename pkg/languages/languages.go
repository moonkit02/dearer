package languages

import (
	"github.com/moonkit02/dearer/pkg/languages/golang"
	"github.com/moonkit02/dearer/pkg/languages/java"
	"github.com/moonkit02/dearer/pkg/languages/javascript"
	"github.com/moonkit02/dearer/pkg/languages/php"
	"github.com/moonkit02/dearer/pkg/languages/python"
	"github.com/moonkit02/dearer/pkg/languages/ruby"
	"github.com/moonkit02/dearer/pkg/scanner/language"
)

func Default() []language.Language {
	return []language.Language{
		golang.Get(),
		java.Get(),
		javascript.Get(),
		php.Get(),
		python.Get(),
		ruby.Get(),
	}
}

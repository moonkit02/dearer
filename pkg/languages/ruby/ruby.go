package ruby

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/moonkit02/dearer/pkg/classification/schema"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/moonkit02/dearer/pkg/scanner/ast/query"
	"github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	detectortypes "github.com/moonkit02/dearer/pkg/scanner/detectors/types"

	"github.com/moonkit02/dearer/pkg/languages/ruby/analyzer"
	"github.com/moonkit02/dearer/pkg/languages/ruby/detectors/object"
	stringdetector "github.com/moonkit02/dearer/pkg/languages/ruby/detectors/string"
	"github.com/moonkit02/dearer/pkg/languages/ruby/pattern"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/datatype"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/insecureurl"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/stringliteral"
	"github.com/moonkit02/dearer/pkg/scanner/language"
)

type implementation struct {
	pattern pattern.Pattern
}

func Get() language.Language {
	return &implementation{}
}

func (*implementation) ID() string {
	return "ruby"
}

func (*implementation) DisplayName() string {
	return "Ruby"
}

func (*implementation) EnryLanguages() []string {
	return []string{"Ruby"}
}

func (*implementation) GoclocLanguages() []string {
	return []string{"Ruby"}
}

func (*implementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorRuby, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*implementation) SitterLanguage() *sitter.Language {
	return ruby.GetLanguage()
}

func (language *implementation) Pattern() language.Pattern {
	return &language.pattern
}

func (*implementation) NewAnalyzer(builder *tree.Builder) language.Analyzer {
	return analyzer.New(builder)
}

package javascript

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/tsx"

	"github.com/moonkit02/dearer/pkg/classification/schema"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/moonkit02/dearer/pkg/scanner/ast/query"
	"github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	detectortypes "github.com/moonkit02/dearer/pkg/scanner/detectors/types"

	"github.com/moonkit02/dearer/pkg/languages/javascript/analyzer"
	"github.com/moonkit02/dearer/pkg/languages/javascript/detectors/object"
	stringdetector "github.com/moonkit02/dearer/pkg/languages/javascript/detectors/string"
	"github.com/moonkit02/dearer/pkg/languages/javascript/pattern"
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
	return "javascript"
}

func (*implementation) DisplayName() string {
	return "JavaScript"
}

func (*implementation) EnryLanguages() []string {
	return []string{"JavaScript", "TypeScript", "TSX"}
}

func (*implementation) GoclocLanguages() []string {
	return []string{"JavaScript", "TypeScript", "JSX"}
}

func (*implementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorJavascript, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*implementation) SitterLanguage() *sitter.Language {
	return tsx.GetLanguage()
}

func (language *implementation) Pattern() language.Pattern {
	return &language.pattern
}

func (*implementation) NewAnalyzer(builder *tree.Builder) language.Analyzer {
	return analyzer.New(builder)
}

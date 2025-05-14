package java

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"

	"github.com/moonkit02/dearer/pkg/classification/schema"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/moonkit02/dearer/pkg/scanner/ast/query"
	"github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	detectortypes "github.com/moonkit02/dearer/pkg/scanner/detectors/types"

	"github.com/moonkit02/dearer/pkg/languages/java/analyzer"
	"github.com/moonkit02/dearer/pkg/languages/java/detectors/object"
	stringdetector "github.com/moonkit02/dearer/pkg/languages/java/detectors/string"
	"github.com/moonkit02/dearer/pkg/languages/java/pattern"
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
	return "java"
}

func (*implementation) DisplayName() string {
	return "Java"
}

func (*implementation) EnryLanguages() []string {
	return []string{"Java"}
}

func (*implementation) GoclocLanguages() []string {
	return []string{"Java"}
}

func (*implementation) NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector {
	return []detectortypes.Detector{
		object.New(querySet),
		datatype.New(detectors.DetectorJava, schemaClassifier),
		stringdetector.New(querySet),
		stringliteral.New(querySet),
		insecureurl.New(querySet),
	}
}

func (*implementation) SitterLanguage() *sitter.Language {
	return java.GetLanguage()
}

func (language *implementation) Pattern() language.Pattern {
	return &language.pattern
}

func (*implementation) NewAnalyzer(builder *tree.Builder) language.Analyzer {
	return analyzer.New(builder)
}

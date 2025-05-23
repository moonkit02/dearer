package language

import (
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/moonkit02/dearer/pkg/classification/schema"
	"github.com/moonkit02/dearer/pkg/scanner/ast/query"
	"github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	detectortypes "github.com/moonkit02/dearer/pkg/scanner/detectors/types"
)

type Language interface {
	ID() string
	DisplayName() string
	EnryLanguages() []string
	GoclocLanguages() []string
	NewBuiltInDetectors(schemaClassifier *schema.Classifier, querySet *query.Set) []detectortypes.Detector
	SitterLanguage() *sitter.Language
	Pattern() Pattern
	NewAnalyzer(builder *tree.Builder) Analyzer
}

type Analyzer interface {
	Analyze(node *sitter.Node, visitChildren func() error) error
}

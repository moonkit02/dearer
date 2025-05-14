package string

import (
	"github.com/moonkit02/dearer/pkg/scanner/ast/query"
	"github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	"github.com/moonkit02/dearer/pkg/scanner/ruleset"
	"github.com/moonkit02/dearer/pkg/util/stringutil"

	"github.com/moonkit02/dearer/pkg/scanner/detectors/common"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/types"
)

type stringDetector struct {
	types.DetectorBase
}

func New(querySet *query.Set) types.Detector {
	return &stringDetector{}
}

func (detector *stringDetector) Rule() *ruleset.Rule {
	return ruleset.BuiltinStringRule
}

func (detector *stringDetector) DetectAt(
	node *tree.Node,
	detectorContext types.Context,
) ([]interface{}, error) {
	switch node.Type() {
	case "string_literal", "character_literal":
		return []interface{}{common.String{
			Value:     stringutil.StripQuotes(node.Content()),
			IsLiteral: true,
		}}, nil
	case "binary_expression":
		if node.Children()[1].Content() == "+" {
			return common.ConcatenateChildStrings(node, detectorContext)
		}
	case "assignment_expression":
		if node.Children()[1].Content() == "+=" {
			return common.ConcatenateAssignEquals(node, detectorContext)
		}
	}

	return nil, nil
}

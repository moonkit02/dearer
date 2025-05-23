package context

import (
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/report/variables"
)

func variableAssignmentExpression(node *parser.Node) *variables.Variable {
	left := node.ChildByFieldName("left")
	if left == nil {
		return nil
	}

	right := node.ChildByFieldName("right")
	if right == nil {
		return nil
	}

	if left.Type() != "variable_name" {
		return nil
	}

	variable := resolveBaseNode(right)
	variable.Name = left.Child(0).Content()

	return variable
}

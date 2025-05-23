package language

import "github.com/moonkit02/dearer/pkg/scanner/ast/tree"

type PatternBase struct{}

func (*PatternBase) IsLeaf(node *tree.Node) bool {
	return false
}

func (*PatternBase) TranslateContent(fromNodeType, toNodeType, content string) string {
	return content
}

func (*PatternBase) IsRoot(node *tree.Node) bool {
	return true
}

func (*PatternBase) ShouldSkipNode(node *tree.Node) bool {
	return false
}

func (*PatternBase) IsContainer(node *tree.Node) bool {
	return false
}

func (*PatternBase) FixupVariableDummyValue(input []byte, node *tree.Node, dummyValue string) string {
	return dummyValue
}

func (*PatternBase) AnonymousParentTypes() []string {
	return nil
}

func (*PatternBase) AdjustInput(input string) string {
	return input
}

func (*PatternBase) FixupMissing(node *tree.Node) string {
	return ""
}

func (*PatternBase) IsVariable(node *tree.Node, dummyValue string) bool {
	return node.Content() == dummyValue
}

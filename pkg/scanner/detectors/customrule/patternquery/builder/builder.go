package builder

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	sitter "github.com/smacker/go-tree-sitter"

	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/scanner/ast"
	"github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	asttree "github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/customrule/patternquery/builder/bytereplacer"
	"github.com/moonkit02/dearer/pkg/scanner/language"
)

type InputParams struct {
	VariableNames     []string
	Variables         []language.PatternVariable
	MatchNodeOffset   int
	UnanchoredOffsets []int
}

type Result struct {
	Query           string
	VariableNames   []string
	ParamToVariable map[string]string
	EqualParams     [][]string
	ParamToContent  map[string]map[string]string
	RootVariable    *language.PatternVariable
}

type builder struct {
	sitterLanguage   *sitter.Language
	patternLanguage  language.Pattern
	stringBuilder    strings.Builder
	idGenerator      nodeid.Generator
	inputParams      InputParams
	variableToParams map[string][]string
	paramToContent   map[string]map[string]string
	matchNode        *asttree.Node
}

func Build(
	language language.Language,
	input string,
	focusedVariable string,
) (*Result, error) {
	patternLanguage := language.Pattern()
	processedInput, inputParams, err := processInput(patternLanguage, input)
	if err != nil {
		return nil, err
	}

	tree, err := ast.Parse(context.TODO(), language, processedInput)
	if err != nil {
		return nil, err
	}

	fixupResult, err := fixupInput(
		patternLanguage,
		processedInput,
		inputParams.Variables,
		tree.RootNode(),
	)
	if err != nil {
		return nil, err
	}

	if fixupResult.Changed() {
		if log.Trace().Enabled() {
			log.Trace().Msgf("fixedInput -> %s", fixupResult.Value())
		}

		tree, err = ast.Parse(context.TODO(), language, fixupResult.Value())
		if err != nil {
			return nil, err
		}

		inputParams.MatchNodeOffset = fixupResult.Translate(inputParams.MatchNodeOffset)
		for i := range inputParams.UnanchoredOffsets {
			inputParams.UnanchoredOffsets[i] = fixupResult.Translate(inputParams.UnanchoredOffsets[i])
		}
	}

	root := tree.RootNode()

	var foundRoot bool
	root.Walk(func(rootNode *asttree.Node, visitChildren func() error) error { //nolint:errcheck
		if foundRoot {
			return nil
		}

		if patternLanguage.IsRoot(rootNode) {
			root = rootNode
			foundRoot = true
			return nil
		} else {
			return visitChildren()
		}
	})

	builder := builder{
		sitterLanguage:   language.SitterLanguage(),
		patternLanguage:  patternLanguage,
		stringBuilder:    strings.Builder{},
		idGenerator:      &nodeid.IntGenerator{},
		inputParams:      *inputParams,
		variableToParams: make(map[string][]string),
		paramToContent:   make(map[string]map[string]string),
	}

	builder.setMatchNode(
		inputParams.MatchNodeOffset,
		focusedVariable,
		tree.RootNode(),
	)
	if builder.matchNode == nil {
		return nil, fmt.Errorf("match node not found")
	}

	result, err := builder.build(root)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func fixupInput(
	patternLanguage language.Pattern,
	byteInput []byte,
	variables []language.PatternVariable,
	rootNode *asttree.Node,
) (*bytereplacer.Result, error) {
	replacer := bytereplacer.New(byteInput)
	insideError := false

	err := rootNode.Walk(func(node *asttree.Node, visitChildren func() error) error {
		oldInsideError := insideError
		if node.IsError() {
			insideError = true
		}
		if err := visitChildren(); err != nil {
			return err
		}
		insideError = oldInsideError

		if !insideError && !node.IsMissing() {
			return nil
		}

		var newValue string

		if insideError {
			variable := getVariableFor(node, patternLanguage, variables)
			if variable == nil {
				return nil
			}

			if log.Trace().Enabled() {
				log.Trace().Msgf("attempting pattern fixup. node: %s", node.Debug())
			}

			newValue = patternLanguage.FixupVariableDummyValue(byteInput, node, variable.DummyValue)
			if newValue == variable.DummyValue {
				return nil
			}
			variable.DummyValue = newValue
		} else {
			if log.Trace().Enabled() {
				log.Trace().Msgf("attempting pattern fixup (missing node). node: %s", node.Debug())
			}

			newValue = patternLanguage.FixupMissing(node)
			if newValue == "" {
				return nil
			}
		}

		return replacer.Replace(node.ContentStart.Byte, node.ContentEnd.Byte, []byte(newValue))
	})

	return replacer.Done(), err
}

func (builder *builder) build(rootNode *asttree.Node) (*Result, error) {
	if len(rootNode.Children()) == 0 {
		variable := builder.getVariableFor(rootNode)
		if variable != nil {
			return &Result{RootVariable: variable}, nil
		}
	}

	builder.write("(")

	if err := builder.compileNode(rootNode, true, false); err != nil {
		return nil, err
	}

	builder.write(" @root")
	builder.write(")")

	paramToVariable, equalParams := builder.processVariableToParams()

	return &Result{
		Query:           builder.stringBuilder.String(),
		VariableNames:   builder.inputParams.VariableNames,
		ParamToVariable: paramToVariable,
		EqualParams:     equalParams,
		ParamToContent:  builder.paramToContent,
	}, nil
}

func (builder *builder) compileNode(node *asttree.Node, isRoot bool, isLastChild bool) error {
	if node.SitterNode().IsError() {
		return fmt.Errorf(
			"error parsing pattern at %d:%d: %s",
			node.ContentStart.Line,
			node.ContentStart.Column,
			node.Content(),
		)
	}

	nodeAnchoredBefore, nodeAnchoredAfter := builder.patternLanguage.IsAnchored(node)
	anchored := !isRoot && node.IsNamed() && nodeAnchoredBefore

	if anchored && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.ContentStart.Byte) {
		builder.write(". ")
	}

	if variable := builder.getVariableFor(node); variable != nil {
		builder.compileVariableNode(node, variable)
	} else if !node.IsNamed() {
		builder.compileAnonymousNode(node)
	} else if len(node.NamedChildren()) == 0 || builder.patternLanguage.IsLeaf(node) {
		builder.compileLeafNode(node)
	} else if err := builder.compileNodeWithChildren(node); err != nil {
		return err
	}

	if node == builder.matchNode {
		builder.write(" @match")
	}

	if anchored && isLastChild && nodeAnchoredAfter && !slices.Contains(builder.inputParams.UnanchoredOffsets, node.ContentEnd.Byte) {
		builder.write(" .")
	}

	return nil
}

// variable nodes match their type and capture their content
func (builder *builder) compileVariableNode(node *tree.Node, variable *language.PatternVariable) {
	if fieldName := fieldNameFor(builder.sitterLanguage, node); fieldName != "" {
		builder.write(fieldName)
		builder.write(": ")
	}

	if variable.Name == "_" {
		builder.write("(_)")
		return
	}

	paramName := builder.newParam()
	builder.variableToParams[variable.Name] = append(builder.variableToParams[variable.Name], paramName)

	builder.write("[")

	for _, nodeType := range variable.NodeTypes {
		builder.write("(")
		builder.write(nodeType)
		builder.write(")")
	}

	builder.write("] @")
	builder.write(paramName)
}

// Anonymous nodes match their content as a literal
func (builder *builder) compileAnonymousNode(node *asttree.Node) {
	if !slices.Contains(builder.patternLanguage.AnonymousParentTypes(), node.Parent().Type()) {
		return
	}

	builder.write(strconv.Quote(node.Content()))
}

// Leaves match their type and content
func (builder *builder) compileLeafNode(node *asttree.Node) {
	if !slices.Contains(builder.patternLanguage.LeafContentTypes(), node.Type()) {
		builder.write("[")

		for _, nodeType := range builder.patternLanguage.NodeTypes(node) {
			builder.write(" (")
			builder.write(nodeType)
			builder.write(" )")
		}

		builder.write("]")
		return
	}

	paramName := builder.newParam()
	paramContent := make(map[string]string)
	builder.paramToContent[paramName] = paramContent

	builder.write("[")

	for _, nodeType := range builder.patternLanguage.NodeTypes(node) {
		paramContent[nodeType] = builder.patternLanguage.TranslateContent(
			node.Type(),
			nodeType, node.Content(),
		)

		builder.write(" (")
		builder.write(nodeType)
		builder.write(" )")
	}

	builder.write("] @")
	builder.write(paramName)
}

// Nodes with children match their type and child nodes
func (builder *builder) compileNodeWithChildren(node *asttree.Node) error {
	builder.write("[")

	var children []*asttree.Node
	if slices.Contains(builder.patternLanguage.AnonymousParentTypes(), node.Type()) {
		children = node.Children()
	} else {
		children = node.NamedChildren()
	}

	lastNode := children[len(children)-1]

	for _, nodeType := range builder.patternLanguage.NodeTypes(node) {
		builder.write("(")
		builder.write(nodeType)

		for _, child := range node.Children() {
			builder.write(" ")

			if err := builder.compileNode(child, false, child == lastNode); err != nil {
				return err
			}
		}

		builder.write(")")
	}

	builder.write("]")

	return nil
}

func (builder *builder) processVariableToParams() (map[string]string, [][]string) {
	paramToVariable := make(map[string]string)
	var equalParams [][]string

	for variableName, paramNames := range builder.variableToParams {
		if len(paramNames) > 1 {
			equalParams = append(equalParams, paramNames)
		}

		for _, paramName := range paramNames {
			paramToVariable[paramName] = variableName
		}
	}

	return paramToVariable, equalParams
}

func (builder *builder) getVariableFor(node *asttree.Node) *language.PatternVariable {
	return getVariableFor(node, builder.patternLanguage, builder.inputParams.Variables)
}

func getVariableFor(
	node *asttree.Node,
	patternLanguage language.Pattern,
	variables []language.PatternVariable,
) *language.PatternVariable {
	if patternLanguage.IsContainer(node) {
		return nil
	}

	for i, variable := range variables {
		if patternLanguage.IsVariable(node, variable.DummyValue) {
			return &variables[i]
		}
	}

	return nil
}

func (builder *builder) write(value string) {
	builder.stringBuilder.WriteString(value)
}

func (builder *builder) newParam() string {
	return "param" + builder.idGenerator.GenerateId()
}

func (builder *builder) setMatchNode(
	offset int,
	focusedVariable string,
	node *asttree.Node,
) {
	err := node.Walk(func(node *asttree.Node, visitChildren func() error) error {
		if focusedVariable != "" {
			if variable := builder.getVariableFor(node); variable != nil && variable.Name == focusedVariable {
				builder.matchNode = node
				return nil
			}
		} else {
			if node.ContentStart.Byte == offset && !builder.patternLanguage.IsContainer(node) {
				builder.matchNode = node
				return nil
			}
		}

		return visitChildren()
	})

	// walk itself shouldn't trigger an error, and we aren't creating any
	if err != nil {
		panic(err)
	}
}

func fieldNameFor(sitterLanguage *sitter.Language, node *tree.Node) string {
	parent := node.Parent()
	if parent == nil {
		return ""
	}

	// the following is a workaround until
	// https://github.com/tree-sitter/tree-sitter/pull/2104 is released
	for i := 1; ; i++ {
		name := sitterLanguage.FieldName(i)
		if name == "" {
			return ""
		}

		if parent.ChildByFieldName(name) == node {
			return name
		}
	}
}

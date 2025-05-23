package tsx

import (
	"strings"

	"github.com/smacker/go-tree-sitter/typescript/tsx"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/moonkit02/dearer/pkg/detectors/types"
	typescript_datatype "github.com/moonkit02/dearer/pkg/detectors/typescript/datatype"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/interfacedetector"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/moonkit02/dearer/pkg/report/values"
	"github.com/moonkit02/dearer/pkg/report/variables"
	"github.com/moonkit02/dearer/pkg/util/file"
)

var (
	environmentVariablesQuery = parser.QueryMustCompile(tsx.GetLanguage(), `
	(member_expression
		object: (member_expression) @object
		property: (property_identifier) @key) @node

	(subscript_expression
		object: (member_expression) @object
		index: (string) @key) @node

	(variable_declarator
		name: (object_pattern (shorthand_property_identifier_pattern) @key @node)
			value: (member_expression) @object)
`)
)

type detector struct {
	idGenerator nodeid.Generator
}

func New(idGenerator nodeid.Generator) types.Detector {
	return &detector{
		idGenerator: idGenerator,
	}
}

func (detector *detector) AcceptDir(dir *file.Path) (bool, error) {
	return true, nil
}

func (detector *detector) ProcessFile(file *file.FileInfo, dir *file.Path, report report.Report) (bool, error) {
	if file.Language != "TSX" {
		return false, nil
	}

	tree, err := parser.ParseFile(file, file.Path, tsx.GetLanguage())
	if err != nil {
		return false, err
	}
	defer tree.Close()

	typescript_datatype.Discover(report, tree, tsx.GetLanguage(), detector.idGenerator)

	if err := annotate(tree, environmentVariablesQuery); err != nil {
		return false, err
	}

	if err := interfacedetector.Detect(&interfacedetector.Request{
		Tree:             tree,
		Report:           report,
		AcceptExpression: acceptExpression,
		DetectorType:     detectors.DetectorTsx,
		PathAllowed:      false,
	}); err != nil {
		return false, err
	}

	return true, nil
}

func annotate(tree *parser.Tree, environmentVariablesQuery *sitter.Query) error {
	if err := annotateEnvironmentVariables(tree, environmentVariablesQuery); err != nil {
		return err
	}

	return tree.Annotate(func(node *parser.Node, value *values.Value) {
		switch node.Type() {
		case "template_substitution":
			value.Append(node.FirstChild().Value())

			return
		case "binary_expression":
			if node.FirstUnnamedChild().Content() == "+" {
				value.Append(node.ChildByFieldName("left").Value())
				value.Append(node.ChildByFieldName("right").Value())

				return
			}
		case "identifier", "property_identifier":
			value.AppendVariableReference(variables.VariableName, node.Content())

			return

		case "string", "template_string":
			node.EachPart(func(text string) error { //nolint:all,errcheck
				value.AppendString(text)

				return nil
			}, func(child *parser.Node) error {
				value.Append(child.Value())

				return nil
			})

			return
		}

		value.AppendUnknown(node.ChildValueParts())
	})
}

func annotateEnvironmentVariables(tree *parser.Tree, query *sitter.Query) error {
	return tree.Query(query, func(captures parser.Captures) error {
		if captures["object"].Content() != "process.ENV" {
			return nil
		}

		node := captures["node"]
		keyNode := captures["key"]
		key := stripQuotes(keyNode.Content())

		value := values.New()
		value.AppendVariableReference(variables.VariableEnvironment, key)
		node.SetValue(value)

		return nil
	})
}

func stripQuotes(value string) string {
	return strings.Trim(value, `'"`)
}

func acceptExpression(node *parser.Node) bool {
	lastNode := node
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		switch parent.Type() {
		case "decorator":
			// @MyDecorator("ignored")
			return false
		case "pair":
			if parent.ChildByFieldName("key").Equal(lastNode) {
				// { 'ignored.domain': '...' }
				return false
			}
		case "import_statement":
			// import * from 'ignored'
			return false
		case "subscript_expression":
			// something['ignored.domain']
			return false
		case "jsx_element", "jsx_self_closing_element":
			// <img src="ignored.domain"/>
			return false
		}

		lastNode = parent
	}

	return true
}

package run

import (
	"os"
	"regexp"

	"github.com/moonkit02/dearer/pkg/classification/schema"
	"github.com/moonkit02/dearer/pkg/commands"
	"github.com/moonkit02/dearer/pkg/engine"
	"github.com/moonkit02/dearer/pkg/engine/implementation"
	"github.com/moonkit02/dearer/pkg/languages"
	"github.com/moonkit02/dearer/pkg/scanner/ast/query"
	"github.com/moonkit02/dearer/pkg/scanner/ast/tree"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/common"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/types"
	"github.com/moonkit02/dearer/pkg/scanner/language"
	"github.com/moonkit02/dearer/pkg/scanner/ruleset"
	"github.com/moonkit02/dearer/pkg/util/regex"
	"github.com/moonkit02/dearer/pkg/util/stringutil"
)

type Object = common.Object
type Property = common.Property
type String = common.String

type Engine = engine.Engine

type Analyzer = language.Analyzer
type Language = language.Language
type Pattern = language.Pattern
type PatternBase = language.PatternBase
type PatternVariable = language.PatternVariable
type Scope = language.Scope

type Query = query.Query
type Set = query.Set

type Rule = ruleset.Rule

type Classifier = schema.Classifier

type Builder = tree.Builder
type Node = tree.Node

type Context = types.Context
type Detection = types.Detection
type Detector = types.Detector
type DetectorBase = types.DetectorBase

var BuiltinObjectRule = ruleset.BuiltinObjectRule
var BuiltinStringRule = ruleset.BuiltinStringRule

func GetNonVirtualObjects(detectorContext types.Context, node *tree.Node) ([]*types.Detection, error) {
	return common.GetNonVirtualObjects(detectorContext, node)
}

func StripQuotes(input string) string {
	return stringutil.StripQuotes(input)
}

func ConcatenateChildStrings(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	return common.ConcatenateChildStrings(node, detectorContext)
}

func ConcatenateAssignEquals(node *tree.Node, detectorContext types.Context) ([]interface{}, error) {
	return common.ConcatenateAssignEquals(node, detectorContext)
}

func ProjectObject(node *tree.Node, detectorContext types.Context, objectNode *tree.Node, objectName string, propertyName string, isPropertyAccess bool) ([]interface{}, error) {
	return common.ProjectObject(node, detectorContext, objectNode, objectName, propertyName, isPropertyAccess)
}

func ReplaceAllWithSubmatches(pattern *regexp.Regexp, input string, replace func(submatches []string) (string, error)) (string, error) {
	return regex.ReplaceAllWithSubmatches(pattern, input, replace)
}

func NewScope(parent *language.Scope) *language.Scope {
	return language.NewScope(parent)
}

func NewEngine(languages []Language) Engine {
	return implementation.New(languages)
}

func DefaultLanguages() []Language {
	return languages.Default()
}

func Run(version, commitSHA string, engine Engine) {
	err := commands.NewApp(version, commitSHA, engine).Execute()

	if err != nil {
		// error messages are printed by the framework
		os.Exit(1)
	}
}

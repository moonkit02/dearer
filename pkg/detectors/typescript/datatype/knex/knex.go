package knex

import (
	"github.com/moonkit02/dearer/pkg/detectors/javascript/util"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/report"
	sitter "github.com/smacker/go-tree-sitter"
)

func Discover(report report.Report, tree *parser.Tree, language *sitter.Language) {
	knexImports := util.GetImports(tree, language, []string{"knex"})

	if len(knexImports) == 0 {
		return
	}

	detectFunctionTypes(report, tree, language, knexImports)
	detectTableDeclarationModule(report, tree, language)
}

package cache

import (
	"github.com/smacker/go-tree-sitter/ruby"

	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/moonkit02/dearer/pkg/report/frameworks/rails"
	"github.com/moonkit02/dearer/pkg/util/file"
)

var (
	language = ruby.GetLanguage()

	query = parser.QueryMustCompile(language, `
		(assignment
			left: (call) @target
			right: (right_assignment_list (simple_symbol) @type)
			(#match? @target "^config\\.cache_store$")) @node
	`)
)

func ExtractCaches(file *file.FileInfo, report report.Report) error {
	tree, err := parser.ParseFile(file, file.Path, language)
	if err != nil {
		return err
	}
	defer tree.Close()

	return tree.Query(query, func(captures parser.Captures) error {
		typeSymbol := captures["type"]
		typeName := typeSymbol.Content()[1:]

		report.AddFramework(detectors.DetectorRails, rails.TypeCache, rails.Cache{
			Type: typeName,
		}, typeSymbol.Source(false))

		return nil
	})
}

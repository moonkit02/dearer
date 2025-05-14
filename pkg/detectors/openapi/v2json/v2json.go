package v2json

import (
	"github.com/moonkit02/dearer/pkg/detectors/openapi/json"
	"github.com/moonkit02/dearer/pkg/detectors/openapi/queries"
	"github.com/moonkit02/dearer/pkg/detectors/openapi/reportadder"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	reporttypes "github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/report/operations/operationshelper"
	"github.com/moonkit02/dearer/pkg/report/schema/schemahelper"
	"github.com/moonkit02/dearer/pkg/util/file"
	"github.com/smacker/go-tree-sitter/javascript"
)

var queryParameters = parser.QueryMustCompile(javascript.GetLanguage(), `
(_
	(
      pair
        key:
            (string) @helperName
            (#match? @helperName "^\"name\"$")
         value:
            (string) @param_name
	)
	(
      pair
        key:
            (string) @helperType
            (#match? @helperType "^\"type\"$")
         value:
            (string) @param_type
	)
)
`)

func ProcessFile(idGenerator nodeid.Generator, file *file.FileInfo, report reporttypes.Report) (bool, error) {

	tree, err := parser.ParseFile(file, file.Path, javascript.GetLanguage())
	if err != nil {
		return false, err
	}
	defer tree.Close()

	foundSchemas := make(map[parser.Node]*schemahelper.Schema)

	nodeIDMap := nodeid.New(tree, idGenerator)
	nodeIDMap.Annotate()

	err = queries.AnnotateV2Paramaters(nodeIDMap, tree, foundSchemas, queryParameters)
	if err != nil {
		return false, err
	}

	err = json.AnnotateOperationId(nodeIDMap, tree, foundSchemas)
	if err != nil {
		return false, err
	}

	err = json.AnnotateObjects(nodeIDMap, tree, foundSchemas)
	if err != nil {
		return false, err
	}

	foundPaths := make(map[parser.Node]*operationshelper.Operation)
	err = json.AnnotatePaths(tree, foundPaths)
	if err != nil {
		return false, err
	}

	servers := findServers(tree)
	reportadder.AddOperations(file, report, foundPaths, servers)
	reportadder.AddSchema(file, report, foundSchemas, idGenerator)

	return true, err
}

package v3yaml

import (
	"github.com/moonkit02/dearer/pkg/detectors/openapi/queries"
	"github.com/moonkit02/dearer/pkg/detectors/openapi/reportadder"
	yamlparser "github.com/moonkit02/dearer/pkg/detectors/openapi/yaml"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	reporttypes "github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/report/operations/operationshelper"
	"github.com/moonkit02/dearer/pkg/report/schema/schemahelper"
	"github.com/moonkit02/dearer/pkg/util/file"
	"github.com/smacker/go-tree-sitter/yaml"
)

var queryParameters = parser.QueryMustCompile(yaml.GetLanguage(), `
(_
	(
      block_mapping_pair
        key:
            (flow_node) @helperName
            (#match? @helperName "^name$")
         value:
            (flow_node) @param_name
	)
	(
      block_mapping_pair
        key:
            (flow_node) @helperSchema
            (#match? @helperSchema "^schema$")
         value:
            (block_node (block_mapping) @param_schema)
	)
)
`)

func ProcessFile(idGenerator nodeid.Generator, file *file.FileInfo, report reporttypes.Report) (bool, error) {
	tree, err := parser.ParseFile(file, file.Path, yaml.GetLanguage())
	if err != nil {
		return false, err
	}
	defer tree.Close()

	nodeIDMap := nodeid.New(tree, idGenerator)
	nodeIDMap.Annotate()

	foundSchemas := make(map[parser.Node]*schemahelper.Schema)
	err = queries.AnnotateV3Paramaters(nodeIDMap, tree, foundSchemas, queryParameters)
	if err != nil {
		return false, err
	}

	err = yamlparser.AnnotateOperationId(nodeIDMap, tree, foundSchemas)
	if err != nil {
		return false, err
	}

	err = yamlparser.AnnotateObjects(nodeIDMap, tree, foundSchemas)
	if err != nil {
		return false, err
	}

	foundPaths := make(map[parser.Node]*operationshelper.Operation)
	err = yamlparser.AnnotatePaths(tree, foundPaths)
	if err != nil {
		return false, err
	}

	servers := queries.FindUrls(file)
	reportadder.AddOperations(file, report, foundPaths, servers)
	reportadder.AddSchema(file, report, foundSchemas, idGenerator)

	return true, err
}

package queries

import (
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/report/schema"
	"github.com/moonkit02/dearer/pkg/report/schema/schemahelper"
	"github.com/moonkit02/dearer/pkg/util/stringutil"
	sitter "github.com/smacker/go-tree-sitter"
)

func AnnotateV2Paramaters(nodeIDMap *nodeid.Map, tree *parser.Tree, foundValues map[parser.Node]*schemahelper.Schema, query *sitter.Query) error {

	captures := tree.QueryMustPass(query)

	for _, capture := range captures {
		if capture["param_type"] == nil || capture["param_name"] == nil {
			continue
		}

		if stringutil.StripQuotes(capture["helperName"].Content()) != "name" ||
			stringutil.StripQuotes(capture["helperType"].Content()) != "type" {
			continue
		}

		nameNode := capture["param_name"]
		typeNode := capture["param_type"]

		foundValues[*nameNode] = &schemahelper.Schema{
			Source: nameNode.Source(true),
			Value: schema.Schema{
				FieldName: nameNode.Content(),
				FieldType: typeNode.Content(),
				FieldUUID: nodeIDMap.ValueForNode(nameNode.ID()),
			},
		}

	}

	return nil
}

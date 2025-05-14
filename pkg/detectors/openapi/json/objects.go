package json

import (
	"github.com/moonkit02/dearer/pkg/detectors/openapi/queries"
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/report/schema/schemahelper"
	"github.com/smacker/go-tree-sitter/javascript"
)

var queryObjects = parser.QueryMustCompile(javascript.GetLanguage(), `
(_
	(
      pair
        key:
            (string) @param_object_name
         value:
            (object
            	(pair
                	key:
                    	(string) @helperProperties
                        (#match? @helperProperties "^\"properties\"$")
                    value:
                    	(object) @param_object_properties
                )
            )
	)
)
`)

type ObjectChildMatcher struct {
}

func (childMatcher ObjectChildMatcher) Match(input *parser.Node) *parser.Node {
	return input
}

func AnnotateObjects(nodeIDMap *nodeid.Map, tree *parser.Tree, foundValues map[parser.Node]*schemahelper.Schema) error {
	return queries.AnnotateObjects(queries.ObjectsRequest{
		Tree:        tree,
		Query:       queryObjects,
		FoundValues: foundValues,
		ChildMatch:  ObjectChildMatcher{},
		NodeIDMap:   nodeIDMap,
	})
}

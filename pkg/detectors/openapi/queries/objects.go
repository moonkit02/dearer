package queries

import (
	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/report/schema"
	"github.com/moonkit02/dearer/pkg/report/schema/schemahelper"
	"github.com/moonkit02/dearer/pkg/util/stringutil"
	sitter "github.com/smacker/go-tree-sitter"
)

type ObjectsRequest struct {
	Tree        *parser.Tree
	Query       *sitter.Query
	FoundValues map[parser.Node]*schemahelper.Schema
	ChildMatch  ChildMatch
	NodeIDMap   *nodeid.Map
}

func AnnotateObjects(request ObjectsRequest) error {
	captures := request.Tree.QueryMustPass(request.Query)

	for _, capture := range captures {
		if capture["param_object_name"] == nil || capture["param_object_properties"] == nil {
			continue
		}

		if stringutil.StripQuotes(capture["helperProperties"].Content()) != "properties" {
			continue
		}

		objectNameNode := capture["param_object_name"]
		objectPropertiesNode := capture["param_object_properties"]

		for i := 0; i < objectPropertiesNode.ChildCount(); i++ {
			property := objectPropertiesNode.Child(i)
			propertyName := property.ChildByFieldName("key")

			if propertyName == nil {
				continue
			}

			propertyValue := request.ChildMatch.Match(property.ChildByFieldName("value"))

			if propertyValue == nil {
				continue
			}

			propertyType := ""

			for j := 0; j < propertyValue.ChildCount(); j++ {
				propertyAttribute := propertyValue.Child(j)

				attributeKey := propertyAttribute.ChildByFieldName("key")

				if attributeKey == nil {
					continue
				}

				if stringutil.StripQuotes(attributeKey.Content()) == "type" {
					propertyTypeNode := propertyAttribute.ChildByFieldName("value")
					if propertyTypeNode == nil {
						continue
					}

					propertyType = propertyTypeNode.Content()
				}
			}

			request.FoundValues[*propertyName] = &schemahelper.Schema{
				Source: propertyName.Source(true),
				Value: schema.Schema{
					FieldName:  propertyName.Content(),
					FieldUUID:  request.NodeIDMap.ValueForNode(propertyName.ID()),
					FieldType:  propertyType,
					ObjectName: objectNameNode.Content(),
					ObjectUUID: request.NodeIDMap.ValueForNode(objectNameNode.ID()),
				},
			}
		}
	}

	return nil
}

package datatype

import (
	"github.com/moonkit02/dearer/pkg/parser"
	parserdatatype "github.com/moonkit02/dearer/pkg/parser/datatype"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
	"github.com/moonkit02/dearer/pkg/report"
	"github.com/moonkit02/dearer/pkg/report/detections"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	schemadatatype "github.com/moonkit02/dearer/pkg/report/schema/datatype"
)

func Discover(report report.Report, tree *parser.Tree, idGenerator nodeid.Generator) {
	datatypes := make(map[parser.NodeID]*schemadatatype.DataType)
	helperDatatypes := make(map[parser.NodeID]*schemadatatype.DataType)

	addProperties(tree, helperDatatypes)
	linkProperties(tree, datatypes, helperDatatypes)
	scopeProperties(datatypes, idGenerator)
	addObjects(tree, datatypes)

	parserdatatype.PruneMap(datatypes)

	// remove this
	for nodeID, datatype := range datatypes {
		if datatype.Name == "this" {
			delete(datatypes, nodeID)
		}
	}

	report.AddDataType(detections.TypeSchema, detectors.DetectorJavascript, idGenerator, datatypes, nil)
}

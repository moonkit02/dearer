package writer

import (
	"fmt"
	"io"
	"log"

	classification "github.com/moonkit02/dearer/pkg/classification"
	classificationschema "github.com/moonkit02/dearer/pkg/classification/schema"
	zerolog "github.com/rs/zerolog/log"

	"github.com/moonkit02/dearer/pkg/parser"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"

	"github.com/moonkit02/dearer/pkg/report/dependencies"
	"github.com/moonkit02/dearer/pkg/report/detections"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/moonkit02/dearer/pkg/report/frameworks"
	"github.com/moonkit02/dearer/pkg/report/interfaces"
	"github.com/moonkit02/dearer/pkg/report/operations"
	"github.com/moonkit02/dearer/pkg/report/schema"
	"github.com/moonkit02/dearer/pkg/report/schema/datatype"
	"github.com/moonkit02/dearer/pkg/report/secret"
	"github.com/moonkit02/dearer/pkg/report/source"

	"github.com/moonkit02/dearer/pkg/util/jsonlines"
)

type StoredSchema struct {
	Value  schema.Schema
	Source *source.Source
	Parent *parser.Node
}

type StoredSchemaNodes = map[*parser.Node]*StoredSchema

type SchemaGroup struct {
	Node         *parser.Node
	ParentSchema StoredSchema
	DetectorType detectors.Type
	Schemas      StoredSchemaNodes
}

type Detectors struct {
	Classifier    *classification.Classifier
	File          io.Writer
	StoredSchemas *SchemaGroup
}

func (report *Detectors) AddInterface(
	detectorType detectors.Type,
	data interfaces.Interface,
	source source.Source,
) {
	detection := &detections.Detection{DetectorType: detectorType, Value: data, Source: source, Type: detections.TypeInterface}
	classifiedDetection, err := report.Classifier.Interfaces.Classify(*detection)
	if err != nil {
		zerolog.Debug().Msgf("classification interfaces error from %s: %s", detection.Source.Filename, err)
		return
	}

	classifiedDetection.Type = detections.TypeInterfaceClassified
	report.Add(classifiedDetection)
}

func (report *Detectors) AddDataType(detectionType detections.DetectionType, detectorType detectors.Type, idGenerator nodeid.Generator, values map[parser.NodeID]*datatype.DataType, parent *parser.Node) {
	classifiedDatatypes := make(map[parser.NodeID]*datatype.ClassifiedDatatype, 0)
	for nodeID, target := range values {
		classification := report.Classifier.Schema.Classify(classificationschema.ClassificationRequest{
			Value:        target.ToClassificationRequestDetection(),
			Filename:     target.GetNode().Source(false).Filename,
			DetectorType: detectorType,
		})

		classifiedDatatypes[nodeID] = datatype.BuildClassifiedDatatype(target, classification)
	}

	if detectionType == detections.TypeCustom {
		datatype.ExportClassified(report, detections.TypeCustomClassified, detectorType, idGenerator, classifiedDatatypes, parent)
	} else {
		datatype.ExportClassified(report, detections.TypeSchemaClassified, detectorType, idGenerator, classifiedDatatypes, nil)
	}
}

func (report *Detectors) SchemaGroupBegin(detectorType detectors.Type, node *parser.Node, schema schema.Schema, source *source.Source, parent *parser.Node) {
	if report.SchemaGroupIsOpen() {
		zerolog.Warn().Msg("schema group already open")
	}
	report.StoredSchemas = &SchemaGroup{
		Node: node,
		ParentSchema: StoredSchema{
			Value:  schema,
			Source: source,
			Parent: parent,
		},
		DetectorType: detectorType,
		Schemas:      make(StoredSchemaNodes),
	}
}

func (report *Detectors) SchemaGroupIsOpen() bool {
	return report.StoredSchemas != nil
}

func (report *Detectors) SchemaGroupShouldClose(tableName string) bool {
	if report.StoredSchemas == nil {
		return false
	}
	return tableName != report.StoredSchemas.ParentSchema.Value.ObjectName
}

func (report *Detectors) SchemaGroupAddItem(node *parser.Node, schema schema.Schema, source *source.Source) {
	report.StoredSchemas.Schemas[node] = &StoredSchema{Value: schema, Source: source, Parent: report.StoredSchemas.ParentSchema.Parent}
}

func (report *Detectors) SchemaGroupEnd(idGenerator nodeid.Generator) {
	if !report.SchemaGroupIsOpen() {
		return
	}

	// Build child data types
	childDataTypes := map[string]datatype.DataTypable{}
	for node, storedSchema := range report.StoredSchemas.Schemas {
		schema := storedSchema.Value

		childName := schema.FieldName
		childDataTypes[childName] = &datatype.DataType{
			Node:       node,
			Name:       childName,
			Type:       schema.SimpleFieldType,
			TextType:   schema.FieldType,
			Properties: map[string]datatype.DataTypable{},
			UUID:       schema.FieldUUID,
		}
	}
	// Build parent data type
	parentDataType := &datatype.DataType{
		Node:       report.StoredSchemas.Node,
		Name:       report.StoredSchemas.ParentSchema.Value.ObjectName,
		Type:       "",
		TextType:   "",
		Properties: childDataTypes,
		UUID:       report.StoredSchemas.ParentSchema.Value.ObjectUUID,
	}
	classifiedDatatypes := make(map[parser.NodeID]*datatype.ClassifiedDatatype, 0)

	parentClassificationRequest := classificationschema.ClassificationRequest{DetectorType: report.StoredSchemas.DetectorType, Value: parentDataType.ToClassificationRequestDetection(), Filename: report.StoredSchemas.ParentSchema.Source.Filename}
	parentClassification := report.Classifier.Schema.Classify(parentClassificationRequest)
	classifiedDatatypes[report.StoredSchemas.Node.ID()] = datatype.BuildClassifiedDatatype(parentDataType, parentClassification)

	// Export classified data types
	datatype.ExportClassified(report, detections.TypeSchemaClassified, report.StoredSchemas.DetectorType, idGenerator, classifiedDatatypes, report.StoredSchemas.ParentSchema.Parent)

	// Clear the map of stored schema detection information
	report.StoredSchemas = nil
}

func (report *Detectors) AddSecretLeak(
	secret secret.Secret,
	source source.Source,
) {
	report.AddDetection(detections.TypeSecretleak, detectors.DetectorGitleaks, source, secret)
}

func (report *Detectors) AddDetection(detectionType detections.DetectionType, detectorType detectors.Type, source source.Source, value interface{}) {
	data := &detections.Detection{
		Type:         detectionType,
		DetectorType: detectorType,
		Source:       source,
		Value:        value,
	}

	report.Add(data)
}

func (report *Detectors) AddDependency(
	detectorType detectors.Type,
	detectorLanguage detectors.Language,
	dependency dependencies.Dependency,
	source source.Source,
) {

	detection := &detections.Detection{
		DetectorType:     detectorType,
		DetectorLanguage: detectorLanguage,
		Value:            dependency,
		Source:           source,
		Type:             detections.TypeDependency,
	}
	classifiedDetection, err := report.Classifier.Dependencies.Classify(*detection)
	if err != nil {
		report.AddError(source.Filename, fmt.Errorf("classification dependencies error: %s", err))
		return
	}

	classifiedDetection.Type = detections.TypeDependencyClassified
	report.Add(classifiedDetection)
}

func (report *Detectors) AddFramework(
	detectorType detectors.Type,
	frameworkType frameworks.Type,
	data interface{},
	source source.Source,
) {
	detection := &detections.Detection{DetectorType: detectorType, Value: data, Source: source, Type: detections.TypeFramework}
	classifiedDetection, err := report.Classifier.Frameworks.Classify(*detection)
	if err != nil {
		report.AddError(source.Filename, fmt.Errorf("classification frameworks error: %s", err))
		return
	}

	classifiedDetection.Type = detections.TypeFrameworkClassified
	report.Add(classifiedDetection)
}

func (report *Detectors) AddError(filePath string, err error) {
	report.Add(&detections.ErrorDetection{
		Type:    detections.TypeError,
		Message: err.Error(),
		File:    filePath,
	})
}

func (report *Detectors) AddOperation(
	detectorType detectors.Type,
	operation operations.Operation,
	source source.Source,
) {
	report.AddDetection(detections.TypeOperation, detectorType, source, operation)
}

func (report *Detectors) Add(data interface{}) {
	detectionsToAdd := []interface{}{data}

	err := jsonlines.Encode(report.File, &detectionsToAdd)
	if err != nil {
		log.Printf("failed to encode data line %s", err)
	}
}

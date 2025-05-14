package report

import (
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
)

type Report interface {
	detections.ReportDetection
	schema.ReportSchema
	datatype.ReportDataType
	AddInterface(detectorType detectors.Type, data interfaces.Interface, source source.Source)
	AddFramework(detectorType detectors.Type, frameworkType frameworks.Type, data interface{}, source source.Source)
	AddDependency(detectorType detectors.Type, detectorLanguage detectors.Language, dependency dependencies.Dependency, source source.Source)
	AddSecretLeak(secret secret.Secret, source source.Source)
	AddOperation(detectorType detectors.Type, operation operations.Operation, source source.Source)
	AddError(filePath string, err error)
}

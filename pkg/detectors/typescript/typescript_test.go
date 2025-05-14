package typescript_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/moonkit02/dearer/pkg/detectors/typescript"
	detectortypes "github.com/moonkit02/dearer/pkg/report/detectors"

	"github.com/moonkit02/dearer/pkg/detectors"
	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorJavascript

func TestDetectorReportGeneral(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: typescript.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "general"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportDatatype(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: typescript.New(&nodeid.IntGenerator{Counter: 0})}}

	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "datatype"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportKnex(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: typescript.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "datatype_knex"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Frameworks)
}

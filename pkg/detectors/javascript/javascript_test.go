package javascript_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/moonkit02/dearer/pkg/detectors/javascript"
	detectortypes "github.com/moonkit02/dearer/pkg/report/detectors"

	"github.com/moonkit02/dearer/pkg/detectors"
	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorJavascript

func TestDetectorReportGeneral(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: javascript.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "general"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func TestDetectorReportDatatypes(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: javascript.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "datatypes"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

func BenchmarkDatatypes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var registrations = []detectors.InitializedDetector{{
		Type:     detectorType,
		Detector: javascript.New(&nodeid.IntGenerator{Counter: 0})}}
	for i := 0; i < b.N; i++ {
		testhelper.Extract(b, filepath.Join("testdata", "datatypes_performance"), registrations, detectorType)
	}
	b.StopTimer()
}

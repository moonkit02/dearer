package ipynb_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/moonkit02/dearer/pkg/detectors/ipynb"
	detectortypes "github.com/moonkit02/dearer/pkg/report/detectors"

	"github.com/moonkit02/dearer/pkg/detectors"
	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
)

const detectorType = detectortypes.DetectorIPYNB

func TestDetectorReportInterfaces(t *testing.T) {
	var registrations = []detectors.InitializedDetector{{
		Type:     detectortypes.DetectorIPYNB,
		Detector: ipynb.New(&nodeid.IntGenerator{Counter: 0})}}
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "notebooks"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

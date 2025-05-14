package envfile_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"
	"github.com/moonkit02/dearer/pkg/report/detectors"
)

const detectorType = detectors.DetectorEnvFile

var registrations = testhelper.RegistrationFor(detectorType)

func TestDetectorReportVariables(t *testing.T) {
	detectorReport := testhelper.Extract(t, filepath.Join("testdata", "variables"), registrations, detectorType)

	cupaloy.SnapshotT(t, detectorReport.Detections)
}

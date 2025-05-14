package composerjson_test

import (
	"path/filepath"
	"testing"

	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"
	"github.com/moonkit02/dearer/pkg/report/detectors"
	"github.com/bradleyjkemp/cupaloy"
)

const detectorType = detectors.DetectorDependencies

var registrations = testhelper.RegistrationFor(detectorType)

func TestDependenciesReport(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata"), registrations, detectorType)
	cupaloy.SnapshotT(t, report.Dependencies)
}

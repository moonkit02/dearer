package beego_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"
	reportdetectors "github.com/moonkit02/dearer/pkg/report/detectors"
)

const detectorType = reportdetectors.DetectorBeego

var registrations = testhelper.RegistrationFor(detectorType)

func TestDetectorReportDatabases(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "beego"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}

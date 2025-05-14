package dotnet_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"

	reportdetectors "github.com/moonkit02/dearer/pkg/report/detectors"
)

const detectorType = reportdetectors.DetectorDotnet

var registrations = testhelper.RegistrationFor(detectorType)

func TestDetectorReportDbContexts(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "project", "db_contexts", "multiple"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}

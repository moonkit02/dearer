package rails_test

import (
	"path/filepath"
	"testing"

	"github.com/bradleyjkemp/cupaloy"

	"github.com/moonkit02/dearer/pkg/detectors"
	detectortypes "github.com/moonkit02/dearer/pkg/report/detectors"

	"github.com/moonkit02/dearer/pkg/detectors/internal/testhelper"
	"github.com/moonkit02/dearer/pkg/detectors/rails"
	"github.com/moonkit02/dearer/pkg/parser/nodeid"
)

var detectorType = detectortypes.DetectorRails
var (
	registrations = []detectors.InitializedDetector{{Type: detectorType, Detector: rails.New(&nodeid.IntGenerator{Counter: 0})}}
)

func TestBuildReportSingleDatabase(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "database", "single"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}

func TestBuildReportMultipleDatabases(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "database", "multiple"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}

func TestBuildReportStorageProviders(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "storage"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}

func TestBuildReportCaches(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "cache"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Frameworks)
}

func TestBuildReportDatabaseSchema(t *testing.T) {
	report := testhelper.Extract(t, filepath.Join("testdata", "schema"), registrations, detectorType)

	cupaloy.SnapshotT(t, report.Detections)
}

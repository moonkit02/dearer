package detectors_test

import (
	"testing"

	"github.com/moonkit02/dearer/pkg/languages/php"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/testhelper"
)

func TestPHPObjects(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/class.php")
	runTest(t, "object_no_class", "object", "testdata/no_class.php")
}

func TestPHPString(t *testing.T) {
	runTest(t, "string", "string", "testdata/string.php")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, php.Get(), detectorType, fileName)
}

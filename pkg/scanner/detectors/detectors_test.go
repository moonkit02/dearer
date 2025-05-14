package detectors_test

import (
	"testing"

	"github.com/moonkit02/dearer/pkg/languages/ruby"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/testhelper"
)

func TestDatatypeDetector(t *testing.T) {
	runTest(t, "datatype", "datatype", "testdata/datatype.rb")
}

func TestDatatypeDetectorInvalidDetection(t *testing.T) {
	runTest(t, "datatype", "datatype", "testdata/invalid_datatype.java")
}

func TestInsecureUrlDetector(t *testing.T) {
	runTest(t, "insecure_url", "insecure_url", "testdata/insecureurl.rb")
}

func runTest(t *testing.T, name string, detectorType, fileName string) {
	testhelper.RunTest(t, name, ruby.Get(), detectorType, fileName)
}

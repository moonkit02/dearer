package detectors_test

import (
	"testing"

	"github.com/moonkit02/dearer/pkg/languages/python"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/testhelper"
)

func TestPythonObjects(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/class.py")
	runTest(t, "object_no_class", "object", "testdata/no_class.py")
}

func TestPythonString(t *testing.T) {
	runTest(t, "string", "string", "testdata/string.py")
	runTest(t, "string_literal", "string", "testdata/string_literal.py")
}

func runTest(t *testing.T, name, detectorType, fileName string) {
	testhelper.RunTest(t, name, python.Get(), detectorType, fileName)
}

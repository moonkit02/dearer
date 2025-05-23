package detectors_test

import (
	"testing"

	"github.com/moonkit02/dearer/pkg/languages/ruby"
	"github.com/moonkit02/dearer/pkg/scanner/detectors/testhelper"
)

func TestRubyObjectDetector(t *testing.T) {
	runTest(t, "object_class", "object", "testdata/object_class.rb")
	runTest(t, "object_hash", "object", "testdata/object_hash.rb")
	runTest(t, "object_projection", "object", "testdata/object_projection.rb")
}

func TestRubyStringDetector(t *testing.T) {
	runTest(t, "string_assign_eq", "string", "testdata/string_assign_eq.rb")
	runTest(t, "string_literal", "string", "testdata/string_literal.rb")
	runTest(t, "string_non_literal", "string", "testdata/string_non_literal.rb")
}

func runTest(t *testing.T, name string, detectorType, fileName string) {
	testhelper.RunTest(t, name, ruby.Get(), detectorType, fileName)
}

package rules_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/moonkit02/dearer/e2e/internal/testhelper"
)

func buildRulesTestCase(testName, path, ruleID string) testhelper.TestCase {
	arguments := []string{
		"scan",
		path,
		"--only-rule=" + ruleID,
		"--format=yaml",
		"--disable-default-rules",
		"--exit-code=0",
		"--external-rule-dir=" + filepath.Join("e2e", "rules", "testdata", "rules"),
	}

	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(testName, arguments, options)
}

func buildRulesTestCaseJsonV2(testName, path, ruleID string) testhelper.TestCase {
	arguments := []string{
		"scan",
		path,
		"--only-rule=" + ruleID,
		"--format=jsonv2",
		"--disable-default-rules",
		"--exit-code=0",
		"--external-rule-dir=" + filepath.Join("e2e", "rules", "testdata", "rules"),
	}

	options := testhelper.TestCaseOptions{}

	return testhelper.NewTestCase(testName, arguments, options)
}

func runRulesTest(folderPath string, ruleID string, t *testing.T) {
	snapshotDirectory := ".snapshots"

	testDataDir := fmt.Sprintf("testdata/data/%s", folderPath)

	testCases := []testhelper.TestCase{}
	testCases = append(testCases,
		buildRulesTestCase(
			testDataDir,
			filepath.Join("e2e", "rules", testDataDir),
			ruleID,
		),
	)

	testhelper.RunTestsWithSnapshotSubdirectory(t, testCases, snapshotDirectory)
}

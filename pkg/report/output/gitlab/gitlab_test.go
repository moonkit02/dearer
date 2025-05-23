package gitlab

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"

	securitytypes "github.com/moonkit02/dearer/pkg/report/output/security/types"
	util "github.com/moonkit02/dearer/pkg/util/output"
)

func TestJuiceShopSarif(t *testing.T) {
	securityOutput, err := os.ReadFile("testdata/juice-shop-security-report.json")
	if err != nil {
		t.Fatalf("failed to read file, err: %s", err)
	}

	var securityFindings map[string][]securitytypes.Finding
	err = json.Unmarshal(securityOutput, &securityFindings)
	if err != nil {
		t.Fatalf("couldn't unmarshal file output: %s", err)
	}

	startTime, _ := time.Parse("2006-01-02T15:04:05", "2006-01-02T15:04:05")
	endTime, _ := time.Parse("2006-01-02T15:04:05", "2006-01-02T15:05:05")

	res, err := ReportGitLab(securityFindings, startTime, endTime)
	if err != nil {
		t.Fatalf("failed to generate security output, err: %s", err)
	}

	output, err := util.ReportJSON(res)
	if err != nil {
		t.Fatalf("failed to generate JSON output, err: %s", err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(output), "", "\t")
	if err != nil {
		t.Fatalf("error indenting output, err: %s", err)
	}
	cupaloy.SnapshotT(t, prettyJSON.String())
}

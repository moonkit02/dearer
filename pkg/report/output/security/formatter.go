package security

import (
	"fmt"
	"time"

	"github.com/hhatto/gocloc"

	"github.com/moonkit02/dearer/cmd/bearer/build"
	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/engine"
	"github.com/moonkit02/dearer/pkg/flag"
	dataflowtypes "github.com/moonkit02/dearer/pkg/report/output/dataflow/types"
	"github.com/moonkit02/dearer/pkg/report/output/gitlab"
	"github.com/moonkit02/dearer/pkg/report/output/html"
	"github.com/moonkit02/dearer/pkg/report/output/reviewdog"
	"github.com/moonkit02/dearer/pkg/report/output/sarif"
	outputtypes "github.com/moonkit02/dearer/pkg/report/output/types"
	outputhandler "github.com/moonkit02/dearer/pkg/util/output"
)

type Formatter struct {
	ReportData   *outputtypes.ReportData
	Config       settings.Config
	engine       engine.Engine
	GoclocResult *gocloc.Result
	StartTime    time.Time
	EndTime      time.Time
}

type JsonV2Output struct {
	Source   string                `json:"source" yaml:"source"`
	Version  string                `json:"version" yaml:"version"`
	Findings RawFindings           `json:"findings" yaml:"findings"`
	Expected ExpectedDetections    `json:"expected_findings,omitempty" yaml:"expected_findings,omitempty"`
	Errors   []dataflowtypes.Error `json:"errors" yaml:"errors"`
}

func NewFormatter(
	reportData *outputtypes.ReportData,
	config settings.Config,
	engine engine.Engine,
	goclocResult *gocloc.Result,
	startTime time.Time,
	endTime time.Time,
) *Formatter {
	return &Formatter{
		ReportData:   reportData,
		Config:       config,
		engine:       engine,
		GoclocResult: goclocResult,
		StartTime:    startTime,
		EndTime:      endTime,
	}
}

func (f Formatter) Format(format string) (output string, err error) {
	switch format {
	case flag.FormatEmpty:
		output = BuildReportString(f.ReportData, f.Config, f.engine, f.GoclocResult).String()
	case flag.FormatSarif:
		sarifContent, sarifErr := sarif.ReportSarif(f.ReportData.FindingsBySeverity, f.Config.Rules)
		if sarifErr != nil {
			return output, fmt.Errorf("error generating sarif report %s", sarifErr)
		}
		return outputhandler.ReportJSON(sarifContent)
	case flag.FormatReviewDog:
		sastContent, reviewdogErr := reviewdog.ReportReviewdog(f.ReportData.FindingsBySeverity)
		if reviewdogErr != nil {
			return output, fmt.Errorf("error generating reviewdog report %s", reviewdogErr)
		}
		return outputhandler.ReportJSON(sastContent)
	case flag.FormatGitLabSast:
		sastContent, sastErr := gitlab.ReportGitLab(f.ReportData.FindingsBySeverity, f.StartTime, f.EndTime)
		if sastErr != nil {
			return output, fmt.Errorf("error generating gitlab-sast report %s", sastErr)
		}
		return outputhandler.ReportJSON(sastContent)
	case flag.FormatJSON:
		return outputhandler.ReportJSON(f.ReportData.FindingsBySeverity)
	case flag.FormatJSONV2:
		return outputhandler.ReportJSON(JsonV2Output{
			Source:   "Bearer",
			Version:  build.Version,
			Findings: f.ReportData.RawFindings,
			Expected: f.ReportData.ExpectedDetections,
			Errors:   f.ReportData.Dataflow.Errors,
		})
	case flag.FormatYAML:
		return outputhandler.ReportYAML(f.ReportData.FindingsBySeverity)
	case flag.FormatHTML:
		title := "Security Report"
		body, securityErr := html.ReportSecurityHTML(f.ReportData.FindingsBySeverity)
		if securityErr != nil {
			return output, securityErr
		}

		output, err = html.ReportHTMLWrapper(title, body)
		if err != nil {
			err = fmt.Errorf("could not generate html page %s", err)
		}
	}

	return output, err
}

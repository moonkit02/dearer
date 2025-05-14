package privacy

import (
	"fmt"

	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/flag"
	"github.com/moonkit02/dearer/pkg/report/output/html"
	outputtypes "github.com/moonkit02/dearer/pkg/report/output/types"
	outputhandler "github.com/moonkit02/dearer/pkg/util/output"
)

type Formatter struct {
	ReportData *outputtypes.ReportData
	Config     settings.Config
}

func NewFormatter(reportData *outputtypes.ReportData, config settings.Config) *Formatter {
	return &Formatter{
		ReportData: reportData,
		Config:     config,
	}
}

func (f Formatter) Format(format string) (output string, err error) {
	switch format {
	case flag.FormatEmpty, flag.FormatCSV:
		stringBuilder, err := BuildCsvString(f.ReportData, f.Config)
		if err != nil {
			return output, err
		}
		output = stringBuilder.String()
	case flag.FormatJSON:
		return outputhandler.ReportJSON(f.ReportData.PrivacyReport)
	case flag.FormatYAML:
		return outputhandler.ReportYAML(f.ReportData.PrivacyReport)
	case flag.FormatHTML:
		title := "Privacy Report"
		body, err := html.ReportPrivacyHTML(f.ReportData.PrivacyReport)
		if err != nil {
			return output, err
		}

		output, err = html.ReportHTMLWrapper(title, body)
		if err != nil {
			return output, fmt.Errorf("could not generate html page %s", err)
		}
	}

	return output, err
}

package detectors

import (
	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/flag"
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
	case flag.FormatEmpty, flag.FormatJSON:
		return outputhandler.ReportJSON(f.ReportData.Detectors)
	case flag.FormatYAML:
		return outputhandler.ReportYAML(f.ReportData.Detectors)
	}

	return output, err
}

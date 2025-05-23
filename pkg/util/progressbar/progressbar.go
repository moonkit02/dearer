package progressbar

import (
	"github.com/moonkit02/dearer/pkg/commands/process/settings"
	"github.com/moonkit02/dearer/pkg/util/output"
	"github.com/schollz/progressbar/v3"
)

func GetProgressBar(filesLength int, config *settings.Config) *progressbar.ProgressBar {
	hideProgress := config.Scan.HideProgressBar || config.Scan.Quiet || config.Debug
	return progressbar.NewOptions(filesLength,
		progressbar.OptionSetVisibility(!hideProgress),
		progressbar.OptionSetWriter(output.ErrorWriter()),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(15),
		progressbar.OptionEnableColorCodes(false),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionOnCompletion(func() {
			output.ErrorWriter().Write([]byte("\n")) //nolint:all,errcheck
		}),
		progressbar.OptionSetDescription(" └"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

package ui

import (
	"github.com/schollz/progressbar/v3"
)

// Returns a progress bar with default options
func GetProgressBar(length int, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions(length,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█", // "=", "█"
			SaucerHead:    "",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}))
}

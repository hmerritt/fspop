package ui

import (
	"bufio"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-wordwrap"
)

// Maximum width of any line, in character count.
const MaxLineLength = 75

// Extend cli.ColoredUi struct and interface
type Ui struct {
	*cli.ColoredUi
	SuccessColor cli.UiColor
}

func GetUi() *Ui {
	return &Ui{
		&cli.ColoredUi{
			InfoColor:  cli.UiColorCyan,
			ErrorColor: cli.UiColorRed,
			WarnColor:  cli.UiColorYellow,
			Ui: &cli.BasicUi{
				Reader:      bufio.NewReader(os.Stdin),
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
			},
		},
		cli.UiColorGreen,
	}
}

// Outputs green text
func (u *Ui) Success(message string) {
	u.Ui.Output(u.Colorize(message, cli.UiColorGreen))
}

// Add color to a string
//
// Imported directly from github.com/mitchellh/cli/blob/v1.1.2/ui_colored.go#L62
func (u *Ui) Colorize(message string, uc cli.UiColor) string {
	const noColor = -1

	if uc.Code == noColor {
		return message
	}

	attr := []color.Attribute{color.Attribute(uc.Code)}
	if uc.Bold {
		attr = append(attr, color.Bold)
	}

	return color.New(attr...).SprintFunc()(message)
}

// Wraps the given text to maxLineLength.
func WrapAtLength(s string) string {
	return wordwrap.WrapString(s, MaxLineLength)
}

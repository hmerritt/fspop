package command

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/cli"
	"gitlab.com/merrittcorp/fspop/version"
)

func Run() {
	// Initiate new CLI app
	app := cli.NewCLI("fspop", version.GetVersion().VersionNumber())
	app.Args = os.Args[1:]

	getBaseCommand := func() *BaseCommand {
		return &BaseCommand{
			UI: &cliUi{
				&cli.ColoredUi{
					ErrorColor: cli.UiColorRed,
					WarnColor:  cli.UiColorYellow,
					Ui: &cli.BasicUi{
						Reader:      bufio.NewReader(os.Stdin),
						Writer:      os.Stdout,
						ErrorWriter: os.Stderr,
					},
				},
				cli.UiColorGreen,
			},
		}
	}

	// Feed active commands to CLI app
	app.Commands = map[string]cli.CommandFactory{
		"deploy": func() (cli.Command, error) {
			return &DeployCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"display": func() (cli.Command, error) {
			return &DisplayCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"init": func() (cli.Command, error) {
			return &InitCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &ListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
	}

	// Run app
	exitStatus, err := app.Run()
	if err != nil {
		os.Stderr.WriteString(fmt.Sprint(err))
	}

	os.Exit(exitStatus)
}

// Master command type which in present in all commands
//
// Used to standardize UI output
type BaseCommand struct {
	UI *cliUi
}

// Extend cli.Ui interface by adding a 'Success' method,
// this method is used for green output.
type cliUi struct {
	*cli.ColoredUi
	SuccessColor cli.UiColor
}

func (u *cliUi) Success(message string) {
	u.Ui.Output(u.Colorize(message, cli.UiColorGreen))
}

func (u *cliUi) Colorize(message string, uc cli.UiColor) string {
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

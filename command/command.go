package command

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"gitlab.com/merrittcorp/fspop/version"
)

func Run() {
	// Initiate new CLI app
	app := cli.NewCLI("fspop", version.GetVersion().VersionNumber())
	app.Args = os.Args[1:]

	getBaseCommand := func() *BaseCommand {
		return &BaseCommand{
			UI: &cli.ColoredUi{
				ErrorColor: cli.UiColorRed,
				WarnColor:  cli.UiColorYellow,
				Ui: &cli.BasicUi{
					Reader:      bufio.NewReader(os.Stdin),
					Writer:      os.Stdout,
					ErrorWriter: os.Stderr,
				},
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
	UI cli.Ui
}

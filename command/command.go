package command

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"gitlab.com/merrittcorp/fspop/ui"
	"gitlab.com/merrittcorp/fspop/version"
)

func Run() {
	// Initiate new CLI app
	app := cli.NewCLI("fspop", version.GetVersion().VersionNumber())
	app.Args = os.Args[1:]

	getBaseCommand := func() *BaseCommand {
		return &BaseCommand{
			UI: ui.GetUi(),
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

	// Exit without an error if no arguments were passed
	if len(app.Args) == 0 {
		os.Exit(0)
	}

	os.Exit(exitStatus)
}

// Master command type which in present in all commands
//
// Used to standardize UI output
type BaseCommand struct {
	UI *ui.Ui
}

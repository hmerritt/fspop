package command

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"gitlab.com/merrittcorp/fspop/version"
)

func Run() {
	// Initiate new CLI app
	app := cli.NewCLI("fspop", version.GetVersion().VersionNumber())
	app.Args = os.Args[1:]

	// Feed active commands to CLI app
	app.Commands = map[string]cli.CommandFactory{
		"list": func() (cli.Command, error) {
			return &ListCommand{}, nil
		},
	}

	// Run app
	exitStatus, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

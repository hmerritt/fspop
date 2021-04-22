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

	// Run app
	exitStatus, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

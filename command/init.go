package command

import (
	"fmt"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
)

type InitCommand struct{}

func (c *InitCommand) Synopsis() string {
	return "Create a new structure file"
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: fspop init [options] NAME
  
  Init will create a new structure file.

  Structure files can be deployed using the command:

  $ fspop deploy <NAME>

`

	return strings.TrimSpace(helpText)
}

func (c *InitCommand) Run(args []string) int {
	path := "structure.yml"

	if len(args) > 0 {
		path = args[0]
	}

	fmt.Println(message.Green("Success.") + ` Created '` + path + `' structure file.

Structure files can be deployed using the command:

$ fspop deploy ` + path)

	return 0
}

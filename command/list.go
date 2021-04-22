package command

import (
	"fmt"
	"strings"
)

type ListCommand struct{}

func (c *ListCommand) Synopsis() string {
	return "List (potential) structure files in the current directory"
}

func (c *ListCommand) Help() string {
	helpText := `
Usage: fspop list [options] PATH
  
  Lists all YAML files in the current directory.
  
  YAML are identified by their file extension(s):
  - .yml
  - .yaml

  Thease files (could) be fspop structure files
  that can be deployed using the command:

  $ fspop deploy [name-of-file]

`

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Run(args []string) int {
	fmt.Println("LIST")
	return 0
}

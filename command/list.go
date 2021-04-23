package command

import (
	"fmt"
	"log"
	"os"
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
	path := "./"

	if len(args) > 0 {
		path = args[0]
	}

	// TODO: Check if path exists
	//       print warning and fallback to './' if not

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	return 0
}

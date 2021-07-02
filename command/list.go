package command

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"gitlab.com/merrittcorp/fspop/parse"
)

type ListCommand struct {
	*BaseCommand
}

func (c *ListCommand) Synopsis() string {
	return "List structure files in the current directory"
}

func (c *ListCommand) Help() string {
	helpText := `
Usage: fspop list [options] PATH
  
  Lists all YAML files in the current directory.
  
  YAML files are identified by the file extension(s):
  - .yml
  - .yaml

`

	return strings.TrimSpace(helpText)
}

func (c *ListCommand) Run(args []string) int {
	path := "./"

	if len(args) > 0 {
		path = args[0]
	}

	files, err := os.ReadDir(path)
	if err != nil {
		c.UI.Error("Unable to read directory files")
		c.UI.Error(fmt.Sprint(err))
		c.UI.Warn("\nThis is most likely due to a lack of permissions,")
		c.UI.Warn("check you have (at least) read access to this directory.")
		return 1
	}

	yamlFiles := make([]string, 0)

	for _, f := range files {
		name := f.Name()
		if !f.IsDir() {
			if parse.IsYamlFile(name) {
				yamlFiles = append(yamlFiles, name)
			}
		}
	}

	sort.Strings(yamlFiles)

	if len(yamlFiles) == 0 {
		// Get friendly path name.
		// Converts './' to 'current'
		pathFriendly := "'" + path + "'"
		if path == "./" {
			pathFriendly = "current"
		}

		c.UI.Output("No structure files found in the " + pathFriendly + " directory.\n")
		c.UI.Output("Create a structure file using:")
		c.UI.Output("$ fspop init <NAME>")
		return 2
	}

	fmt.Printf("Found %d possible structure files:\n", len(yamlFiles))

	for _, yf := range yamlFiles {
		fmt.Printf("-- %s\n", yf)
	}

	c.UI.Output("\nDeploy a structure file using:")
	c.UI.Output("$ fspop deploy <NAME>")

	return 0
}

package command

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
	"gitlab.com/merrittcorp/fspop/yaml"
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
		path = yaml.AddYamlExtension(args[0])
	}

	// TODO: Add actual error reporting

	// Create structure file
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// Write content into file
	_, err = file.WriteString(yamlFileContent())
	if err != nil {
		fmt.Println(err)
		file.Close()
		return 1
	}

	// Close file
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println(message.Green("Success.") + ` Created '` + path + `' structure file.

Structure files can be deployed using the command:

$ fspop deploy ` + path)

	return 0
}

func yamlFileContent() string {
	return `version: 4

name: fspop-example

data:
  - example: text can be imported like this
  - data_file: /path/to/file

dynamic:
  - dyn:
    - amount: 100
    - data: example
    - type: file
    - name: fspop_example_$num
    - padded: true

structure:
  - file.empty
  - file_data: example
  - $dyn
  - /folder:
    - /sub-folder
`
}

package command

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
	"gitlab.com/merrittcorp/fspop/parse"
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
	path := parse.DefaultYamlFile

	if len(args) > 0 {
		path = parse.AddYamlExtension(args[0])
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
	return `###########################
## fspop structure file
##
## Usage info:
## $ fspop -h
###########################
version: 4

name: fspop-structure

data:
  - example: text can be imported like this
  - data_file: /path/to/file
  - data_url: https://example.com/data/from/url
  - data_actual: https://via.placeholder.com/400/771796

dynamic:
  - dyn:
    - amount: 10
    - data: example
    - type: file
    - name: fspop_$num.dynamic
    - padded: true
    - start: 95

  - dyn_folders:
    - amount: 10
    - type: folder
    - name: fspop_$num_dynamic
    - padded: false
    - start: 5

structure:
  - file.empty
  - file.data: example
  - image.png: data_actual
  - /dynamic-items:
    - /dynamic-folders:
      - $dyn_folders
    - $dyn
`
}

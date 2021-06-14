package command

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/merrittcorp/fspop/parse"
)

type InitCommand struct {
	*BaseCommand
}

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

	// Check if file already exists
	if parse.FileExists(path) {
		// Exit. Dont't overwrite existing files
		c.UI.Error("Structure file '" + path + "' already exists.\n")
		c.UI.Warn("Rename or remove the existing file and try again.")
		return 1
	}

	// Create structure file
	file, err := os.Create(path)
	if err != nil {
		c.UI.Error("Unable to create new structure file '" + path + "'.")
		c.UI.Error(fmt.Sprint(err))
		c.UI.Warn("\nThis is most likely due to a lack of permissions,")
		c.UI.Warn("check you have write access to this directory.")
		return 1
	}

	// Write content into file
	_, err = file.WriteString(yamlFileContent())
	if err != nil {
		file.Close()

		c.UI.Error("Created structure file '" + path + "', but failed to write data into it.")
		c.UI.Error(fmt.Sprint(err))
		c.UI.Warn("\nThis is most likely due to a lack of permissions,")
		c.UI.Warn("check you have write access to this file and directory.")
		return 1
	}

	// Close file
	err = file.Close()
	if err != nil {
		c.UI.Error("Created structure file '" + path + "', but somthing went wrong when finalising it.")
		c.UI.Error(fmt.Sprint(err))
		return 1
	}

	c.UI.Output(c.UI.Colorize("Success.", c.UI.SuccessColor) + " Created '" + path + "' structure file.\n")
	c.UI.Output("Deploy this structure file using the command:")
	c.UI.Output("$ fspop deploy " + path + "")

	return 0
}

func yamlFileContent() string {
	return `###########################
## fspop structure file
##
## Usage info:
## $ fspop -h
##
## my-folder/          - folders need a '/' at the end
## my-file             - items are assumed to be files (unless there is a '/')
## data-file: data_var - files can have custum data by assigning a data variable 'file: your_data_variable'
## $dynamic_item       - dynamic items have a dollar prefix '$' + 'your_dynamic_variable'
###########################
version: 4

name: fspop-structure
entrypoint: fspop/structure

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

my-custom-sub-structure: &myStructureName
  - 1. this structure can be reused
  - 2. data can be imported: example
  - $dyn

structure:
  - file.empty
  - file.data: example
  - image.png: data_actual
  - sub-structure/:
    - *myStructureName
    - reuse/:
      - *myStructureName
  - dynamic-items/:
    - dynamic-folders/:
      - $dyn_folders
    - $dyn
`
}

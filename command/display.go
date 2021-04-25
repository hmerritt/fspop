package command

import (
	"fmt"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
	"gitlab.com/merrittcorp/fspop/parse"
)

type DisplayCommand struct{}

func (c *DisplayCommand) Synopsis() string {
	return "Print the directory tree of a structure file"
}

func (c *DisplayCommand) Help() string {
	helpText := `
Usage: fspop display [options] NAME
  
  Display prints the directory tree of a structure file to the terminal.

  Structure files can be deployed using the command:

  $ fspop deploy <NAME>

`

	return strings.TrimSpace(helpText)
}

func (c *DisplayCommand) Run(args []string) int {
	var path string

	if len(args) == 0 {
		message.Warn("No file entered.")
		message.Warn("Trying default '" + parse.DefaultYamlFile + "' instead.")
		message.Text("")
		path = parse.DefaultYamlFile
	} else {
		path = parse.ElasticExtension(args[0])
	}

	var fileData []byte
	var fileError error

	// Decide if URL or file
	if parse.UseUrl(path) {
		message.Spinner.Start("", " Fetching URL data...")

		fileData, fileError = parse.FetchUrl(path)

		message.Spinner.Stop()

		if fileError != nil {
			message.Error("Unable to fetch URL data.")
			message.Error(fmt.Sprint(fileError))
			fmt.Println()
			message.Warn("Make sure the link is accessible and try again.")
			return 2
		}
	} else {
		fileData, fileError = parse.FetchFile(path)

		if fileError != nil {
			message.Error("Unable to open file.")
			message.Error(fmt.Sprint(fileError))
			fmt.Println()
			message.Warn("Check the file is exists and try again.")
			return 2
		}
	}

	// Parse YAML
	structure, parseErr := parse.ParseYaml(fileData)
	// https://merritt.es/tools/structure.yml

	if parseErr != nil {
		message.Error("Unable to parse YAML file.")
		message.Error(fmt.Sprint(parseErr))
		fmt.Println()
		message.Warn("Check the file is valid YAML and try again.")
		return 2
	}

	fmt.Printf("Version: %s \n", structure.Version)
	fmt.Printf("Name: %s \n", structure.Name)
	fmt.Printf("Data: %v \n", structure.Data)
	fmt.Printf("Data: %v \n", structure.Data.([]interface{})[1].(map[interface{}]interface{})["data_file"])
	fmt.Printf("Dynamic: %v \n", structure.Dynamic)
	fmt.Printf("Structure: %v \n", structure.Structure)

	return 0
}

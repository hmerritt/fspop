package command

import (
	"fmt"
	"strings"

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
	path := parse.DefaultYamlFile

	if len(args) > 0 {
		path = parse.ElasticExtension(args[0])
	}

	var yamlData []byte

	// Decide if URL or file
	if parse.UseUrl(path) {
		fmt.Println("Fetching remote YAML file")
		yamlData = parse.FetchUrl(path)
	} else {
		yamlData = parse.FetchFile(path)
	}

	// TODO: catch fetch errors here

	// Parse YAML
	structure := parse.ParseYaml(yamlData)
	// https://merritt.es/tools/structure.yml

	// TODO: catch parse errors here

	fmt.Printf("Version: %s \n", structure.Version)
	fmt.Printf("Name: %s \n", structure.Name)
	fmt.Printf("Data: %v \n", structure.Data)
	fmt.Printf("Data: %v \n", structure.Data.([]interface{})[1].(map[interface{}]interface{})["data_file"])
	fmt.Printf("Dynamic: %v \n", structure.Dynamic)
	fmt.Printf("Structure: %v \n", structure.Structure)

	return 0
}

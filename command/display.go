package command

import (
	"fmt"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
	"gitlab.com/merrittcorp/fspop/parse"
	"gitlab.com/merrittcorp/fspop/structure"
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

	fsStructure := parse.ParseAndRefineYaml(fileData)
	//parse.ParseAndRefineYaml(fileData)

	fmt.Println(fsStructure.Items)

	fmt.Print("\n\n\n")

	fsStructure.Crawl(func(key string, item structure.FspopItem) {
		fmt.Printf("%v  %s  :  %s\n", (key == item.Path.ToString()), key, item.Path.ToString())
	})

	return 0
}

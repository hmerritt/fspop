package command

import (
	"strings"

	"github.com/disiqueira/gotree"
	"gitlab.com/merrittcorp/fspop/parse"
	"gitlab.com/merrittcorp/fspop/structure"
)

type DisplayCommand struct {
	*BaseCommand
}

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
		// Use default structure file
		// Checks for both '.yaml' and '.yml' extensions
		path = parse.AddYamlExtension(parse.ElasticExtension(parse.DefaultYamlFileName))
		c.UI.Warn("No file entered.")
		c.UI.Warn("Trying default '" + path + "' instead.\n")
	} else {
		path = parse.ElasticExtension(args[0])
	}

	// Fetch structure
	fsStructure := parse.FetchAndParseStructure(path)

	// Single-depth slice to store all nodes + their tree instance
	treeNodes := make(map[string]*gotree.Tree, len(fsStructure.Items))

	// Attempt to build tree structure from FspopStructure.Items
	// Initialise tree struct
	tree := gotree.New(fsStructure.Name)
	treeNodes["/"] = &tree

	// Populate treeNodes structure
	// Crawl each path in structure
	fsStructure.Crawl(func(item structure.FspopItem) {
		treeItem := treeNodes["/"]

		// Breakdown invividual path nodes
		for k, v := range item.Path.PathProgressive() {
			name := item.Path.Path[k]
			// Exists already
			if _, ok := treeNodes[v]; ok {
				treeItem = treeNodes[v]
			} else {
				newDir := (*treeItem).Add(name)
				treeNodes[v] = &newDir
				treeItem = treeNodes[v]
			}
		}
	})

	// Print tree
	c.UI.Output(tree.Print())

	return 0
}

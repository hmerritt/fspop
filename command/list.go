package command

import (
	"fmt"
	"log"
	"os"
	"regexp"
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
		// TODO: Print actual error message here
		log.Fatal(err)
	}

	// TODO: Run regex in a thread for each file
	//       wait at end of for-loop before resuming

	yamlFiles := make([]string, 0)

	for _, f := range files {
		match, _ := regexp.MatchString(".yml|.yaml", f.Name())
		if match && !f.IsDir() {
			yamlFiles = append(yamlFiles, f.Name())
		}
	}

	if len(yamlFiles) == 0 {
		fmt.Println("No structure files found in the current directory.")
		fmt.Println()
		fmt.Println("Create a structure file using:")
		fmt.Println("$ fspop init <name>")
		return 2
	}

	fmt.Printf("Found %d possible structure files:\n", len(yamlFiles))

	for _, yf := range yamlFiles {
		fmt.Printf("-- %s\n", yf)
	}

	fmt.Println("\nDeploy a structure file using:")
	fmt.Println("$ fspop deploy <name>")

	return 0
}

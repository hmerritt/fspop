package command

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
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

	files, err := os.ReadDir(path)
	if err != nil {
		// TODO: Print actual error message here
		log.Fatal(err)
	}

	start := time.Now()

	// TODO: Run regex in a thread for each file
	//       wait at end of for-loop before resuming
	var wg sync.WaitGroup

	yamlFiles := make([]string, 0)

	for _, f := range files {
		if !f.IsDir() {
			wg.Add(1)
			go matchYAMLFile(&wg, &yamlFiles, f.Name())
		}
	}

	wg.Wait()

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

	fmt.Println(time.Since(start))

	return 0
}

func matchYAMLFile(wg *sync.WaitGroup, yamlFiles *[]string, name string) {
	match, _ := regexp.MatchString(".yml|.yaml", name)
	if match {
		*yamlFiles = append(*yamlFiles, name)
	}
	wg.Done()
}

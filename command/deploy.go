package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.com/merrittcorp/fspop/message"
	"gitlab.com/merrittcorp/fspop/parse"
	"gitlab.com/merrittcorp/fspop/structure"
)

type DeployCommand struct{}

func (c *DeployCommand) Synopsis() string {
	return "Deploy (create) a structure"
}

func (c *DeployCommand) Help() string {
	helpText := `
Usage: fspop deploy [options] NAME
  
  Deploy creates the file-structure defined in a structure config-file.

  Structure files can be created using the command:

  $ fspop init <NAME>

`

	return strings.TrimSpace(helpText)
}

func (c *DeployCommand) Run(args []string) int {
	var path string

	if len(args) == 0 {
		// Use default structure file
		// Checks for both '.yaml' and '.yml' extensions
		path = parse.AddYamlExtension(parse.ElasticExtension(parse.DefaultYamlFileName))
		message.Warn("No file entered.")
		message.Warn("Trying default '" + path + "' instead.")
		message.Text("")
	} else {
		path = parse.ElasticExtension(args[0])
	}

	// Fetch structure
	fsStructure := parse.FetchAndParseStructure(path)

	timeStart := time.Now()

	// Print structure stats
	printStructureStats(fsStructure)

	// Initiate progress bar
	bar := message.GetProgressBar(len(fsStructure.Items), " Deploying")

	// Loop each item endpoint
	// An endpoint can be both a file or directory
	for key, item := range fsStructure.Items {
		// time.Sleep(200 * time.Millisecond)

		// Detect a dynamic item
		// Dynamic items are special and treated separately
		if len(item.DynamicKey) > 0 {
			// Check if item 'DynamicKey' actually exists
			if fsDynamicItem, ok := fsStructure.Dynamic[item.DynamicKey]; ok {
				// Create empty fileData
				fileData := make([]byte, 0)

				// If dynamic item has type 'file'
				// Load file data before creating items
				if !fsDynamicItem.IsTypeDirectory() {
					// Resolve 'DataKey' into a data item
					if fsDataItem, ok := fsStructure.GetDataItem(fsDynamicItem.DataKey); ok {
						fileData, _ = resolveDataPayload(fsDataItem.Data)
					}
				}

				// Loop n times
				// n = fsDynamicItem.Count = user defined
				fsDynamicItemMaxCount := fsDynamicItem.Start + fsDynamicItem.Count
				for i := fsDynamicItem.Start; i < fsDynamicItemMaxCount; i++ {
					itemPath, itemParentPath := fsDynamicItem.BuildItemPath(fsStructure.Name, &item.Path, i)

					// Is directory
					if fsDynamicItem.IsTypeDirectory() {
						// Create full directory
						os.MkdirAll(itemPath, os.ModePerm)

						// Is file
					} else {
						// Create parent directory
						os.MkdirAll(itemParentPath, os.ModePerm)

						// Create file
						// TODO: check and print any errors
						newFile, _ := parse.CreateFile(itemPath)

						if len(fileData) > 0 {
							newFile.Write(fileData)
						}

						newFile.Close()
					}
				}
			} //else {
			// 	// Dynamic key does not exist
			// 	fmt.Println("Dynamic key does not exist")
			// }

			bar.Add(1)
			continue
		}

		// Directory
		// Check if endpoint is a directory
		if structure.IsDirectory(key) {
			// Recursively make all directories
			// TODO: check and print any errors
			os.MkdirAll(fmt.Sprintf("%s/%s", fsStructure.Name, key), os.ModePerm)

		} else {
			// File
			// Recursively make all parent directories
			os.MkdirAll(fmt.Sprintf("%s/%s", fsStructure.Name, item.Path.ParentString()), os.ModePerm)

			// Create empty file
			// TODO: check and print any errors
			newFile, _ := parse.CreateFile(fmt.Sprintf("%s/%s", fsStructure.Name, item.Path.ToString()))

			// Resolve 'DataKey' into a data item
			if fsDataItem, ok := fsStructure.GetDataItem(item.DataKey); ok {
				fetchAndWriteToFile(newFile, fsDataItem.Data)
			}

			newFile.Close()
		}

		bar.Add(1)
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("%s in %s\n", message.Green("Structure deployed"), time.Since(timeStart))

	return 0
}

// Print basic structure stats
func printStructureStats(fsStructure *structure.FspopStructure) {
	message.Info("Structure File")
	message.Text("├── Data Variables       " + fmt.Sprint(len(fsStructure.Data)))
	message.Text("├── Dynamic Variables    " + fmt.Sprint(len(fsStructure.Dynamic)))
	message.Text("└── Structure Endpoints  " + fmt.Sprint(len(fsStructure.Items)))
	message.Text("")
}

// Fetch data key payload AND write data to an open file
func fetchAndWriteToFile(file *os.File, dataString string) error {
	// Get data
	data, errData := resolveDataPayload(dataString)

	if errData != nil {
		return errData
	}

	// Write to file
	_, errWrite := file.Write(data)

	if errWrite != nil {
		return errWrite
	}

	return nil
}

// Fetch data key payload from difference sources
// Local file, URL, as-is (text)
func resolveDataPayload(dataString string) ([]byte, error) {
	// Is URL
	if parse.UseUrl(dataString) {
		// Fetch URL data
		dataFromURL, err := parse.FetchUrl(dataString)
		if err == nil {
			return dataFromURL, nil
		} else {
			return nil, err
		}
	} else if parse.FileExists(dataString) {
		// Load file data and place into new file
		dataFromFile, err := parse.FetchFile(dataString)
		if err == nil {
			return dataFromFile, nil
		} else {
			return nil, err
		}
	}

	// Use data as-is (text)
	return []byte(dataString), nil
}

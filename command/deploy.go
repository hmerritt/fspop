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

	// Fetch structure file
	fileData := parse.FetchYaml(path)

	timeStart := time.Now()

	fsStructure := parse.ParseAndRefineYaml(fileData)

	// Print structure stats
	message.Info("Structure File")
	message.Text("├── Data Variables       " + fmt.Sprint(len(fsStructure.Data)))
	message.Text("├── Dynamic Variables    " + fmt.Sprint(len(fsStructure.Dynamic)))
	message.Text("└── Structure Endpoints  " + fmt.Sprint(len(fsStructure.Items)))
	message.Text("")

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
				if !strings.HasPrefix(fsDynamicItem.Type, "dir") || !strings.HasPrefix(fsDynamicItem.Type, "fol") {
					// Check if file has a data key
					if len(fsDynamicItem.DataKey) > 0 {
						// Check if item 'DataKey' actually exists
						if fsDataItem, ok := fsStructure.Data[fsDynamicItem.DataKey]; ok {
							// Is URL
							if parse.UseUrl(fsDataItem.Data) {
								// Fetch URL data
								dataFromURL, _ := parse.FetchUrl(fsDataItem.Data)
								fileData = dataFromURL

							} else if parse.FileExists(fsDataItem.Data) {
								// Load file data and place into new file
								dataFromFile, _ := parse.FetchFile(fsDataItem.Data)
								fileData = dataFromFile

							} else {
								// Treat data as plain text
								fileData = []byte(fsDataItem.Data)
							}
						}
					}
				}

				// Loop n times
				// n = fsDynamicItem.Count = user defined
				fsDynamicItemMaxCount := fsDynamicItem.Start + fsDynamicItem.Count
				for i := fsDynamicItem.Start; i < fsDynamicItemMaxCount; i++ {
					// Build file/directory name
					itemCountPadWidth := 1
					if fsDynamicItem.Padded {
						itemCountPadWidth = len(fmt.Sprint(fsDynamicItemMaxCount))
					}
					itemCountPadded := fmt.Sprintf("%0*d", itemCountPadWidth, i)
					if !strings.Contains(fsDynamicItem.Name, "$num") {
						fsDynamicItem.Name = fmt.Sprintf("%s_$num", fsDynamicItem.Name)
					}
					itemName := strings.ReplaceAll(fsDynamicItem.Name, "$num", itemCountPadded)
					itemParentPath := fmt.Sprintf("%s/%s/", fsStructure.Name, item.Path.ParentString())
					itemPath := fmt.Sprintf("%s/%s", itemParentPath, itemName)

					// Is directory
					if strings.HasPrefix(fsDynamicItem.Type, "dir") || strings.HasPrefix(fsDynamicItem.Type, "fol") {
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

			bar.Add(1)
			continue
		}

		// File
		// Recursively make all parent directories
		os.MkdirAll(fmt.Sprintf("%s/%s", fsStructure.Name, item.Path.ParentString()), os.ModePerm)

		// Create empty file
		// TODO: check and print any errors
		newFile, _ := parse.CreateFile(fmt.Sprintf("%s/%s", fsStructure.Name, item.Path.ToString()))

		// Check if file has a data key
		if len(item.DataKey) > 0 {
			// Check if item 'DataKey' actually exists
			if fsDataItem, ok := fsStructure.Data[item.DataKey]; ok {
				// Is URL
				if parse.UseUrl(fsDataItem.Data) {
					// Fetch URL data
					dataFromURL, _ := parse.FetchUrl(fsDataItem.Data)
					newFile.Write(dataFromURL)

				} else if parse.FileExists(fsDataItem.Data) {
					// Load file data and place into new file
					dataFromFile, _ := parse.FetchFile(fsDataItem.Data)
					newFile.Write(dataFromFile)

				} else {
					// Treat data as plain text
					newFile.WriteString(fsDataItem.Data)
				}
			}
		}

		newFile.Close()

		bar.Add(1)
	}

	fmt.Println()
	fmt.Println()
	fmt.Printf("%s in %s", message.Green("Structure deployed"), time.Since(timeStart))

	return 0
}

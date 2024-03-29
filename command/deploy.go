package command

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"gitlab.com/merrittcorp/fspop/exe"
	"gitlab.com/merrittcorp/fspop/parse"
	"gitlab.com/merrittcorp/fspop/structure"
	"gitlab.com/merrittcorp/fspop/ui"
)

type DeployCommand struct {
	*BaseCommand
}

func (c *DeployCommand) Synopsis() string {
	return "Deploy (create) a structure"
}

func (c *DeployCommand) Help() string {
	helpText := `
Usage: fspop deploy [options] NAME
  
  Deploy creates the file-structure defined in a structure config file.

  Structure files can be created using the command:
  $ fspop init <NAME>

Options:

` + c.Flags().Help()

	return strings.TrimSpace(helpText)
}

func (c *DeployCommand) Flags() *FlagMap {
	return GetFlagMap(FlagNamesGlobal)
}

func (c *DeployCommand) strictExit() {
	if c.Flags().Get("strict").Value == true {
		c.UI.Error("\nAn error occured while using the '--strict' flag.")
		// c.UI.Warn("(There is likely error messages above on what went wrong.)")
		os.Exit(1)
	}
}

func (c *DeployCommand) Run(args []string) int {
	args = c.Flags().Parse(c.UI, args)

	var path string

	if len(args) == 0 {
		// Use default structure file
		// Checks for both '.yaml' and '.yml' extensions
		path = parse.AddYamlExtension(parse.ElasticExtension(parse.DefaultYamlFileName))
		c.UI.Warn("No file entered.")
		c.strictExit()
		c.UI.Warn("Trying default '" + path + "' instead.\n")
	} else {
		path = parse.ElasticExtension(args[0])
	}

	// Fetch structure
	fsStructure := parse.FetchAndParseStructure(path)

	// Print structure stats
	c.UI.Output("Structure File")
	c.UI.Output("├── Data Variables       " + fmt.Sprint(len(fsStructure.Data)))
	c.UI.Output("├── Dynamic Variables    " + fmt.Sprint(len(fsStructure.Dynamic)))
	c.UI.Output("└── Structure Endpoints  " + fmt.Sprint(len(fsStructure.Items)))
	c.UI.Output("")

	// Check if entrypoint directory already exists.
	// If so, only deploy if --force is true
	if stat, err := os.Stat(fsStructure.GetEntrypoint()); err == nil && stat.IsDir() {
		if c.Flags().Get("force").Value == true {
			c.UI.Warn("--force flag is enabled. Existing files/folders may be overwritten.\n")
		} else {
			c.UI.Error("Entrypoint directory '" + fsStructure.GetEntrypoint() + "' already exists.\nfspop does not deploy to existing directories.")
			c.UI.Warn("\nUse '--force' flag to deploy to an existing directory.")
			os.Exit(1)
		}
	}

	// Record the total duration of this command
	timeStart := time.Now()

	// Initiate error slice and counter
	// Collects errors
	// TODO: make this a type + methods
	// errorSlice := make([]error, 0, 1)
	errorCount := 0

	// Initiate progress bar
	bar := ui.GetProgressBar(len(fsStructure.Items), " Deploying")

	// Loop each item endpoint
	// An endpoint can be both a file or directory
	for key, item := range fsStructure.Items {
		// Detect a dynamic item
		// Dynamic items are special and treated separately
		if len(item.DynamicKey) > 0 {
			// Get actual dynamic item from 'DynamicKey'
			fsDynamicItem, ok := fsStructure.Dynamic[item.DynamicKey]

			if ok {
				// Deploy dynamic item
				deployDynamicItem(fsStructure, item, fsDynamicItem, c, bar, &errorCount)
			} else {
				// Error: Dynamic key does not exist
				printError(c, bar, &errorCount, errors.New("dynamic key does not exist: '"+item.DynamicKey+"'"))
			}

			bar.Add(1)
			continue
		}

		deployItem(fsStructure, item, key, c, bar, &errorCount)

		bar.Add(1)
	}

	c.UI.Output("\n")

	// Check for any post-deploy 'Actions' scripts
	if len(fsStructure.Actions) > 0 {
		iRan := 0
		for _, fsAction := range fsStructure.Actions {
			if fsAction.CanRunOnOs() {
				iRan++
				if iRan > 1 {
					c.UI.Output("")
				}

				if len(fsAction.GetKeyOs()) > 0 {
					c.UI.Output(fmt.Sprintf("%s on %s", c.UI.Colorize(ui.WrapAtLength(fsAction.GetKeyWithoutOs(), 0), c.UI.InfoColor), fsAction.GetKeyOs()))
				} else {
					c.UI.Output(ui.WrapAtLength(fsAction.Key, 0))
				}

				printCommand := func(s string) {
					c.UI.Info(ui.WrapAtLength(fmt.Sprintf("  $ %s", s), 4))
				}

				// Check if first item in fsAction is not a command
				// and is actually a script to be executed
				// e.g. /usr/bin/backup.sh
				ok, scriptPath := exe.ScriptExists(fsAction.Script[0], fsStructure.Entrypoint)
				if len(fsAction.Script) == 1 && ok {
					// Run command and print output
					printCommand(fsAction.Script[0])
					err := exe.Run(exe.ScriptCommandExe(scriptPath), scriptPath, ".")
					if err != nil {
						errorCount++
						c.strictExit()
					}
					continue
				}

				// Loop each script command
				for _, command := range fsAction.Script {
					// Run command and print output
					printCommand(command)
					err := exe.Run(exe.GetOsShell(), command, fsStructure.Entrypoint)

					if err != nil {
						errorCount++
						c.strictExit()
					}
				}
			}
		}

		if len(fsStructure.Actions) > 0 {
			c.UI.Output("")
		}
	}

	if errorCount > 0 {
		c.UI.Warn("Use '--strict' flag to stop immediately if any errors occur\n")

		c.UI.Output(fmt.Sprintf("%s in %s", c.UI.Colorize("Structure deployed (with "+fmt.Sprint(errorCount)+" errors)", c.UI.WarnColor), time.Since(timeStart)))
		return 1
	}

	c.UI.Output(fmt.Sprintf("%s in %s", c.UI.Colorize("Structure deployed", c.UI.SuccessColor), time.Since(timeStart)))

	return 0
}

// Main procedure to deploy a dynamic item
func deployDynamicItem(fsStructure *structure.FspopStructure, item *structure.FspopItem, fsDynamicItem *structure.FspopDynamic, c *DeployCommand, bar *progressbar.ProgressBar, errorCount *int) {
	// Create empty fileData
	fileData := make([]byte, 0)

	// If dynamic item has type 'file'
	// Load file data before creating items
	if !fsDynamicItem.IsTypeDirectory() {
		// Resolve 'DataKey' into a data item
		if fsDataItem, ok := fsStructure.GetDataItem(fsDynamicItem.DataKey); ok {
			var err error
			fileData, err = resolveDataPayload(fsDataItem.Data)

			if err != nil {
				printError(c, bar, errorCount, errors.New("unable to get data payload for "+fsDynamicItem.Key+": '"+fsDataItem.Data+"'"))
				fileData = []byte(fsDataItem.Data) // Fallback to whatever the user set fsDataItem.Data to
			}
		}
	}

	// Loop n times
	// n = fsDynamicItem.Count = user defined
	fsDynamicItemMaxCount := fsDynamicItem.Start + fsDynamicItem.Count
	for i := fsDynamicItem.Start; i < fsDynamicItemMaxCount; i++ {
		itemPath, itemParentPath := fsDynamicItem.BuildItemPath(fsStructure.GetEntrypoint(), &item.Path, i)

		// Is directory
		if fsDynamicItem.IsTypeDirectory() {
			// Create full directory
			err := os.MkdirAll(itemPath, os.ModePerm)
			if err != nil {
				printError(c, bar, errorCount, errors.New("unable to make directory for "+fsDynamicItem.Key+": '"+itemPath+"'"))
			}

			// Is file
		} else {
			// Create parent directory
			err := os.MkdirAll(itemParentPath, os.ModePerm)
			if err != nil {
				printError(c, bar, errorCount, errors.New("unable to make directory for "+fsDynamicItem.Key+": '"+itemParentPath+"'"))
				continue
			}

			// Create file
			newFile, err := parse.CreateFile(itemPath)
			if err != nil {
				printError(c, bar, errorCount, errors.New("unable to create file for "+fsDynamicItem.Key+": '"+itemPath+"'"))
				newFile.Close()
				continue
			}

			if len(fileData) > 0 {
				_, err := newFile.Write(fileData)
				if err != nil {
					printError(c, bar, errorCount, errors.New("unable to add data to file for "+fsDynamicItem.Key+": '"+itemPath+"'"))
				}
			}

			newFile.Close()
		}
	}
}

// Main procedure to deploy an item (files/folders only, not dynamic items)
func deployItem(fsStructure *structure.FspopStructure, item *structure.FspopItem, key string, c *DeployCommand, bar *progressbar.ProgressBar, errorCount *int) {
	// Directory
	// Check if endpoint is a directory
	if structure.IsDirectory(key) {
		// Recursively make all directories
		err := os.MkdirAll(fmt.Sprintf("%s/%s", fsStructure.GetEntrypoint(), key), os.ModePerm)
		if err != nil {
			printError(c, bar, errorCount, errors.New("unable to make directory: '"+key+"'"))
		}

	} else {
		// File
		// Recursively make all parent directories
		err := os.MkdirAll(fmt.Sprintf("%s/%s", fsStructure.GetEntrypoint(), item.Path.ParentString()), os.ModePerm)
		if err != nil {
			printError(c, bar, errorCount, errors.New("unable to make directory: '"+item.Path.ParentString()+"'"))
			return
		}

		// Create empty file
		newFile, err := parse.CreateFile(fmt.Sprintf("%s/%s", fsStructure.GetEntrypoint(), item.Path.ToString()))
		if err != nil {
			printError(c, bar, errorCount, errors.New("unable to create file: '"+item.Path.ToString()+"'"))
			newFile.Close()
			return
		}

		// Resolve 'DataKey' into a data item
		if fsDataItem, ok := fsStructure.GetDataItem(item.DataKey); ok {
			err := fetchAndWriteToFile(newFile, fsDataItem.Data)
			if err != nil {
				printError(c, bar, errorCount, errors.New("unable to add data to file ("+fmt.Sprint(err)+"): '"+item.Path.ToString()+"'"))
			}
		}

		newFile.Close()
	}
}

// Prints an error while deploying
// Handles printing aorund the progress-bar
func printError(c *DeployCommand, bar *progressbar.ProgressBar, errorCount *int, err error) {
	// Remove the progress bar from the current line
	bar.Clear()

	// Check if first error
	// Setup error list for the first error
	if *errorCount == 0 {
		c.UI.Error("ERROR:")
	} else {
		fmt.Print("\r\033[A")
	}

	// Print error
	c.UI.Error(ui.WrapAtLength(fmt.Sprintf("  -- %s\n", err), 5))

	// Exif IF --strict
	if c.Flags().Get("strict").Value == true {
		bar.Add(1)
		c.UI.Output("")
		c.strictExit()
	}

	*errorCount++
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

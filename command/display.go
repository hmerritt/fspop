package command

import (
	"fmt"
	"strings"

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
	// var path string

	// if len(args) == 0 {
	// 	message.Warn("No file entered.")
	// 	message.Warn("Trying default '" + parse.DefaultYamlFile + "' instead.")
	// 	message.Text("")
	// 	path = parse.DefaultYamlFile
	// } else {
	// 	path = parse.ElasticExtension(args[0])
	// }

	// var fileData []byte
	// var fileError error

	// // Decide if URL or file
	// if parse.UseUrl(path) {
	// 	message.Spinner.Start("", " Fetching URL data...")

	// 	fileData, fileError = parse.FetchUrl(path)

	// 	message.Spinner.Stop()

	// 	if fileError != nil {
	// 		message.Error("Unable to fetch URL data.")
	// 		message.Error(fmt.Sprint(fileError))
	// 		fmt.Println()
	// 		message.Warn("Make sure the link is accessible and try again.")
	// 		return 2
	// 	}
	// } else {
	// 	fileData, fileError = parse.FetchFile(path)

	// 	if fileError != nil {
	// 		message.Error("Unable to open file.")
	// 		message.Error(fmt.Sprint(fileError))
	// 		fmt.Println()
	// 		message.Warn("Check the file is exists and try again.")
	// 		return 2
	// 	}
	// }

	// fsStructure := parse.ParseAndRefineYaml(fileData)
	// //parse.ParseAndRefineYaml(fileData)

	// // fmt.Println(fsStructure)

	// Mock children items
	fsMockChildrenItems := make([]*structure.FspopItem, 0)
	fsMockChildrenItems = append(fsMockChildrenItems, &structure.FspopItem{
		Path:       *structure.CreateFspopPath([]string{"disone/", "file.wow"}),
		IsDir:      false,
		IsEndpoint: true,
		HasData:    false,
	})
	fsMockChildrenItems = append(fsMockChildrenItems, &structure.FspopItem{
		Path:       *structure.CreateFspopPath([]string{"disone/", "2dir/"}),
		IsDir:      true,
		IsEndpoint: false,
		HasData:    false,
		Data:       "oh no",
	})

	// Mock items
	fsMockItems := make([]*structure.FspopItem, 0)
	fsMockItems = append(fsMockItems, &structure.FspopItem{
		Path:       *structure.CreateFspopPath([]string{"file.wow"}),
		IsDir:      false,
		IsEndpoint: true,
		HasData:    false,
	})
	fsMockItems = append(fsMockItems, &structure.FspopItem{
		Path:       *structure.CreateFspopPath([]string{"anotherfile.waw"}),
		IsDir:      false,
		IsEndpoint: true,
		HasData:    false,
	})
	fsMockItems = append(fsMockItems, &structure.FspopItem{
		Path:       *structure.CreateFspopPath([]string{"firstdir/"}),
		IsDir:      true,
		IsEndpoint: false,
		HasData:    false,
	})
	fsMockItems = append(fsMockItems, &structure.FspopItem{
		Path:       *structure.CreateFspopPath([]string{"disone/"}),
		IsDir:      true,
		IsEndpoint: false,
		HasData:    false,
		Children:   fsMockChildrenItems,
	})

	// Mock structure
	fsMockStructure := &structure.FspopStructure{
		Version: "2",
		Name:    "Mock structure",
		Items:   fsMockItems,
	}

	// = & fsMockStructure[item]
	// getItem, _ := fsMockStructure.Find(structure.CreateFspopPath([]string{"file.wow"}))
	// getItem.Data = "wooo"
	// fmt.Println(*getItem)
	// fmt.Println(fsMockStructure.Items[0])

	_, err := fsMockStructure.Find(structure.CreateFspopPath([]string{"disone/", "2dir/", "wownew/"}))
	fmt.Println(err)

	fsMockStructure.Add(&structure.FspopItem{
		Path:       *structure.CreateFspopPath([]string{"disone/", "2dir/", "wownew/"}),
		IsDir:      true,
		IsEndpoint: true,
		HasData:    false,
	})

	path, err2 := fsMockStructure.Find(structure.CreateFspopPath([]string{"disone/", "2dir/", "wownew/"}))
	fmt.Println(err2)
	fmt.Println(path)

	return 0
}

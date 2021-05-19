package parse

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
	"gitlab.com/merrittcorp/fspop/structure"
	"gopkg.in/yaml.v2"
)

func ParseAndRefineYaml(data []byte) *structure.FspopStructure {
	// Parse YAML
	yamlStructure, parseErr := ParseYaml(data)

	if parseErr != nil {
		message.Error("Unable to parse YAML file.")
		message.Error(fmt.Sprint(parseErr))
		fmt.Println()
		message.Warn("Check the file is valid YAML and try again.")
		os.Exit(2)
	}

	// Refine YAML
	refined := RefineYaml(yamlStructure)

	// TODO: catch errors when refining

	return refined
}

//
// Parse YAML byte array (fetch not included)
//
// Parses into an intermediary structure, this structure is messy
// and needs to be refinded before any real use.
//
func ParseYaml(data []byte) (structure.YamlStructure, error) {
	// Define structure
	structure := structure.YamlStructure{}

	// Parse YAML data
	err := yaml.Unmarshal(data, &structure)
	if err != nil {
		return structure, err
	}

	return structure, nil
}

//
// Refine the messy intermediary YamlStructure into an organsized,
// searchable structure which is used from here-on-out.
//
func RefineYaml(parsedYamlStructure structure.YamlStructure) *structure.FspopStructure {
	refinedStructure := &structure.FspopStructure{
		Version: parsedYamlStructure.Version,
		Name:    parsedYamlStructure.Name,
		Data:    make(map[string]*structure.FspopData),
		Dynamic: make(map[string]*structure.FspopDynamic),
		Items:   make(map[string]*structure.FspopItem),
	}

	// TODO: Refine 'Data'
	// TODO: Refine 'Dynamic'

	// Setup structure items
	fsPath := *structure.CreateFspopPath([]string{})

	callback := func(path structure.FspopPath, dataKey string, dynamicKey string) {
		refinedStructure.Items[path.ToString()] = &structure.FspopItem{
			Path:       path,
			DataKey:    dataKey,
			DynamicKey: dynamicKey,
		}
	}

	// Refine structure items
	RefineYamlItems(parsedYamlStructure.Structure, fsPath, callback)

	// TODO: build directory tree structure

	return refinedStructure
}

//
// Crawl through the 'structure:' items key in the messy parsed yaml
// figuring out whats-what and organising it one item at a time.
// Will detect and output file data and dynamic keys
//
func RefineYamlItems(structureInterface interface{}, pathStart structure.FspopPath, callback func(structure.FspopPath, string, string)) {
	// Unique path for each iteration
	// path := *structure.CreateFspopPath(pathStart.Path)
	path := structure.FspopPath{
		Path: pathStart.Path,
	}

	switch structureInterface.(type) {
	case string:
		itemName := fmt.Sprintf("%v", structureInterface)
		dynamicKey := ""

		// File or Directory name
		path.Append(itemName)

		// Check for a dynamic key
		if !structure.IsDirectory(itemName) && strings.HasPrefix(itemName, "$") {
			dynamicKey = itemName
		}

		callback(path, "", dynamicKey)

	case []interface{}:
		// Use type assertion to loop over []interface{}
		for _, v := range structureInterface.([]interface{}) {
			RefineYamlItems(v, path, callback)
		}

	case map[interface{}]interface{}:
		// Use type assertion to loop over map[string]interface{}
		for key, value := range structureInterface.(map[interface{}]interface{}) {
			// Interface 'key' is a directory name
			// Create a new unique path for each iteration,
			// prevents 'path' being carried forward and messing
			// with the callback later.
			path.Append(fmt.Sprintf("%v", key))

			// Check for file with a data variable
			// Use 'key' string value as data key
			if !structure.IsDirectory(fmt.Sprintf("%v", key)) {
				dataKey := fmt.Sprintf("%v", value)
				callback(path, dataKey, "")
				continue
			}

			RefineYamlItems(value, path, callback)
		}
	}
}

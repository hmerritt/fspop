package parse

import (
	"fmt"
	"os"

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

/*
 * Parse YAML byte array (fetch not included)
 * parses into an intermediary structure
 * this structure is messy and needs to be
 * refinded before any real use.
 */
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

/*
 * Refine the messy intermediary YamlStructure
 * into an organsized, searchable structure
 * which is used from here-on-out.
 */
func RefineYaml(parsedYamlStructure structure.YamlStructure) *structure.FspopStructure {
	refinedStructure := &structure.FspopStructure{
		Version: parsedYamlStructure.Version,
		Name:    parsedYamlStructure.Name,
		Data:    make(map[string]*structure.FspopData),
		Dynamic: make(map[string]*structure.FspopDynamic),
		Items:   make(map[string]*structure.FspopItem),
	}

	// Structure
	fsPath := *structure.CreateFspopPath([]string{})

	callback := func(path structure.FspopPath) {
		refinedStructure.Items[path.ToString()] = &structure.FspopItem{
			Path: path,
		}
	}

	RefineYamlItems(parsedYamlStructure.Structure, fsPath, callback)

	// TODO: Refine 'Data'
	// TODO: Refine 'Dynamic'

	return refinedStructure
}

func RefineYamlItems(structureInterface interface{}, pathStart structure.FspopPath, callback func(structure.FspopPath)) {
	// Unique path for each iteration
	path := *structure.CreateFspopPath(pathStart.Path)

	switch structureInterface.(type) {
	case string:
		// File or Directory name
		path.Append(fmt.Sprintf("%v", structureInterface))
		callback(path)

	case []interface{}:
		// Use type assertion to loop over []interface{}
		for _, v := range structureInterface.([]interface{}) {
			RefineYamlItems(v, path, callback)
		}

	case map[interface{}]interface{}:
		// Use type assertion to loop over map[string]interface{}
		for key, value := range structureInterface.(map[interface{}]interface{}) {
			// Interface 'key' is a directory name
			// create a new unique path for each iteration,
			// prevents 'path' being carried forward and messing
			// with the callback later.
			newPath := *structure.CreateFspopPath(path.Path)
			newPath.Append(fmt.Sprintf("%v", key))

			RefineYamlItems(value, path, callback)
		}
	}
}

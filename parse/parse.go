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
		Tree:    structure.StartTree(),
		// Data:    []structure.FspopData,
		// Dynamic: []structure.FspopDynamic,
		// Items:   []structure.FspopItem,
	}

	// Structure
	fsPath := structure.FspopStructurePath{
		Path: []string{},
	}

	callback := func(path structure.FspopStructurePath, isEndpoint bool) {
		refinedStructure.Items = append(refinedStructure.Items, &structure.FspopItem{
			Path:       path,
			IsDir:      structure.IsDirectory(path.Last()),
			IsEndpoint: isEndpoint,
			HasData:    false,
			Data:       "",
		})
	}

	RefineYamlItems(parsedYamlStructure.Structure, fsPath, callback)

	// TODO: Refine 'Data'
	// TODO: Refine 'Dynamic'

	return refinedStructure
}

func RefineYamlItems(structureInterface interface{}, pathStart structure.FspopStructurePath, callback func(structure.FspopStructurePath, bool)) {
	// Unique path for each iteration
	path := structure.FspopStructurePath{
		Path: pathStart.Path,
	}

	switch structureInterface.(type) {
	case string:
		// File or Directory name
		path.Append(fmt.Sprintf("%v", structureInterface))
		callback(path, true)

	case []interface{}:
		// Use type assertion to loop over []interface{}
		for _, v := range structureInterface.([]interface{}) {
			RefineYamlItems(v, path, callback)
		}

	case map[interface{}]interface{}:
		// Use type assertion to loop over map[string]interface{}
		for key, value := range structureInterface.(map[interface{}]interface{}) {
			// Interface 'key' is a directory name
			path.Append(fmt.Sprintf("%v", key))
			callback(path, false)

			RefineYamlItems(value, path, callback)
		}
	}
}

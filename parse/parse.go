package parse

import (
	"fmt"
	"os"

	"gitlab.com/merrittcorp/fspop/message"
	"gitlab.com/merrittcorp/fspop/structure"
	"gopkg.in/yaml.v2"
)

func FetchAndParseStructure(path string) *structure.FspopStructure {
	// Fetch structure file
	fileData := FetchYaml(path)

	// Parse YAML and refine into a useable structure
	fsStructure := ParseAndRefineYaml(fileData)

	return fsStructure
}

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
	refined, refineErr := RefineYaml(yamlStructure)

	if refineErr != nil {
		message.Error("Structure file is invalid or has missing parts.")
		message.Error(fmt.Sprint(refineErr))
		fmt.Println()
		message.Warn("Check the structure file is valid and try again.")
		os.Exit(2)
	}

	return refined
}

// Parse YAML byte array (fetch not included)
//
// Parses into an intermediary structure, this structure is messy
// and needs to be refinded before any real use.
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

// Refine the messy intermediary YamlStructure into an organsized,
// searchable structure which is used from here-on-out.
func RefineYaml(parsedYamlStructure structure.YamlStructure) (*structure.FspopStructure, error) {
	refinedStructure := &structure.FspopStructure{
		Version: parsedYamlStructure.Version,
		Name:    parsedYamlStructure.Name,
		Data:    make(map[string]*structure.FspopData),
		Dynamic: make(map[string]*structure.FspopDynamic),
		Items:   make(map[string]*structure.FspopItem),
	}

	// Check if parsed YAML is a valid structure file
	isValid, isValidErr := structure.IsValid(parsedYamlStructure)
	if !isValid {
		return refinedStructure, isValidErr
	}

	if parsedYamlStructure.Data != nil {
		// Refine 'data:' items
		refineYamlData(&parsedYamlStructure, refinedStructure)
	}

	if parsedYamlStructure.Dynamic != nil {
		// Refine 'dynamic:' items
		refineYamlDynamic(&parsedYamlStructure, refinedStructure)
	}

	// Setup structure items
	if parsedYamlStructure.Structure != nil {
		// Refine 'structure:' items
		refineYamlItems(&parsedYamlStructure, refinedStructure)
	}

	// TODO: build directory tree structure

	return refinedStructure, nil
}

// Refine 'data:' key in yaml structure file
func refineYamlData(parsedYamlStructure *structure.YamlStructure, refinedStructure *structure.FspopStructure) {
	callback := func(fsData structure.FspopData) {
		refinedStructure.Data[fsData.Key] = &fsData
	}
	parsedYamlStructure.TraverseData(callback)
}

// Refine 'dynamic:' key in yaml structure file
func refineYamlDynamic(parsedYamlStructure *structure.YamlStructure, refinedStructure *structure.FspopStructure) {
	callback := func(fsDynamic structure.FspopDynamic) {
		refinedStructure.Dynamic[fsDynamic.Key] = &fsDynamic
	}
	parsedYamlStructure.TraverseDynamic(callback)
}

// Crawl through the 'structure:' items key in the messy parsed yaml
// figuring out whats-what and organising it one item at a time.
// Will detect and output file data and dynamic keys.
func refineYamlItems(parsedYamlStructure *structure.YamlStructure, refinedStructure *structure.FspopStructure) {
	callback := func(fsItem structure.FspopItem) {
		refinedStructure.Items[fsItem.Path.ToString()] = &fsItem
	}
	parsedYamlStructure.TraverseStructureEndpoints(callback)
}

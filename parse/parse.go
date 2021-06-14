package parse

import (
	"fmt"
	"os"

	"gitlab.com/merrittcorp/fspop/structure"
	"gitlab.com/merrittcorp/fspop/ui"
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

	UI := ui.GetUi()

	if parseErr != nil {
		UI.Error("Unable to parse YAML file.")
		UI.Error(fmt.Sprint(parseErr))
		UI.Warn("\nCheck the file is valid YAML and try again.")
		os.Exit(2)
	}

	// Refine YAML
	refined, refineErr := RefineYaml(yamlStructure)

	if refineErr != nil {
		UI.Error("Structure file is invalid or has missing parts.")
		UI.Error(fmt.Sprint(refineErr))
		UI.Warn("\nCheck the structure file is valid and try again.")
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
		Version:    parsedYamlStructure.Version,
		Name:       parsedYamlStructure.Name,
		Entrypoint: parsedYamlStructure.Entrypoint,
		Data:       make(map[string]*structure.FspopData),
		Dynamic:    make(map[string]*structure.FspopDynamic),
		Items:      make(map[string]*structure.FspopItem),
	}

	// Check if parsed YAML is a valid structure file
	isValid, isValidErr := parsedYamlStructure.IsValid()
	if !isValid {
		return refinedStructure, isValidErr
	}

	// Refine 'data:' items
	if parsedYamlStructure.Data != nil {
		callback := func(fsData structure.FspopData) {
			refinedStructure.Data[fsData.Key] = &fsData
		}
		parsedYamlStructure.TraverseData(callback)
	}

	// Refine 'dynamic:' items
	if parsedYamlStructure.Dynamic != nil {
		callback := func(fsDynamic structure.FspopDynamic) {
			refinedStructure.Dynamic[fsDynamic.Key] = &fsDynamic
		}
		parsedYamlStructure.TraverseDynamic(callback)
	}

	// Refine 'structure:' items
	if parsedYamlStructure.Structure != nil {
		callback := func(fsItem structure.FspopItem) {
			refinedStructure.Items[fsItem.Path.ToString()] = &fsItem
		}
		parsedYamlStructure.TraverseStructureEndpoints(callback)
	}

	// TODO: build directory tree structure

	return refinedStructure, nil
}

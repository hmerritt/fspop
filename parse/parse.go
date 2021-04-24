package parse

import (
	"gopkg.in/yaml.v2"
)

type FspopStructure struct {
	Version   string
	Name      string
	Data      interface{}
	Dynamic   interface{}
	Structure interface{}
}

func ParseYaml(path string) FspopStructure {
	// Get YAML
	data := FetchYaml(path)

	// Define structure
	structure := FspopStructure{}

	// Parse YAML data
	yaml.Unmarshal([]byte(data), &structure)

	return structure
}

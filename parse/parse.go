package parse

import (
	"gitlab.com/merrittcorp/fspop/structure"
	"gopkg.in/yaml.v2"
)

func ParseYaml(data string) structure.FspopStructure {
	// Get YAML
	data := FetchYaml(path)

	// Define structure
	structure := structure.FspopStructure{}

	// Parse YAML data
	yaml.Unmarshal([]byte(data), &structure)

	return structure
}

package parse

import (
	"gitlab.com/merrittcorp/fspop/structure"
	"gopkg.in/yaml.v2"
)

func ParseYaml(data []byte) structure.FspopStructure {
	// Define structure
	structure := structure.FspopStructure{}

	// Parse YAML data
	yaml.Unmarshal(data, &structure)

	return structure
}

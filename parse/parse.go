package parse

import (
	"gitlab.com/merrittcorp/fspop/structure"
	"gopkg.in/yaml.v2"
)

func ParseYaml(data []byte) (structure.FspopStructure, error) {
	// Define structure
	structure := structure.FspopStructure{}

	// Parse YAML data
	err := yaml.Unmarshal(data, &structure)
	if err != nil {
		return structure, err
	}

	return structure, nil
}

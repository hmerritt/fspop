package structure

import (
	"fmt"
	"strings"
)

// Verifies integraty of a structure
//
// Checks for required keys and expected types
func IsValid(parsedYamlStructure YamlStructure) bool {
	// TODO: return array of all problems (not just static bool)

	// Required keys

	// Structure key
	if parsedYamlStructure.Structure == nil {
		return false
	}

	// Expected types

	// Data key, check if it exists
	if parsedYamlStructure.Data != nil {
		// Data key should be an []interface{}
		dataMapType, ok := parsedYamlStructure.Data.([]interface{})
		if !ok {
			return false
		}

		if len(dataMapType) > 0 {
			for _, dataMap := range dataMapType {
				// Invididual data key should have the type map[interface{}]interface{}
				// data key + value
				if _, ok := dataMap.(map[interface{}]interface{}); !ok {
					return false
				}
			}
		}
	}

	// Dynamic key, check if it exists
	if parsedYamlStructure.Dynamic != nil {
		// Dynamic key should be an []interface{}
		dynamicMapType, ok := parsedYamlStructure.Dynamic.([]interface{})
		if !ok {
			return false
		}

		if len(dynamicMapType) > 0 {
			for _, dynamicMap := range dynamicMapType {
				// Invididual dynamic key should have the type map[interface{}]interface{}
				// dynamic key + value map
				dynamicItemMapType, ok := dynamicMap.(map[interface{}]interface{})
				if !ok {
					return false
				}

				for _, dynamicItemMap := range dynamicItemMapType {
					// Dynamic item map should be an []interface{}
					dynamicItemMapType, ok := dynamicItemMap.([]interface{})
					if !ok {
						return false
					}

					for _, dynamicValueMap := range dynamicItemMapType {
						dynamicValueMapType, ok := dynamicValueMap.(map[interface{}]interface{})
						if !ok {
							return false
						}

						// Check each dynamic item value
						// Verify type, some items need to be an int
						for variable, value := range dynamicValueMapType {
							// Ditermine variable name
							switch strings.ToLower(fmt.Sprint(variable)) {
							case "amount":
								// Ammount must be an int
								if _, ok := value.(int); !ok {
									return false
								}
							case "count":
								// Count must be an int
								if _, ok := value.(int); !ok {
									return false
								}
							case "padded":
								// Padded must be a bool
								if _, ok := value.(bool); !ok {
									return false
								}
							case "start":
								// Start must be an int
								if _, ok := value.(int); !ok {
									return false
								}
							}
						}
					}
				}
			}
		}
	}

	return true
}

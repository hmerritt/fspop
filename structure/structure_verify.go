package structure

import (
	"errors"
	"fmt"
	"strings"
)

// Verifies integraty of a structure
//
// Checks for required keys and expected types
func IsValid(parsedYamlStructure YamlStructure) (bool, error) {
	// TODO: return array of all problems (not just bool and first error)

	// Required keys

	// Structure key
	if parsedYamlStructure.Structure == nil {
		return false, errors.New("'structure:' key not found in structure file")
	}

	// Expected types

	// Data key, check if it exists
	if parsedYamlStructure.Data != nil {
		// Data key should be an []interface{}
		dataMapType, ok := parsedYamlStructure.Data.([]interface{})
		if !ok {
			return false, errors.New("'data:' key format is invalid")
		}

		if len(dataMapType) > 0 {
			for _, dataMap := range dataMapType {
				// Invididual data key should have the type map[interface{}]interface{}
				// data key + value
				if _, ok := dataMap.(map[interface{}]interface{}); !ok {
					return false, errors.New("'data:' key format is invalid")
				}
			}
		}
	}

	// Dynamic key, check if it exists
	if parsedYamlStructure.Dynamic != nil {
		// Dynamic key should be an []interface{}
		dynamicMapType, ok := parsedYamlStructure.Dynamic.([]interface{})
		if !ok {
			return false, errors.New("'dynamic:' key format is invalid")
		}

		if len(dynamicMapType) > 0 {
			for _, dynamicMap := range dynamicMapType {
				// Invididual dynamic key should have the type map[interface{}]interface{}
				// dynamic key + value map
				dynamicItemMapType, ok := dynamicMap.(map[interface{}]interface{})
				if !ok {
					return false, errors.New("'dynamic:' key format is invalid")
				}

				for key, dynamicItemMap := range dynamicItemMapType {
					// Dynamic item map should be an []interface{}
					dynamicItemMapType, ok := dynamicItemMap.([]interface{})
					if !ok {
						return false, errors.New("dynamic '" + fmt.Sprint(key) + "' key format is invalid")
					}

					for _, dynamicValueMap := range dynamicItemMapType {
						dynamicValueMapType, ok := dynamicValueMap.(map[interface{}]interface{})
						if !ok {
							return false, errors.New("dynamic '" + fmt.Sprint(key) + "' key format is invalid")
						}

						// Check each dynamic item value
						// Verify type, some items need to be an int
						for variable, value := range dynamicValueMapType {
							// Ditermine variable name
							switch strings.ToLower(fmt.Sprint(variable)) {
							case "amount":
								// Ammount must be an int
								if _, ok := value.(int); !ok {
									return false, errors.New("dynamic '" + fmt.Sprint(key) + "' key, 'ammount' value should be a number")
								}
							case "count":
								// Count must be an int
								if _, ok := value.(int); !ok {
									return false, errors.New("dynamic '" + fmt.Sprint(key) + "' key, 'count' value should be a number")
								}
							case "padded":
								// Padded must be a bool
								if _, ok := value.(bool); !ok {
									return false, errors.New("dynamic '" + fmt.Sprint(key) + "' key, 'padded' value should be true/false")
								}
							case "start":
								// Start must be an int
								if _, ok := value.(int); !ok {
									return false, errors.New("dynamic '" + fmt.Sprint(key) + "' key, 'start' value should be a number")
								}
							}
						}
					}
				}
			}
		}
	}

	return true, nil
}

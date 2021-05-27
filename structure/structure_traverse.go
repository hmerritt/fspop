package structure

import (
	"fmt"
	"strings"
)

// Traverse each data key & value in an unrefined YAML structure
func (fsYamlStruct *YamlStructure) TraverseData(callback func(FspopData)) {
	// Skip traversal if data key is empty
	if fsYamlStruct.Data == nil {
		return
	}

	// Iterate each map individually
	for _, dataMap := range fsYamlStruct.Data.([]interface{}) {
		// Get key and value from map
		for key, value := range dataMap.(map[interface{}]interface{}) {
			callback(FspopData{
				Key:  fmt.Sprint(key),
				Data: fmt.Sprint(value),
			})
		}
	}
}

// Traverse each dynamic key map in an unrefined YAML structure
func (fsYamlStruct *YamlStructure) TraverseDynamic(callback func(FspopDynamic)) {
	// Skip traversal if data key is empty
	if fsYamlStruct.Dynamic == nil {
		return
	}

	// Iterate each dynamic variable map individually
	for _, dynamicMap := range fsYamlStruct.Dynamic.([]interface{}) {
		// Get dynamic key and it's values in map form
		for key, dynamicItemMap := range dynamicMap.(map[interface{}]interface{}) {
			// Create dynamic key struct
			fsDynamic := FspopDynamic{
				// Make sure key has '$' prefix
				// Removes exsting '$' and then adds it back
				Key: fmt.Sprintf("$%s", strings.TrimPrefix(fmt.Sprint(key), "$")),
			}

			// Iterate all dynamic values
			for _, dynamicValueMap := range dynamicItemMap.([]interface{}) {
				// Get dynamic item variables from map
				for variable, value := range dynamicValueMap.(map[interface{}]interface{}) {
					// Ditermine variable name and place value in the correct place
					switch strings.ToLower(fmt.Sprint(variable)) {
					case "amount":
						fsDynamic.Count = value.(int)
					case "count":
						fsDynamic.Count = value.(int)
					case "data":
						fsDynamic.DataKey = fmt.Sprint(value)
					case "type":
						fsDynamic.Type = strings.ToLower(fmt.Sprint(value))
					case "name":
						fsDynamic.Name = fmt.Sprint(value)
					case "padded":
						fsDynamic.Padded = value.(bool)
					case "start":
						fsDynamic.Start = value.(int)
					}
				}
			}

			callback(fsDynamic)
		}
	}
}

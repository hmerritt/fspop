package structure

import (
	"fmt"
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

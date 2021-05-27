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

// Traverse each structure endpoint as an 'FspopItem'
// note: only endpoints are traversed and not each stage of the structure
func (fsYamlStruct *YamlStructure) TraverseStructure(callback func(FspopItem)) {
	// Crawl each endpoint recursivly
	crawlYamlStructureItems(fsYamlStruct.Structure, *CreateFspopPath([]string{}), callback)
}

// Crawl through the 'structure:' items key in the messy parsed yaml
// figuring out whats-what and organising it one item at a time.
// Will detect and output file data and dynamic keys.
func crawlYamlStructureItems(structureInterface interface{}, pathStart FspopPath, callback func(FspopItem)) {
	// Unique path for each iteration
	// Make a deep copy of path array
	path := FspopPath{
		Path: deepCopy(pathStart.Path),
	}

	switch structureInterface.(type) {
	case string:
		itemName := fmt.Sprintf("%v", structureInterface)
		dynamicKey := ""

		// File or Directory name
		path.Append(itemName)

		// Check for a dynamic key
		if !IsDirectory(itemName) && strings.HasPrefix(itemName, "$") {
			dynamicKey = itemName
		}

		callback(FspopItem{
			Path:       path,
			DataKey:    "",
			DynamicKey: dynamicKey,
		})

	case []interface{}:
		// Use type assertion to loop over []interface{}
		for _, v := range structureInterface.([]interface{}) {
			crawlYamlStructureItems(v, path, callback)
		}

	case map[interface{}]interface{}:
		// Use type assertion to loop over map[string]interface{}
		for key, value := range structureInterface.(map[interface{}]interface{}) {
			// Interface 'key' is a directory name
			// Create a new unique path for each iteration,
			// prevents 'path' being carried forward and messing
			// with the callback later.
			path.Append(fmt.Sprintf("%v", key))

			// Check for file with a data variable
			// Use 'key' string value as data key
			if !IsDirectory(fmt.Sprintf("%v", key)) {
				dataKey := fmt.Sprintf("%v", value)
				callback(FspopItem{
					Path:       path,
					DataKey:    dataKey,
					DynamicKey: "",
				})
				continue
			}

			crawlYamlStructureItems(value, path, callback)
		}
	}
}

// Make a deep copy of a string slice
func deepCopy(arr []string) []string {
	newArr := make([]string, len(arr))
	copy(newArr, arr)
	return newArr
}

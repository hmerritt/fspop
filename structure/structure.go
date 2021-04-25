package structure

import (
	"fmt"
	"strings"
)

type FspopStructure struct {
	Version   string
	Name      string
	Data      interface{}
	Dynamic   interface{}
	Structure interface{}
}

type FspopStructureAction func(string)

func Crawl(structure interface{}, path string, callback FspopStructureAction) {
	// fmt.Printf("Version: %s \n", structure.Version)
	// fmt.Printf("Name: %s \n", structure.Name)
	// fmt.Printf("Data: %v \n", structure.Data)
	// fmt.Printf("Data: %v \n", structure.Data.([]interface{})[1].(map[interface{}]interface{})["data_file"])
	// fmt.Printf("Dynamic: %v \n", structure.Dynamic)
	// fmt.Printf("Structure: %v \n", structure.Structure)
	switch structure.(type) {
	case string:
		//fmt.Printf("%v is an interface \n", structure)
		fmt.Printf("%s %v\n", path, structure)

	case []interface{}:
		//fmt.Printf("%v is a slice of interface \n ", structure)
		for _, v := range structure.([]interface{}) { // use type assertion to loop over []interface{}
			Crawl(v, path, func(path string) {})
		}

	case map[interface{}]interface{}:
		//fmt.Printf("%v is a map \n\n\n", structure)

		// Recurse interface
		for key, value := range structure.(map[interface{}]interface{}) { // use type assertion to loop over map[string]interface{}
			//fmt.Println(key)
			fmt.Printf("%s %v\n", path, key)
			newPath := path + fmt.Sprintf("%v", key)
			Crawl(value, newPath, func(path string) {})
		}
	}
}

func IsDirectory(str string) bool {
	return (strings.HasPrefix(str, "/") || strings.HasSuffix(str, "/"))
}

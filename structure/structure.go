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

func Crawl(structure interface{}, pathStart *FspopStructurePath, callback func(string)) {
	// fmt.Printf("Version: %s \n", structure.Version)
	// fmt.Printf("Name: %s \n", structure.Name)
	// fmt.Printf("Data: %v \n", structure.Data)
	// fmt.Printf("Data: %v \n", structure.Data.([]interface{})[1].(map[interface{}]interface{})["data_file"])
	// fmt.Printf("Dynamic: %v \n", structure.Dynamic)
	// fmt.Printf("Structure: %v \n", structure.Structure)

	// Unique path for each iteration
	path := &FspopStructurePath{
		Path: pathStart.Path,
	}

	switch structure.(type) {
	case string:
		//fmt.Printf("%v is an interface \n", structure)
		//fmt.Printf("%s %v\n", cliArrow, structure)
		path.Append(fmt.Sprintf("%v", structure))
		fmt.Printf("%v\n", path.ToString())

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
			//fmt.Printf("%s %v\n", cliArrow, key)
			path.Append(fmt.Sprintf("%v", key))
			fmt.Printf("%v\n", path.ToString())

			Crawl(value, path, func(path string) {})
		}
	}
}

func IsDirectory(path string) bool {
	return (strings.HasPrefix(path, "/") || strings.HasSuffix(path, "/"))
}

func StandardizeDirectory(path string) string {
	if strings.HasPrefix(path, "/") {
		path = path[1:] + "/"
	}
	return path
}

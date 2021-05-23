package parse

import (
	"fmt"
	"path/filepath"
	"strings"
)

var YamlExtensions = [2]string{".yaml", ".yml"}

var DefaultYamlFileName = "structure"
var DefaultYamlFile = fmt.Sprintf("%s.%s", DefaultYamlFileName, YamlExtensions[1])

func FileExtension(filename string) string {
	return filepath.Ext(filename)
}

func IsYamlFile(filename string) bool {
	for _, yamlext := range YamlExtensions {
		if strings.HasSuffix(filename, yamlext) {
			return true
		}
	}
	return false
}

func AddYamlExtension(filename string) string {
	if IsYamlFile(filename) {
		return filename
	}

	return filename + YamlExtensions[1]
}

func ElasticExtension(filename string) string {
	if IsYamlFile(filename) || FileExists(filename) {
		return filename
	}

	// Check for an existing file using all yaml extensions
	for yamlExtIndex := range YamlExtensions {
		elastic := filename + YamlExtensions[yamlExtIndex]
		if FileExists(elastic) {
			return elastic
		}
	}

	// Return original filename if nothing found
	return filename
}

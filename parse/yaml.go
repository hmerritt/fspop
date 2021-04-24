package parse

import (
	"path/filepath"
)

var YamlExtensions = [2]string{".yaml", ".yml"}

func FileExtension(filename string) string {
	return filepath.Ext(filename)
}

func IsYamlFile(filename string) bool {
	extension := FileExtension(filename)
	for _, yamlext := range YamlExtensions {
		if extension == yamlext {
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

	elastic := AddYamlExtension(filename)

	if FileExists(elastic) {
		return elastic
	} else {
		return filename
	}
}

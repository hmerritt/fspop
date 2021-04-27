package structure

import (
	"strings"
)

type YamlStructure struct {
	Version   string
	Name      string
	Data      interface{}
	Dynamic   interface{}
	Structure interface{}
}

type FspopData struct {
	Key  string
	Data string
}

type FspopDynamic struct {
	Key    string
	Count  int
	Data   FspopData
	Type   string
	Name   string
	Padded bool
}

type FspopItem struct {
	Path       FspopStructurePath
	IsDir      bool
	IsEndpoint bool // Tree endpoint (a file, or a directory with no sub-directories)
	HasData    bool
	Data       string
	Children   []FspopItem
}

type FspopStructure struct {
	Version string
	Name    string
	Data    []FspopData
	Dynamic []FspopDynamic
	Items   []FspopItem
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

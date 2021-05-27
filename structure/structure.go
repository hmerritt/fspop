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
	Key     string
	Count   int
	DataKey string
	Type    string
	Name    string
	Padded  bool
	Start   int
}

type FspopItem struct {
	Path       FspopPath
	DynamicKey string
	DataKey    string
}

type FspopStructure struct {
	Version string
	Name    string
	Data    map[string]*FspopData
	Dynamic map[string]*FspopDynamic
	Items   map[string]*FspopItem
}

func GetDefaultDynamicItem() *FspopDynamic {
	return &FspopDynamic{
		Count:  1,
		Type:   "file",
		Name:   "$num",
		Padded: true,
		Start:  1,
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

// Checks if a path exists in Items
func (fsStruct *FspopStructure) Exists(pathToFind *FspopPath) bool {
	if _, ok := fsStruct.Items[pathToFind.ToString()]; ok {
		return true
	} else {
		return false
	}
}

// Crawl each value in a structure
func (fsStruct *FspopStructure) Crawl(callback func(FspopItem)) {
	for _, v := range fsStruct.Items {
		callback(*v)
	}
}

// Count item endpoints
//
// Ideally, this would count all unique nodes not just endpoints.
func (fsStruct *FspopStructure) Count() int {
	return len(fsStruct.Items)
}

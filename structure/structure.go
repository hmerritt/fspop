package structure

import (
	"errors"
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
	Children   []*FspopItem // Re-create tree structure
}

type FspopStructure struct {
	Version string
	Name    string
	Data    []FspopData
	Dynamic []FspopDynamic
	Items   []*FspopItem // map[FspopPath]FspopItem
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

/*
 * Find an FspopItem through it's Path
 * would like to use a recursive function for this but would require
 * an extra param such as, Find(items *[]FspopItem, ...) which is
 * not great to use.
 *
 */
func (fsStruct *FspopStructure) Find(pathToFind *FspopStructurePath) (*FspopItem, error) {
	items := fsStruct.Items
	loop := true

	// Traverse until a find, or til end
	for loop {
		// Loop FspopItem slice
		for _, i := range items {
			// Match path
			if i.Path.ToString() == pathToFind.ToString() {
				// Found!
				i.Data = "changed"
				return i, nil

			} else if len(i.Children) > 0 {
				// Recurse deeper if item has children
				items = i.Children
			} else {
				// Item not found, end loop
				loop = false
			}
		}
	}

	// Path not found :(
	return &FspopItem{}, errors.New("path not found")
}

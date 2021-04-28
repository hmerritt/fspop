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
	itemsSlice := [][]*FspopItem{fsStruct.Items}
	count := 0

	// Traverse until a find, or til end
	for count <= len(itemsSlice)-1 {
		items := itemsSlice[count]

		// Loop FspopItem slice
		for _, i := range items {
			// Match path
			if i.Path.ToString() == pathToFind.ToString() {
				// Found!
				return i, nil

			} else if len(i.Children) > 0 {
				// Recurse deeper if item has children
				itemsSlice = append(itemsSlice, i.Children)
			}
		}

		count++
	}

	// Path not found :(
	return &FspopItem{}, errors.New("path not found")
}

/*
 * Add an FspopItem to the structure
 */
func (fsStruct *FspopStructure) Add(itemToAdd *FspopItem) error {
	if len(itemToAdd.Path.Path) < 1 {
		return errors.New("fspopitem not added to structure. path is empty")
	} else if len(itemToAdd.Path.Path) == 1 {
		fsStruct.Items = append(fsStruct.Items, itemToAdd)
		return nil
	}

	// Find parent
	parentPath := CreateFspopPath(itemToAdd.Path.Path[:itemToAdd.Path.Length()-1])
	parent, err := fsStruct.Find(parentPath)

	if err == nil && !parent.IsEndpoint {
		parent.Children = append(parent.Children, itemToAdd)
		return nil
	} else {
		return errors.New("fspopitem not added to structure")
	}
}

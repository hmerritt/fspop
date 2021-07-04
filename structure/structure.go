package structure

import (
	"fmt"
	"strings"
)

type YamlStructure struct {
	Version    string
	Name       string
	Entrypoint string
	Actions    interface{}
	Data       interface{}
	Dynamic    interface{}
	Structure  interface{}
}

type FspopAction struct {
	Key    string
	Script []string
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
	Version    string
	Name       string
	Entrypoint string
	Actions    []*FspopAction
	Data       map[string]*FspopData
	Dynamic    map[string]*FspopDynamic
	Items      map[string]*FspopItem
}

//
// FspopData methods
//

//
// FspopDynamic methods
//

// Returns true if dynamic item has type of 'directory'.
// Checks for multiple ways of defining directory type
func (fsDynamic *FspopDynamic) IsTypeDirectory() bool {
	return (strings.HasPrefix(fsDynamic.Type, "dir") || strings.HasPrefix(fsDynamic.Type, "fol"))
}

// Returns item path and it's parent directory path for an individual item
func (fsDynamic *FspopDynamic) BuildItemPath(entrypoint string, path *FspopPath, i int) (string, string) {
	// Maximum count any item can be in this dynamic set
	fsDynamicItemMaxCount := fsDynamic.Start + fsDynamic.Count

	// Padding of item count, is relative to max count '1' => '01'
	// Set default padding of '1'
	itemCountPadWidth := 1

	// If item has padding set to true,
	// calculate new padding width from max item count
	if fsDynamic.Padded {
		itemCountPadWidth = len(fmt.Sprint(fsDynamicItemMaxCount))
	}

	// Get count plus correct padding
	itemCountPadded := fmt.Sprintf("%0*d", itemCountPadWidth, i)

	// Check if item name does NOT have '$num'
	// If not, append '$num' to each item is different
	if !strings.Contains(fsDynamic.Name, "$num") {
		fsDynamic.Name = fmt.Sprintf("%s_$num", fsDynamic.Name)
	}

	// Replace '$num' with actual count + any padding
	itemName := strings.ReplaceAll(fsDynamic.Name, "$num", itemCountPadded)

	// Build parent directory path of item
	itemParentPath := fmt.Sprintf("%s/%s/", entrypoint, path.ParentString())

	// Build full item path
	itemPath := fmt.Sprintf("%s/%s", itemParentPath, itemName)

	// Return both item and parent path
	return itemPath, itemParentPath
}

// Returns a dynamic item with basic fallback values
func GetDefaultDynamicItem() *FspopDynamic {
	return &FspopDynamic{
		Count:  1,
		Type:   "file",
		Name:   "$num",
		Padded: true,
		Start:  1,
	}
}

//
// FspopStructure methods
//

// Get entrypoint path
//
// Returns 'Entrypoint' key, or 'Name' as a fallback
func (fsStruct *FspopStructure) GetEntrypoint() string {
	if fsStruct.Entrypoint == "" {
		return strings.TrimSuffix(fsStruct.Name, "/")
	}
	return strings.TrimSuffix(fsStruct.Entrypoint, "/")
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

// Resolves a dataKey into full FspopData item
// Returns 'FspopData' and if the data item exists (as a bool)
func (fsStruct *FspopStructure) GetDataItem(dataKey string) (*FspopData, bool) {
	if len(dataKey) == 0 {
		return nil, false
	}

	fsDataItem, ok := fsStruct.Data[dataKey]
	return fsDataItem, ok
}

//
// Generic structure methods
//

func IsDirectory(path string) bool {
	return (strings.HasPrefix(path, "/") || strings.HasSuffix(path, "/"))
}

func StandardizeDirectory(path string) string {
	if strings.HasPrefix(path, "/") {
		path = path[1:] + "/"
	}
	return path
}

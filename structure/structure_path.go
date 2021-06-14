package structure

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/merrittcorp/fspop/ui"
)

type FspopPath struct {
	Path []string
}

func CreateFspopPath(pathInit []string) *FspopPath {
	return &FspopPath{
		Path: pathInit,
	}
}

func (fsPath *FspopPath) Actual() string {
	return strings.TrimSuffix(fsPath.ToString(), "/")
}

func (fsPath *FspopPath) ParentString() string {
	if fsPath.Length() < 2 {
		return ""
	}
	parent := strings.Join(fsPath.Path[:fsPath.Length()-1], "")
	return strings.TrimSuffix(parent, "/")
}

func (fsPath *FspopPath) ToString() string {
	return strings.Join(fsPath.Path[:], "")
}

func (fsPath *FspopPath) IsEmpty() bool {
	return len(fsPath.Path) == 0
}

func (fsPath *FspopPath) First() string {
	if fsPath.IsEmpty() {
		return ""
	}
	return fsPath.Path[0]
}

func (fsPath *FspopPath) Last() string {
	if fsPath.IsEmpty() {
		return ""
	}
	return fsPath.Path[len(fsPath.Path)-1]
}

func (fsPath *FspopPath) Append(path string) {
	// Can't append if a file has been reached
	if !fsPath.IsEmpty() && !IsDirectory(fsPath.Last()) {
		// TODO: Prevent this from needing os.Exit
		UI := ui.GetUi()
		UI.Error(path)
		UI.Error(fmt.Sprint(fsPath.Path))
		UI.Error("\nfspop internal fatal error. tried to append after a file in fspopPath")
		os.Exit(1)
	}

	if IsDirectory(path) {
		// Fix directory
		path = StandardizeDirectory(path)
	}

	fsPath.Path = append(fsPath.Path, path)
}

func (fsPath *FspopPath) Length() int {
	return len(fsPath.Path)
}

func (fsPath *FspopPath) Name() string {
	if fsPath.IsEmpty() {
		return ""
	}
	return strings.TrimSuffix(fsPath.Last(), "/")
}

//
// Returns slice of path strings which include eachother
// with each iteration.
//
func (fsPath *FspopPath) PathProgressive() []string {
	if fsPath.IsEmpty() {
		return []string{}
	}

	path := []string{}
	for i, _ := range fsPath.Path {
		path = append(path, strings.Join(fsPath.Path[0:i+1], ""))
	}

	return path
}

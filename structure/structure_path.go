package structure

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
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

func (fsPath *FspopPath) ParentPath() string {
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
		message.Error(path)
		message.Error(fmt.Sprint(fsPath.Path))
		panic(errors.New("tried to append after a file in fspopPath"))
		//return
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

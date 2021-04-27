package structure

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/merrittcorp/fspop/message"
)

type FspopStructurePath struct {
	Path []string
}

func CreateFspopStructurePath(pathInit []string) *FspopStructurePath {
	return &FspopStructurePath{
		Path: pathInit,
	}
}

func (fsPath *FspopStructurePath) Actual() string {
	return strings.TrimSuffix(fsPath.ToString(), "/")
}

func (fsPath *FspopStructurePath) ParentPath() string {
	if fsPath.Length() < 2 {
		return ""
	}
	var final string
	for _, value := range fsPath.Path[:len(fsPath.Path)-1] {
		final = final + value
	}
	return strings.TrimSuffix(final, "/")
}

func (fsPath *FspopStructurePath) ToString() string {
	var final string
	for _, value := range fsPath.Path {
		final = final + value
	}
	return final
}

func (fsPath *FspopStructurePath) IsEmpty() bool {
	return len(fsPath.Path) == 0
}

func (fsPath *FspopStructurePath) First() string {
	if fsPath.IsEmpty() {
		return ""
	}
	return fsPath.Path[0]
}

func (fsPath *FspopStructurePath) Last() string {
	if fsPath.IsEmpty() {
		return ""
	}
	return fsPath.Path[len(fsPath.Path)-1]
}

func (fsPath *FspopStructurePath) Name() string {
	if fsPath.IsEmpty() {
		return ""
	}
	return strings.TrimSuffix(fsPath.Last(), "/")
}

func (fsPath *FspopStructurePath) Length() int {
	return len(fsPath.Path)
}

func (fsPath *FspopStructurePath) Append(path string) {
	// Can't append if a file has been reached
	if !fsPath.IsEmpty() && !IsDirectory(fsPath.Last()) {
		message.Error(path)
		message.Error(fmt.Sprint(fsPath.Path))
		panic(errors.New("tried to append after a file in fspopstructurePath"))
		//return
	}

	if IsDirectory(path) {
		// Fix directory
		path = StandardizeDirectory(path)
	}

	fsPath.Path = append(fsPath.Path, path)
}

package structure

import (
	"errors"
	"fmt"

	"gitlab.com/merrittcorp/fspop/message"
)

type FspopStructurePath struct {
	Path []string
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

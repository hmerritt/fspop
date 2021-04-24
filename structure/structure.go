package structure

import "strings"

type FspopStructure struct {
	Version   string
	Name      string
	Data      interface{}
	Dynamic   interface{}
	Structure interface{}
}

func IsDirectory(str string) bool {
	return (strings.HasPrefix(str, "/") || strings.HasSuffix(str, "/"))
}

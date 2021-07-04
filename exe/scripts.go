package exe

import (
	"fmt"
	"path/filepath"

	"gitlab.com/merrittcorp/fspop/parse"
)

// Map of regognized script extensions
var ScriptRecognize = map[string]struct{}{
	".bat":        {},
	".cmd":        {},
	".ps1":        {},
	".powershell": {},
	".sh":         {},
	".bash":       {},
}

// Check if a command is a script
// reutrns true + the path to
var ScriptExists = func(path string, entrypoint string) (bool, string) {
	if _, ok := ScriptRecognize[filepath.Ext(path)]; ok {
		// Path including entrypoint
		pathEntrypoint := filepath.Clean(fmt.Sprintf("%s/%s", entrypoint, path))
		if parse.FileExists(pathEntrypoint) {
			return true, pathEntrypoint
			// Path without entrypoint
		} else if parse.FileExists(path) {
			return true, path
		}
	}
	return false, ""
}

// Detect script runner via it's file extention
// e.g. script.sh -> /bin/sh
var ScriptCommandExe = func(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".bat":
		return "cmd"
	case ".cmd":
		return "cmd"
	case ".powershell":
		return "powershell"
	case ".ps1":
		return "powershell"
	case ".bash":
		return "bash"
	case ".sh":
		return "bash"
	}
	return ""
}

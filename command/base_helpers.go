package command

import (
	"fmt"
)

// Populate map of select flags (defaults to ALL flags)
//
// fspop commands can choose which flags they need
func GetFlagMap(which []string) *FlagMap {
	// Convert which slice to a map
	//
	// This improves performace as map lookups are O(1)
	whichMap := make(map[string]struct{}, len(which))
	for _, i := range which {
		whichMap[i] = struct{}{}
	}

	// Create flag map
	fm := make(FlagMap, len(which))

	addToMap := func(fl *Flag) {
		// Check if flag name exists in whichMap
		_, ok := whichMap[fl.Name]

		// Add to 'fm' if.
		// 'whichMap' map is empty,
		// or flag exists in map.
		if len(whichMap) == 0 || ok {
			fm[fl.Name] = fl
		}
	}

	addToMap(&flagStrict)
	addToMap(&flagForce)

	return &fm
}

// Detect long flags entered with one dash '-'
// and add a dash to prevent a panic when parsing
//
// -strict -> --strict
func flagSingleToDoubleDash(args []string) []string {
	for i, arg := range args {
		for _, fl := range FlagNames {
			if arg == fmt.Sprintf("-%s", fl) {
				args[i] = fmt.Sprintf("--%s", fl)
			}
		}
	}
	return args
}

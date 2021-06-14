package main

import (
	"fmt"

	"gitlab.com/merrittcorp/fspop/command"
	"gitlab.com/merrittcorp/fspop/version"
)

func main() {
	printTitle()

	command.Run()
}

// Prints fspop title + version info
func printTitle() {
	// Get version info
	versionStruct := version.GetVersion()

	// Check if dev release
	// Hide git-sha for main releases
	isDev := false
	if versionStruct.VersionPrerelease != "" {
		isDev = true
	}

	// Get full version string
	versionString := versionStruct.FullVersionNumber(isDev)

	fmt.Println(versionString)
	fmt.Println("(c) MerrittCorp. All rights reserved.")
	fmt.Println()
}

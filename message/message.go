package message

import (
	"fmt"

	"gitlab.com/merrittcorp/fspop/version"
)

func PrintTitle() {
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

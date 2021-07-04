package exe

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"gitlab.com/merrittcorp/fspop/ui"
)

// Execute command via a shell, and log the output
func Run(shell, command, entrypoint string) {
	UI := ui.GetUi()

	// Setup command to exec
	// Set exec directory to the entrypoint
	run := exec.Command(shell, command)
	run.Dir = entrypoint

	// Run command and print output
	out, err := run.CombinedOutput()
	if err != nil {
		UI.Error(fmt.Sprint(err))
	}

	UI.Output(strings.TrimSpace(string(out)))
	UI.Output("")
}

// Returns shell name for OS.
//
// E.G. windows -> powershell
func GetOsShell() string {
	switch runtime.GOOS {
	case "windows":
		if CommandExists("powershell") {
			return "powershell"
		} else {
			return "cmd"
		}
	case "linux":
		if CommandExists("bash") {
			return "bash"
		} else {
			return "sh"
		}
	default:
		return "sh"
	}
}

// Searches for an executable in the PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

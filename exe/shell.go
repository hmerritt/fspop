package exe

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"gitlab.com/merrittcorp/fspop/ui"
)

// Execute command via a shell, and log the output
func Run(shell, command, entrypoint string) error {
	UI := ui.GetUi()

	var run *exec.Cmd

	// Setup command to exec
	//
	// Bash and sh require '-c' option
	if shell == "bash" || shell == "sh" {
		run = exec.Command(shell, "-c", command)
	} else {
		run = exec.Command(shell, command)
	}

	// Set exec directory to the entrypoint
	run.Dir = entrypoint

	// Run command and print output
	out, err := run.CombinedOutput()

	shellPrint := func(s string) {
		UI.Output(fmt.Sprintf("  %s", ui.WrapAtLength(strings.TrimSpace(s), 2)))
	}

	shellPrint(string(out))

	if err != nil {
		shellPrint(fmt.Sprint(err))
		UI.Output("")
		return err
	}

	return nil
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

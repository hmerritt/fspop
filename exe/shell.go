package exe

import (
	"fmt"
	"os/exec"

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

	UI.Output(string(out))
}

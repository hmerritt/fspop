package exe

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"sync"

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

	run.Dir = entrypoint

	stdOut, stdErr, err := runGetStdPipes(run)
	if err != nil {
		return err
	}

	err = run.Start()
	if err != nil {
		return err
	}

	// Stream standard output live
	var wg sync.WaitGroup
	wg.Add(2)
	go runStreamStd(stdOut, UI.Output, &wg)
	go runStreamStd(stdErr, UI.Error, &wg)

	err = run.Wait()
	if err != nil {
		return err
	}

	wg.Wait()

	return nil
}

// Returns both stdout and stderr piped io.Reader
func runGetStdPipes(run *exec.Cmd) (io.Reader, io.Reader, error) {
	// Standard Output pipe
	stdout, err := run.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	// Standard Error pipe
	stderr, err2 := run.StderrPipe()
	if err2 != nil {
		return nil, nil, err2
	}

	// Combine outputs into one io.Reader
	return stdout, stderr, nil
}

func runStreamStd(reader io.Reader, output func(message string), wg *sync.WaitGroup) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		s := scanner.Text()
		output(fmt.Sprintf("    %s", ui.IndentString(s, 4)))
	}

	wg.Done()
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

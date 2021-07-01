package command

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/posener/complete"
	"gitlab.com/merrittcorp/fspop/ui"
)

// Maximum width of any line, in character count.
const maxLineLength = 75

// Slice of all flag names
var FlagNames = []string{flagStrict.Name, flagForce.Name}

// Slice of global flag names
var FlagNamesGlobal = []string{flagStrict.Name, flagForce.Name}

// Master command type which in present in all commands
//
// Used to standardize UI output
type BaseCommand struct {
	UI *ui.Ui
}

type FlagMap map[string]*Flag

type Flag struct {
	Name       string
	Usage      string
	Default    interface{}
	Value      interface{}
	Completion complete.Predictor
}

// flag --strict
//
// Stop after any errors when deploying
var flagStrict = Flag{
	Name:    "strict",
	Usage:   "Stop after any errors or warnings.",
	Default: false,
	Value:   false,
}

// flag --force
//
// Prevents CLI prompts asking confirmation
var flagForce = Flag{
	Name:    "force",
	Usage:   "Bypasses CLI prompts without asking for confirmation.",
	Default: false,
	Value:   false,
}

// Help builds usage string for all flags
func (fm *FlagMap) Help() string {
	var out bytes.Buffer

	for _, flag := range *fm {
		fmt.Fprintf(&out, "  --%s \n      %s\n\n", flag.Name, flag.Usage)
	}

	return strings.TrimRight(out.String(), "\n")
}

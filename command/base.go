package command

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/posener/complete"
	"gitlab.com/merrittcorp/fspop/ui"
)

// Slice of all flag names
var FlagNames = []string{flagStrict.Name, flagForce.Name}

// Slice of global flag names
var FlagNamesGlobal = []string{flagStrict.Name, flagForce.Name}

// Master command type which is present in all commands
//
// Used to standardize UI output
type BaseCommand struct {
	UI *ui.Ui
}

type Flag struct {
	Name       string
	Usage      string
	Default    interface{}
	Value      interface{}
	Completion complete.Predictor
}

type FlagMap map[string]*Flag

func (fm *FlagMap) Get(flagName string) *Flag {
	fl, ok := (*fm)[flagName]
	if ok {
		return fl
	}
	return nil
}

// Help builds usage string for all flags in a FlagMap
func (fm *FlagMap) Help() string {
	var out bytes.Buffer

	for _, flag := range *fm {
		fmt.Fprintf(&out, "  --%s \n      %s\n\n", flag.Name, flag.Usage)
	}

	return strings.TrimRight(out.String(), "\n")
}

// Parse CLI args to FlagMap
func (fm *FlagMap) Parse(UI *ui.Ui, args []string) []string {
	// Struct used to parse flags
	var opts struct {
		Strict bool `short:"s" long:"strict"`
		Force  bool `short:"f" long:"force"`
	}

	// Parse flags from `args'.
	args, err := flags.ParseArgs(&opts, flagSingleToDoubleDash(args))

	if err != nil {
		UI.Error("Unable to parse flag from the arguments entered '" + fmt.Sprint(args[0]) + "'")
		UI.Warn("Flags are entered with double dashes '--', for example '--strict'")
		os.Exit(1)
	}

	updateFmWithOps := func(flagName string, value interface{}) {
		// Check if flag name exists in fm
		_, ok := (*fm)[flagName]

		// Update 'fm' if flag exists in map.
		if ok {
			(*fm)[flagName].Value = value
		}
	}

	updateFmWithOps("strict", opts.Strict)
	updateFmWithOps("force", opts.Force)

	return args
}

// flag definitions

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

// +build tools

// This file ensures tool dependencies are kept in sync.
// This is the recommended way of doing this according to:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

//go:generate go install github.com/mitchellh/gox
import (
	_ "github.com/mitchellh/gox"
)

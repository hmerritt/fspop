package command

import "github.com/mitchellh/go-wordwrap"

// Wraps the given text to maxLineLength.
func wrapAtLength(s string) string {
	return wordwrap.WrapString(s, maxLineLength)
}

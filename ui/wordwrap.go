package ui

import (
	"bytes"
	"unicode"
)

// Maximum width of any line, in character count.
const MaxLineLength = 75

//
// Imported directly from, https://github.com/mitchellh/go-wordwrap/blob/v1.0.1/wordwrap.go
//
// Modified subtly to fit the needs of fspop
//

const nbsp = 0xA0

// WrapString wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func WrapString(s string, lim uint, indentSpaces int) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current uint
	var wordBuf, spaceBuf bytes.Buffer
	var wordBufLen, spaceBufLen uint

	for _, char := range s {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+spaceBufLen > lim {
					current = 0
				} else {
					current += spaceBufLen
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
				spaceBufLen = 0
			} else {
				current += spaceBufLen + wordBufLen
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}
			buf.WriteRune(char)

			// Add custom indent
			for i := 0; i < indentSpaces; i++ {
				buf.WriteRune(' ')
			}

			current = 0
		} else if unicode.IsSpace(char) && char != nbsp {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBufLen + wordBufLen
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}

			spaceBuf.WriteRune(char)
			spaceBufLen++
		} else {
			wordBuf.WriteRune(char)
			wordBufLen++

			if current+wordBufLen+spaceBufLen > lim && wordBufLen < lim {
				buf.WriteRune('\n') // break

				// Add custom indent
				for i := 0; i < indentSpaces; i++ {
					buf.WriteRune(' ')
				}

				current = 0
				spaceBuf.Reset()
				spaceBufLen = 0
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+spaceBufLen <= lim {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}

	return buf.String()
}

// Wraps the given text to maxLineLength.
//
// Define the number of indent spaces added on-wrap
func WrapAtLength(s string, indentSpaces int) string {
	return WrapString(s, MaxLineLength, indentSpaces)
}

// Indent string for ALL line-breaks.
func IndentString(s string, indentSpaces int) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	for _, char := range s {
		buf.WriteRune(char)

		if char == '\n' {
			// Add custom indent
			for i := 0; i < indentSpaces; i++ {
				buf.WriteRune(' ')
			}
		}
	}

	return buf.String()
}

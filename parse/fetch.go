package parse

import (
	"os"
)

func FetchFile(filepath string) []byte {
	// Read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	// Return byte data
	return data
}

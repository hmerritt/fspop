package parse

import (
	"net/http"
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

func FetchUrl(url string) []byte {
	// Fetch URL
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	content := make([]byte, res.ContentLength)
	_, err = res.Body.Read(content)
	if err != nil {
		panic(err)
	}

	return content
}

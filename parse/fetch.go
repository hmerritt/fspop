package parse

import (
	"net/http"
	"os"
	"strings"
)

func IsUrl(str string) bool {
	if strings.HasPrefix(str, "www.") || strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://") {
		return true
	}
	return false
}

func UseUrl(str string) bool {
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		return IsUrl(str)
	}

	return false
}

func FetchYaml(path string) []byte {
	if UseUrl(path) {
		return FetchUrl(path)
	} else {
		return FetchFile(path)
	}
}

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

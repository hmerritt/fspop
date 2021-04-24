package parse

import (
	"net/http"
	"os"
	"strings"
)

func IsUrl(path string) bool {
	return (strings.HasPrefix(path, "www.") || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"))
}

func UseUrl(path string) bool {
	return (!FileExists(path) && IsUrl(path))
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FetchYaml(path string) []byte {
	if UseUrl(path) {
		return FetchUrl(path)
	} else {
		return FetchFile(path)
	}
}

func FetchFile(filepath string) []byte {
	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	return data
}

func FetchUrl(url string) []byte {
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

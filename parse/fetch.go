package parse

import (
	"errors"
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

func FetchYaml(path string) ([]byte, error) {
	if UseUrl(path) {
		return FetchUrl(path)
	} else {
		return FetchFile(path)
	}
}

func FetchFile(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	return data, err
}

func FetchUrl(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(res.Status, "1") && !strings.HasPrefix(res.Status, "2") {
		return nil, errors.New("request returned a bad http status code: " + res.Status + ".")
	}

	content := make([]byte, res.ContentLength)
	_, err = res.Body.Read(content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

//func HttpStatusCodeError

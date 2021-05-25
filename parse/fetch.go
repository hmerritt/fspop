package parse

import (
	"errors"
	"os"
	"strings"

	"github.com/imroc/req"
)

func IsUrl(path string) bool {
	return (strings.HasPrefix(path, "www.") || strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://"))
}

func UseUrl(path string) bool {
	return (!FileExists(path) && IsUrl(path))
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && !stat.IsDir()
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
	res, err := req.Get(url)
	if err != nil {
		return nil, err
	}

	resStatus := res.Response().Status

	if !strings.HasPrefix(resStatus, "1") && !strings.HasPrefix(resStatus, "2") {
		return nil, errors.New("request returned a bad http status code: " + resStatus + ".")
	}

	return res.Bytes(), nil
}

// Create file and return open file object
func CreateFile(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	return file, err
}

//func HttpStatusCodeError

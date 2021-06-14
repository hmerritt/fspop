package parse

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/imroc/req"
	"gitlab.com/merrittcorp/fspop/ui"
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

// Fetch and parse YAML file locally or from a URL
func FetchYaml(path string) []byte {
	var data []byte
	var err error

	UI := ui.GetUi()

	// Decide if URL or file
	if UseUrl(path) {
		ui.Spinner.Start("", " Fetching URL data...")

		data, err = FetchUrl(path)

		ui.Spinner.Stop()

		if err != nil {
			UI.Error("Unable to fetch URL data.")
			UI.Error(fmt.Sprint(err))
			UI.Warn("\nMake sure the link is accessible and try again.")
			os.Exit(2)
		}
	} else {
		data, err = FetchFile(path)

		if err != nil {
			UI.Error("Unable to open file.")
			UI.Error(fmt.Sprint(err))
			UI.Warn("\nCheck the file is exists and try again.")
			os.Exit(2)
		}
	}

	return data
}

// Create file and return open file object
func CreateFile(path string) (*os.File, error) {
	file, err := os.Create(path)
	return file, err
}

//go:build mage

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/merrittcorp/fspop/version"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace
type Npm mg.Namespace

var Aliases = map[string]interface{}{
	"build": Build.Release,
}

// ----------------------------------------------------------------------------
// Test
// ----------------------------------------------------------------------------

// Runs Go tests
func Test() error {
	log := NewLogger()
	defer log.End()
	return RunSync([][]string{
		{"gotestsum", "--format", "pkgname", "--", "--cover", "./..."},
	})
}

func Bench() error {
	log := NewLogger()
	defer log.End()
	return RunSync([][]string{
		{"gotestsum", "--format", "pkgname", "--", "--cover", "-bench", ".", "-benchmem", "./..."},
	})
}

// ----------------------------------------------------------------------------
// Build
// ----------------------------------------------------------------------------

func (Build) Clean() error {
	log := NewLogger()
	defer log.End()
	log.Info("cleaning bin directory")

	if err := os.RemoveAll("bin"); err != nil {
		return log.Error(err)
	}

	return fs.WalkDir(os.DirFS("npm"), ".", func(filePath string, d fs.DirEntry, err error) error {
		if d.IsDir() ||
			filePath == "fspop/bin/fspop" ||
			!strings.HasPrefix(filePath, "fspop") ||
			!strings.HasPrefix(d.Name(), "fspop") {
			return nil
		}
		return os.Remove(path.Join("npm", filePath))
	})
}

func (Build) Debug() error {
	log := NewLogger()
	defer log.End()
	log.Info("compiling debug binary")
	return RunSync([][]string{
		{"go", "build", "-ldflags", "-s -w", "."},
	})
}

func (Build) Release() error {
	mg.Deps(Build.Clean, BumpVersion)

	log := NewLogger()
	defer log.End()
	log.Info("compiling release binary")
	return RunSync([][]string{
		{"gox",
			"-osarch",
			"darwin/amd64 darwin/arm64 freebsd/amd64 freebsd/arm linux/amd64 linux/arm64 netbsd/amd64 netbsd/arm openbsd/amd64 windows/amd64",
			"-gocmd",
			"go",
			"-ldflags",
			LdFlagString(),
			"-tags",
			"fspop",
			"-output",
			"bin/{{.OS}}_{{.Arch}}/fspop",
			"."},
	})
}

// ----------------------------------------------------------------------------
// Release
// ----------------------------------------------------------------------------

func Release() error {
	log := NewLogger()
	defer log.End()

	releaseVersion := GetEnv("RELEASE_VERSION", version.Version)
	log.Info("release version: ", releaseVersion)

	releaseArchs := []string{
		"darwin_amd64",
		"darwin_arm64",
		"freebsd_amd64",
		"freebsd_arm",
		"linux_amd64",
		"linux_arm64",
		"netbsd_amd64",
		"netbsd_arm",
		"openbsd_amd64",
		"windows_amd64"}

	releaseArchsNpmDirectories := []string{
		"fspop-darwin-x64",
		"fspop-darwin-arm64",
		"",
		"",
		"fspop-linux-x64",
		"fspop-linux-arm64",
		"",
		"",
		"",
		"fspop-win32-x64"}

	// Zip
	for archIndex, arch := range releaseArchs {
		log.Info(fmt.Sprintf("%13s", arch))
		archLog := fmt.Sprintf("%13s |", arch)
		binDirPath := fmt.Sprintf("bin/%s", arch)
		binFilePath := ""

		// Search each binary path for the release binary,
		// ensure binary file exists and check it's file size.
		log.Debug(archLog, "checking release binary")
		binPathfiles, err := os.ReadDir(binDirPath)
		if err != nil {
			return log.Error(archLog, "error reading directory:", err)
		}
		for _, file := range binPathfiles {
			if file.IsDir() {
				continue
			}
			if strings.Contains(file.Name(), "fspop") {
				fileInfo, err := file.Info()
				if err != nil {
					return log.Error(archLog, "error getting file info:", err)
				}
				if fileInfo.Size() < 1000000 {
					return log.Error(archLog, "release binary is too small:", err)
				}
				binFilePath = fmt.Sprintf("%s/%s", binDirPath, file.Name())
				break
			}
		}
		if binFilePath == "" {
			return log.Error(archLog, "failed to find release binary", arch)
		}

		log.Debug(archLog, "zip for release")
		zipFileName := fmt.Sprintf("fspop_%s_%s.zip", releaseVersion, arch)
		zipFilePath := fmt.Sprintf("bin/%s", zipFileName)

		err = ZipFiles(zipFilePath, binFilePath)
		if err != nil {
			return log.Error(archLog, "failed to zip binary", err)
		}

		// Copy binary into npm directory
		if releaseArchsNpmDirectories[archIndex] != "" {
			log.Debug(archLog, "copy binary into npm directory")
			CopyFile(binFilePath, path.Join("npm", releaseArchsNpmDirectories[archIndex], path.Base(binFilePath)))
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Housekeeping
// ----------------------------------------------------------------------------

// Bootstraps required packages (installs required linux/macOS packages if needed)
func Bootstrap() error {
	log := NewLogger()
	defer log.End()

	// Install mage bootstrap (the recommended, as seen in https://magefile.org)
	if !ExecExists("mage") && ExecExists("git") {
		log.Info("installing mage")
		tmpDir := "__tmp_mage"

		if err := sh.Run("git", "clone", "https://github.com/magefile/mage", tmpDir); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		if err := os.Chdir(tmpDir); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		if err := sh.Run("go", "run", "bootstrap.go"); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		if err := os.Chdir("../"); err != nil {
			return log.Error("error: installing mage: ", err)
		}

		os.RemoveAll(tmpDir)
	}

	// Install Go dependencies
	log.Info("installing go dependencies")
	return RunSync([][]string{
		{"go", "mod", "vendor"},
		{"go", "mod", "tidy"},
		{"go", "generate", "-tags", "tools", "tools/tools.go"},
	})
}

// Update all Go dependencies
func UpdateDeps() error {
	log := NewLogger()
	defer log.End()
	return RunSync([][]string{
		{"go", "get", "-u", "all"},
		{"go", "mod", "tidy"},
	})
}

// Bumps (syncs) patch version to the commit count (see `version/version_base.go`)
func BumpVersion() error {
	log := NewLogger()
	defer log.End()

	// Get the total commit count
	commitCountString, err := sh.Output("git", "rev-list", "--count", "HEAD")

	if err != nil {
		return log.Error("failed to get commit count:", err)
	}

	commitCount, err := strconv.Atoi(commitCountString)

	if err != nil {
		return log.Error("failed to parse commit count:", commitCountString, err)
	}

	filesWithVersion := []string{
		"version/version_base.go",
		"npm/fspop/package.json",
		"npm/fspop-darwin-arm64/package.json",
		"npm/fspop-darwin-x64/package.json",
		"npm/fspop-linux-arm64/package.json",
		"npm/fspop-linux-x64/package.json",
		"npm/fspop-win32-x64/package.json",
	}

	for index, versionFile := range filesWithVersion {
		versionFileContent, err := os.ReadFile(versionFile)

		if err != nil {
			return log.Error("failed open version file:", err)
		}

		// Change regex based on file type
		var versionMatchRegex, versionReplaceRegex, versionReplaceSprintf string
		switch path.Ext(versionFile) {
		case ".go":
			versionMatchRegex = `= "(\d+).(\d+).(\d+)"`
			versionReplaceRegex = `= "\d+.\d+.\d+"`
			versionReplaceSprintf = `= "%s.%s.%d"`
		case ".json":
			versionMatchRegex = `"version": "(\d+).(\d+).(\d+)"`
			versionReplaceRegex = `"version": "\d+.\d+.\d+"`
			versionReplaceSprintf = `"version": "%s.%s.%d"`
		}

		// Extract current patch version
		versionMatch := regexp.MustCompile(versionMatchRegex).FindStringSubmatch(string(versionFileContent))

		if len(versionMatch) < 4 {
			return log.Error("failed to parse version:", versionFile, versionMatch)
		}

		majorVersion := versionMatch[1]
		minorVersion := versionMatch[2]
		patchVersionCurrent := versionMatch[3]
		patchVersion := commitCount

		versionCurrent := fmt.Sprintf("%s.%s.%s", majorVersion, minorVersion, patchVersionCurrent)
		versionNew := fmt.Sprintf("%s.%s.%d", majorVersion, minorVersion, patchVersion)

		if index == 0 {
			log.Info("bumping version", versionCurrent, "->", versionNew)
		}

		// Update version file
		versionFileContent = regexp.MustCompile(versionReplaceRegex).ReplaceAll(versionFileContent, []byte(fmt.Sprintf(versionReplaceSprintf, majorVersion, minorVersion, patchVersion)))

		if err := os.WriteFile(versionFile, versionFileContent, 0644); err != nil {
			return log.Error("failed to write version file:", err)
		}
	}

	return nil
}

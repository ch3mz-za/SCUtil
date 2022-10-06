package common

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func CleanInput(input string) string {
	os := runtime.GOOS
	switch os {
	case "windows":
		return strings.Replace(input, "\r\n", "", -1)
	case "darwin":
		return strings.Replace(input, "\n", "", -1)
	default:
		return input
	}
}

func FindDir(root, target string) (string, error) {

	var gamePath string
	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if dir.IsDir() && filepath.Base(path) == target {
			gamePath = path
			return nil
		}
		return nil
	})
	return gamePath, err
}

func ListAllFilesAndDirs(dir string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func UserHomeDir() string {
	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	return home
}

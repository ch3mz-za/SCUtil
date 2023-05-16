package common

import (
	"fmt"
	"io"
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

		if dir.IsDir() && strings.HasSuffix(path, target) {
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

func CopyFile(src string, dst string) error {
	fin, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		return err
	}
	return nil
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to create directory: %s - Error: %s", dir, err.Error())
		}
	}
}

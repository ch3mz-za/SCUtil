package common

import (
	"fmt"
	"io"
	"io/fs"
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
	target = filepath.Clean(target)
	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		cleanPath := filepath.Clean(path)
		if dir.IsDir() && (cleanPath == target || strings.HasSuffix(cleanPath, string(filepath.Separator)+target)) {
			gamePath = path
			return fs.SkipDir
		}
		return nil
	})
	return gamePath, err
}

func ListAllFilesAndDirs(dir string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(dir)
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
	if home == "" {
		home, _ = os.UserHomeDir()
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
	_, err := os.Stat(path)
	return err == nil
}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to create directory: %s - Error: %s", dir, err.Error())
		}
	}
}

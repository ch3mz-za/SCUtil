package util

import (
	"fmt"
	"os"
)

func Ptr[T any](v T) *T {
	return &v
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("unable to create directory: %s", dir)
		}
	}
	return nil
}

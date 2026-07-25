package scu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ch3mz-za/SCUtil/internal/common"
)

// ClearUserFolder - Clears all the data in the USER folder with the option to exclude control mappings
func ClearUserFolder(version string, exclusionsEnabled bool) error {
	userDir := filepath.Join(GetGameDir(), version, "USER")
	exclusion := filepath.Join(userDir, "Client", "0", "Controls")

	err := filepath.Walk(userDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {

				if exclusionsEnabled && isWithin(path, exclusion) {
					return nil
				}

				if err := os.Remove(path); err != nil {
					return fmt.Errorf("remove %s: %w", path, err)
				}
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("error removing USER directory:\n %s", err.Error())
	}
	return nil
}

// ClearAllDataExceptP4k - Clears all the data around the Data.p4k file
func ClearAllDataExceptP4k(version string) error {
	gameDir := filepath.Join(GetGameDir(), version)

	files, err := common.ListAllFilesAndDirs(gameDir)
	if err != nil {
		return fmt.Errorf("unable to list directories and files in %s: %w", gameDir, err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".p4k") {
			continue
		}
		filePath := filepath.Join(gameDir, f.Name())

		if err := os.RemoveAll(filePath); err != nil {
			return fmt.Errorf("remove %s: %w", filePath, err)
		}
	}
	return nil
}

func isWithin(path, root string) bool {
	rel, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	return err == nil && rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// ClearEasyAntiCheatFolder - Clears the EasyAntiCheat folder for the specified version
func ClearEasyAntiCheatFolder(version string) error {
	easyAntiCheatDir := filepath.Join(GetGameDir(), version, "EasyAntiCheat")
	if common.Exists(easyAntiCheatDir) {
		if err := os.RemoveAll(easyAntiCheatDir); err != nil {
			return err
		}
	}
	return nil
}

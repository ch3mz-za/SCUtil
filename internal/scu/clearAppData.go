package scu

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/ch3mz-za/SCUtil/internal/common"
)

// ClearStarCitizenAppData - Clears the game's date within AppData
func ClearStarCitizenAppData(enableExclusions bool) error {
	scAppDataDir := filepath.Join(common.UserHomeDir(), "AppData", "Local", "Star Citizen")
	files, _ := common.ListAllFilesAndDirs(scAppDataDir)

	var exclusion []string
	if enableExclusions {
		exclusion = append(exclusion, "GraphicsSettings.json")
	}

	if len(files) == 0 {
		return errors.New("Star Citizen AppData is empty")
	} else {
		for _, f := range files {
			if err := deleteAllFilesWithExclusions(filepath.Join(scAppDataDir, f.Name()), exclusion...); err != nil {
				return err
			}
		}
	}
	return nil
}

// ClearRsiLauncherAppData - Clears the game's launcher data within AppData
func ClearRsiLauncherAppData() *[]string {
	var filesRemoved []string
	for _, folder := range []string{"rsilauncher", "RSI Launcher"} {
		rsiLauncherDir := filepath.Join(common.UserHomeDir(), "AppData", "Roaming", folder)
		files, _ := common.ListAllFilesAndDirs(rsiLauncherDir)
		for _, f := range files {
			if err := os.RemoveAll(filepath.Join(rsiLauncherDir, f.Name())); err != nil {
				continue
			}
			filesRemoved = append(filesRemoved, f.Name())
		}
	}
	return &filesRemoved
}

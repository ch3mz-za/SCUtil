package scu

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/pkg/common"
	p4k "github.com/ch3mz-za/SCUtil/pkg/p4kReader"
)

const (
	p4kFileNameDir   string        = "./p4k_filenames"
	p4kSearchResults string        = "./p4k_search_results"
	ctrlMapFileExt   string        = ".xml"
	twoSecondDur     time.Duration = 2 * time.Second

	// Directories
	UserDir                  string = "USER"
	UserBackupDir            string = "BACKUPS/UserFolder"
	ControlMappingsDir       string = UserDir + "/Client/0/Controls/Mappings"
	ControlMappingsBackupDir string = "BACKUPS/ControlMappings"
	ScreenshotsDir           string = "ScreenShots"
	ScreenshotsBackupDir     string = "BACKUPS/Screenshots"
)

var (
	GameDir string = ""
	AppDir  string = ""
)

// ClearAllDataExceptP4k - Clears all the data around the Data.p4k file
func ClearAllDataExceptP4k(version string) error {
	gameDir := filepath.Join(GameDir, version)

	files, err := common.ListAllFilesAndDirs(gameDir)
	if err != nil {
		return errors.New("unable to list directories and files")
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".p4k") {
			continue
		}
		fPAth := filepath.Join(gameDir, f.Name())

		err := os.RemoveAll(fPAth)
		if err != nil {
			continue
		}
	}
	return nil
}

// ClearUserFolder - Clears all the data in the USER folder with the option to exclude control mappings
func ClearUserFolder(version string, exclusionsEnabled bool) error {
	userDir := filepath.Join(GameDir, version, "USER")
	exclusion := filepath.Join(userDir, "Client", "0", "Controls")

	err := filepath.Walk(userDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {

				if exclusionsEnabled && strings.HasPrefix(path, exclusion) {
					return nil
				}

				if err := os.Remove(path); err != nil {
					return nil
				}
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("error removing USER directory:\n %s", err.Error())
	}
	return nil
}

// GetP4kFilenames - Gets all the filenames from the Data.p4k file and writes them to a specific folder
func GetP4kFilenames(version string) error {
	gameDir := filepath.Join(GameDir, version)
	p4k.GetP4kFilenames(gameDir, p4kFileNameDir)
	return nil
}

// SearchP4kFilenames - Search for specific filenames within the Data.p4k file
func SearchP4kFilenames(version, phrase string) error {
	gameDir := filepath.Join(GameDir, version)
	results, err := p4k.SearchP4kFilenames(gameDir, phrase)
	if err != nil {
		return fmt.Errorf("unable to search files: %s", err.Error())
	}

	filename := strings.ReplaceAll(phrase, "\\", "_") + ".txt"
	p4k.MakeDir(p4kSearchResults)
	p4k.WriteStringsToFile(filepath.Join(p4kSearchResults, filename), results)
	return nil
}

// ClearStarCitizenAppData - Clears the game's date within AppData
func ClearStarCitizenAppData() (*[]string, error) {
	scAppDataDir := filepath.Join(common.UserHomeDir(), "AppData", "Local", "Star Citizen")
	files, _ := common.ListAllFilesAndDirs(scAppDataDir)
	filesRemoved := make([]string, 0, len(files))

	if len(files) == 0 {
		return &filesRemoved, errors.New("Star Citizen AppData is empty")
	} else {
		for _, f := range files {
			filename := filepath.Join(scAppDataDir, f.Name())
			err := os.RemoveAll(filename)
			if err != nil {
				continue
			}
			filesRemoved = append(filesRemoved, f.Name())
		}
	}
	return &filesRemoved, nil
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

// BackupControlMappings - Backup game control mappings
func BackupControlMappings(version string) error {
	mappingsDir := filepath.Join(GameDir, version, ControlMappingsDir)
	backupDir := filepath.Join(AppDir, ControlMappingsBackupDir, version)
	if err := backupFiles(mappingsDir, backupDir, true, ctrlMapFileExt); err != nil {
		return err
	}
	return nil
}

// GetBackedUpControlMappings - Retrieve a list of all the backed-up control mappings
func GetBackedUpControlMappings(version string) (*[]string, error) {
	backupDir := filepath.Join(AppDir, ControlMappingsBackupDir, version)
	files, err := os.ReadDir(backupDir)
	if err != nil {
		return &[]string{}, errors.New("unable to open backup directory")
	}

	items := make([]string, 0, len(files))
	for _, f := range files {
		items = append(items, f.Name())
	}
	return &items, nil
}

// RestoreControlMappings - Restores a specified control mapping for a specific game version
func RestoreControlMappings(version string, filename string) error {
	mappingsDir := filepath.Join(GameDir, version, ControlMappingsDir)
	backupDir := filepath.Join(AppDir, ControlMappingsBackupDir, version)
	return restoreFile(backupDir, mappingsDir, filename)
}

// BackupScreenshots - Backup all screenshots for specific game version
func BackupScreenshots(version string) error {
	screenshotDir := filepath.Join(GameDir, version, ScreenshotsDir)
	backupDir := filepath.Join(AppDir, ScreenshotsBackupDir, version)
	return backupFiles(screenshotDir, backupDir, false, ".jpg")
}

// BackupUserDirectory - Backup the USER directory
func BackupUserDirectory(version string) error {
	userDir := filepath.Join(GameDir, version, UserDir)
	backupDir := filepath.Join(AppDir, UserBackupDir, version, UserDir)
	return BackupDirectory(userDir, backupDir)
}

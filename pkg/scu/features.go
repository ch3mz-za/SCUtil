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
	ctrlMapFileExt string        = ".xml"
	twoSecondDur   time.Duration = 2 * time.Second

	// Directories
	UserDir                  string = "USER"
	UserBackupDir            string = "BACKUPS/UserFolder"
	ControlMappingsDir       string = UserDir + "/Client/0/Controls/Mappings"
	ControlMappingsBackupDir string = "BACKUPS/ControlMappings"
	ScreenshotsDir           string = "ScreenShots"
	ScreenshotsBackupDir     string = "BACKUPS/Screenshots"
	CharactersDir            string = UserDir + "/Client/0/CustomCharacters"
	CharactersBackupDir      string = "BACKUPS/CustomCharacters"
	P4kSearchResultsDir      string = "P4kResults/Searches"
	P4kFilenameResultsDir    string = "P4kResults/AllFileNames/%s/AllP4kFilenames.txt"
)

var (
	GameDir string = ""
	AppDir  string = ""
)

func GetGameVersions() []string {
	dirs, err := os.ReadDir(GameDir)
	if err != nil {
		return []string{}
	}

	versions := make([]string, 0, len(dirs))
	for _, d := range dirs {
		versions = append(versions, d.Name())
	}
	return versions
}

// ClearAllDataExceptP4k - Clears all the data around the Data.p4k file
func ClearAllDataExceptP4k(version string) error {
	gameDir := filepath.Join(GameDir, version)

	files, err := common.ListAllFilesAndDirs(gameDir)
	if err != nil {
		return fmt.Errorf("unable to list directories and files in %s: %w", gameDir, err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".p4k") {
			continue
		}
		filePath := filepath.Join(gameDir, f.Name())

		err := os.RemoveAll(filePath)
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
	resultsDir := filepath.Join(AppDir, fmt.Sprintf(P4kFilenameResultsDir, version))
	return p4k.GetP4kFilenames(gameDir, resultsDir)
}

// SearchP4kFilenames - Search for specific filenames within the Data.p4k file
func SearchP4kFilenames(version, phrase string) error {
	gameDir := filepath.Join(GameDir, version)
	filename := strings.ReplaceAll(phrase, "\\", "_") + ".txt"
	resultsDir := filepath.Join(AppDir, P4kSearchResultsDir, version)
	common.MakeDir(resultsDir)
	resultsDir = filepath.Join(resultsDir, filename)

	err := p4k.SearchP4kFilenames(gameDir, phrase, resultsDir)
	if err != nil {
		return err
	}

	return nil
}

// ClearStarCitizenAppData - Clears the game's date within AppData
func ClearStarCitizenAppData() error {
	scAppDataDir := filepath.Join(common.UserHomeDir(), "AppData", "Local", "Star Citizen")
	files, _ := common.ListAllFilesAndDirs(scAppDataDir)

	if len(files) == 0 {
		return errors.New("Star Citizen AppData is empty")
	} else {
		for _, f := range files {
			if err := deleteAllFilesWithExclusions(filepath.Join(scAppDataDir, f.Name()), "GraphicsSettings.json"); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteAllFilesWithExclusions(dir string, exclusions ...string) error {
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				for _, ex := range exclusions {
					if strings.HasSuffix(path, ex) {
						return nil
					}
				}

				if err := os.Remove(path); err != nil {
					return nil
				}
			}

			return nil
		})
	if err != nil {
		return fmt.Errorf("error removing directory '%s':\n %s", dir, err.Error())
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

// BackupControlMappings - Backup game control mappings
func BackupControlMappings(version string) error {
	mappingsDir := filepath.Join(GameDir, version, ControlMappingsDir)
	backupDir := filepath.Join(AppDir, ControlMappingsBackupDir, version)
	return backupFiles(mappingsDir, backupDir, true, ctrlMapFileExt)
}

// GetFilesListFromDir - Retrieve a list of all the files listed at a directory
func GetFilesListFromDir(dir string) (*[]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return &[]string{}, errors.New("no backup directory found")
	}

	items := make([]string, 0, len(files))
	for _, f := range files {
		items = append(items, f.Name())
	}
	return &items, nil
}

// GetBackedUpControlMappings - Retrieve a list of all the backed-up control mappings
func GetBackedUpControlMappings(version string) (*[]string, error) {
	return GetFilesListFromDir(filepath.Join(AppDir, ControlMappingsBackupDir, version))
}

// RestoreControlMappings - Restores a specified control mapping for a specific game version
func RestoreControlMappings(version string, filename string) error {
	mappingsFilePath := filepath.Join(GameDir, version, ControlMappingsDir, filename)
	backupFilePath := filepath.Join(AppDir, ControlMappingsBackupDir, version, filename)
	return restoreFile(backupFilePath, mappingsFilePath, true)
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

// BackupUserCharacters - Backup the custom characters in the USER directory
func BackupUserCharacters(version string) error {
	charDir := filepath.Join(GameDir, version, CharactersDir)
	backupDir := filepath.Join(AppDir, CharactersBackupDir, version)
	return backupFiles(charDir, backupDir, false, ".chf")
}

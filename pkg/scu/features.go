package scu

import (
	"errors"
	"fmt"
	"log"
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
	controlMappingsDir       string = "USER/Client/0/Controls/Mappings"
	controlMappingsBackupDir string = "BACKUPS/ControlMappings"
	screenshotsDir           string = "ScreenShots"
	screenshotsBackupDir     string = "BACKUPS/Screenshots"
)

var (
	RootDir string = ""
)

// func GetGameDir() (string, error) {
// 	return common.FindDir(RootDir, string(version))
// }

func ClearAllDataExceptP4k(version string) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("unable to find game directory")
	}

	// TODO: Check how you can update the status line
	fmt.Printf("\nGame directory found: %s\n", gameDir)

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
			fmt.Println("ERROR: " + err.Error())
			continue
		}
		// TODO: Remove println like this
		println("Deleted: " + fPAth)

	}
	return nil
}

func ClearUserFolder(version string, exclusionsEnabled bool) error {
	userDir, err := common.FindDir(RootDir, string(version))
	if err != nil || userDir == "" {
		return errors.New("unable to find game directory")
	}
	exclusion := filepath.Join(userDir, "USER", "Client", "0", "Controls")

	err = filepath.Walk(filepath.Join(userDir, "USER"),
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

func GetP4kFilenames(version string) error {

	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("unable to find game directory")
	}

	p4k.GetP4kFilenames(gameDir, p4kFileNameDir)
	return nil
}

func SearchP4kFilenames(version, phrase string) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("unable to find game directory")
	}

	results, err := p4k.SearchP4kFilenames(gameDir, phrase)
	if err != nil {
		return fmt.Errorf("unable to search files: %s", err.Error())
	}

	filename := strings.ReplaceAll(phrase, "\\", "_") + ".txt"
	p4k.MakeDir(p4kSearchResults)
	p4k.WriteStringsToFile(filepath.Join(p4kSearchResults, filename), results)
	return nil
}

func ClearStarCitizenAppData() (*[]string, error) {
	var filesRemoved []string

	scAppDataDir := filepath.Join(common.UserHomeDir(), "AppData", "Local", "Star Citizen")
	files, _ := common.ListAllFilesAndDirs(scAppDataDir)

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

func BackupControlMappings(version string) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("unable to find game directory")
	}
	log.Printf("\nUser directory found: %s\n", gameDir)
	mappingsDir := filepath.Join(gameDir, controlMappingsDir)
	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), controlMappingsBackupDir, string(version))
	if err := backupFiles(mappingsDir, backupDir, true, ".xml"); err != nil {
		return err
	}
	return nil
}

func GetBackedUpControlMappings(version string) (*[]string, error) {

	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return nil, errors.New("unable to find game directory")
	}

	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), controlMappingsBackupDir, string(version))

	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, errors.New("unable to open backup directory")
	}

	items := make([]string, 0, len(files))
	for _, f := range files {
		items = append(items, f.Name())
	}
	return &items, nil
}

func RestoreControlMappings(version string, filename string) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("unable to find game directory")
	}
	mappingsDir := filepath.Join(gameDir, controlMappingsDir)
	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), controlMappingsBackupDir, string(version))
	return restoreFiles(backupDir, mappingsDir, filename)
}

func BackupScreenshots(version string) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("unable to find game directory")
	}

	fmt.Printf("\nGame directory found: %s\n", gameDir)
	screenshotDir := filepath.Join(gameDir, screenshotsDir)
	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), screenshotsBackupDir, string(version))
	backupFiles(screenshotDir, backupDir, false, ".jpg")
	return nil
}

package scu

import (
	"bufio"
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

func ClearAllDataExceptP4k(version GameVersion) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("Unable to find game directory")
	}

	// TODO: Check how you can update the status line
	fmt.Printf("\nGame directory found: %s\n", gameDir)

	files, err := common.ListAllFilesAndDirs(gameDir)
	if err != nil {
		return errors.New("Unable to list directories and files")
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
		println("Deleted: " + fPAth)

	}
	return nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ClearUserFolder(version GameVersion, exclusionsEnabled bool) error {
	userDir, err := common.FindDir(RootDir, string(version))
	if err != nil || userDir == "" {
		return errors.New("Unable to find game directory")
	}
	exclusion := filepath.Join(userDir, "USER", "Client", "0", "Controls")
	fmt.Printf("\nUser directory found: %s\n", userDir)

	err = filepath.Walk(filepath.Join(userDir, "USER"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {

				if exclusionsEnabled && strings.HasPrefix(path, exclusion) {
					fmt.Println("Excluding: " + path)
					return nil
				}

				if err := os.Remove(path); err != nil {
					fmt.Printf("Unable to remove file: %s | error: %s\n", path, err.Error())
					return nil
				}
				fmt.Println("Removing: " + path)

			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("Error removing USER directory: %s\n", err.Error())
	}

	fmt.Println("Cleared USER directory")
	return nil
}

func GetP4kFilenames(version GameVersion) error {

	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("Unable to find game directory")
	}

	fmt.Printf("\nGame directory found: %s\n", gameDir)
	p4k.GetP4kFilenames(gameDir, p4kFileNameDir)
	return nil
}

func SearchP4kFilenames(version GameVersion) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return errors.New("Unable to find game directory")
	}

	reader := bufio.NewReader(os.Stdin)
	phrase, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("Invalid input: %s\n", err.Error())
	}

	// TODO: Fix this for GUI
	phrase = "kaas"
	results, err := p4k.SearchP4kFilenames(gameDir, phrase)
	if err != nil {
		return fmt.Errorf("Unable to search files: %s\n", err.Error())
	}

	println()

	filename := strings.ReplaceAll(phrase, "\\", "_") + ".txt"
	p4k.MakeDir(p4kSearchResults)
	p4k.WriteStringsToFile(filepath.Join(p4kSearchResults, filename), results)
	return nil
}

func ClearStarCitizenAppData() {
	scAppDataDir := filepath.Join(common.UserHomeDir(), "AppData", "Local", "Star Citizen")
	files, _ := common.ListAllFilesAndDirs(scAppDataDir)

	if len(files) == 0 {
		fmt.Println("Star Citizen AppData is empty")
	} else {
		fmt.Println("Clearing Star Citizen AppData directory")
		for _, f := range files {
			filename := filepath.Join(scAppDataDir, f.Name())
			err := os.RemoveAll(filename)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				continue
			}
			println("Deleted: " + filename)
		}
	}
}

func ClearRsiLauncherAppData() {

	for _, folder := range []string{"rsilauncher", "RSI Launcher"} {
		rsiLauncherDir := filepath.Join(common.UserHomeDir(), "AppData", "Roaming", folder)
		files, _ := common.ListAllFilesAndDirs(rsiLauncherDir)

		if len(files) == 0 {
			fmt.Printf("RSI Launcher AppData folder (%s) is empty!\n", folder)
			time.Sleep(twoSecondDur)
		} else {

			fmt.Println("Clearing Star Citizen AppData directory")
			for _, f := range files {
				filename := filepath.Join(rsiLauncherDir, f.Name())
				err := os.RemoveAll(filename)
				if err != nil {
					fmt.Println("ERROR: " + err.Error())
					continue
				}
				println("Deleted: " + filename)
			}
		}
	}
}

func BackupControlMappings(version GameVersion) error {
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

func GetBackedUpControlMappings(version GameVersion) (*[]string, error) {
	var items []string
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		return nil, errors.New("unable to find game directory")
	}

	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), controlMappingsBackupDir, string(version))

	files, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, errors.New("unable to open backup directory: " + err.Error())
	}

	for _, f := range files {
		items = append(items, f.Name())
	}
	return &items, nil
}

// func RestoreControlMappings(version GameVersion) error {
// 	gameDir, err := common.FindDir(RootDir, string(version))
// 	if err != nil || gameDir == "" {
// 		return errors.New("Unable to find game directory")
// 	}
// 	mappingsDir := filepath.Join(gameDir, controlMappingsDir)
// 	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), controlMappingsBackupDir, string(version))

// 	files, err := os.ReadDir(backupDir)
// 	if err != nil {
// 		return errors.New("Unable to open backup directory: " + err.Error())
// 	}

// 	var restoreMenuOpts []disp.MenuStringOption
// 	for _, f := range files {
// 		if strings.HasSuffix(strings.ToLower(f.Name()), ctrlMapFileExt) {
// 			restoreMenuOpts = append(restoreMenuOpts, disp.MenuStringOption(f.Name()))
// 		}
// 	}

// 	if len(restoreMenuOpts) == 0 {
// 		return errors.New("No objects found for restoration")
// 	}

// 	restoreMenuOpts = append(restoreMenuOpts, optBack)
// 	menuOption := disp.NewStringOptionMenu("Select file to restore", restoreMenuOpts).Run()
// 	if menuOption == optBack {
// 		return
// 	}

// 	restoreFiles(backupDir, mappingsDir, string(menuOption))
// }

func BackupScreenshots(version GameVersion) error {
	gameDir, err := common.FindDir(RootDir, string(version))
	if err != nil || gameDir == "" {
		errors.New("Unable to find game directory")

	}

	fmt.Printf("\nGame directory found: %s\n", gameDir)
	screenshotDir := filepath.Join(gameDir, screenshotsDir)
	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), screenshotsBackupDir, string(version))
	backupFiles(screenshotDir, backupDir, false, ".jpg")
	return nil
}

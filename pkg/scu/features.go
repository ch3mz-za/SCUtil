package scu

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/pkg/common"
	disp "github.com/ch3mz-za/SCUtil/pkg/display"
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

func ClearAllDataExceptP4k() {
	choice := ptuOrLiveMenu.Run()
	if choice == optBack {
		return
	}

	gameDir, err := common.FindDir(RootDir, string(choice))
	if err != nil || gameDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}

	fmt.Printf("\nGame directory found: %s\n", gameDir)

	files, err := common.ListAllFilesAndDirs(gameDir)
	if err != nil {
		fmt.Println("Unable to list directories and files")
		disp.EnterToContinue()
		return
	}

	if disp.YesOrNo("Clear all data except p4k file", true) {
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
		disp.EnterToContinue()
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ClearUserFolerWithExclusions() {
	ClearUserFolder(true)
}

func ClearUserFolerWithoutExclusions() {
	ClearUserFolder(false)
}

func ClearUserFolder(exclusionsEnabled bool) {
	choice := ptuOrLiveMenu.Run()
	if choice == optBack {
		return
	}
	userDir, err := common.FindDir(RootDir, string(choice))
	if err != nil || userDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}
	exclusion := filepath.Join(userDir, "USER", "Client", "0", "Controls")
	fmt.Printf("\nUser directory found: %s\n", userDir)

	if disp.YesOrNo("Clear user folder", true) {
		err := filepath.Walk(filepath.Join(userDir, "USER"),
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
			fmt.Printf("Error removing USER directory: %s\n", err.Error())
		}

		fmt.Println("Cleared USER directory")

	}
	disp.EnterToContinue()
}

func GetP4kFilenames() {
	choice := ptuOrLiveMenu.Run()
	if choice == optBack {
		return
	}
	gameDir, err := common.FindDir(RootDir, string(choice))
	if err != nil || gameDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}

	fmt.Printf("\nGame directory found: %s\n", gameDir)
	p4k.GetP4kFilenames(gameDir, p4kFileNameDir)
	disp.EnterToContinue()
}

func SearchP4kFilenames() {
	choice := ptuOrLiveMenu.Run()
	if choice == optBack {
		return
	}
	gameDir, err := common.FindDir(RootDir, string(choice))
	if err != nil || gameDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}

	disp.ClearTerminal()
	fmt.Print("Please enter your search phrase\n0. Back\n-> ")

	reader := bufio.NewReader(os.Stdin)
	phrase, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Invalid input: %s\n", err.Error())
		disp.EnterToContinue()
	}
	phrase = common.CleanInput(phrase)
	if phrase == "0" {
		return
	}

	results, err := p4k.SearchP4kFilenames(gameDir, phrase)
	if err != nil {
		fmt.Printf("Unable to search files: %s\n", err.Error())
		disp.EnterToContinue()
	}

	println()

	filename := strings.ReplaceAll(phrase, "\\", "_") + ".txt"
	p4k.MakeDir(p4kSearchResults)
	p4k.WriteStringsToFile(filepath.Join(p4kSearchResults, filename), results)
	disp.EnterToContinue()
}

func ClearStarCitizenAppData() {
	if !disp.YesOrNo("Clear Star Citizen App Data", true) {
		return
	}

	disp.ClearTerminal()
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
				fmt.Println("ERROR: " + err.Error())
				continue
			}
			println("Deleted: " + filename)
		}
	}
	disp.EnterToContinue()
}

func ClearRsiLauncherAppData() {

	if !disp.YesOrNo("Clear RSI Launcher data", true) {
		return
	}

	disp.ClearTerminal()
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
	disp.EnterToContinue()
}

func BackupOrRestoreControlMappings() {
	choice := ptuOrLiveMenu.Run()
	if choice == optBack {
		return
	}
	gameDir, err := common.FindDir(RootDir, string(choice))
	if err != nil || gameDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}

	disp.ClearTerminal()
	fmt.Printf("\nGame directory found: %s\n", gameDir)
	mappingsDir := filepath.Join(gameDir, controlMappingsDir)
	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), controlMappingsBackupDir, string(choice))

	var backupOpt disp.MenuStringOption = "Backup"
	var restoreOpt disp.MenuStringOption = "Restore"

	option := disp.NewStringOptionMenu("Control Mappings", []disp.MenuStringOption{backupOpt, restoreOpt}).Run()
	switch option {

	// Backup control mappings
	case backupOpt:
		backupFiles(mappingsDir, backupDir, true, ".xml")

	// Restore control mappings
	case restoreOpt:

		files, err := os.ReadDir(backupDir)
		if err != nil {
			fmt.Println("Unable to open backup directory: " + err.Error())
			disp.EnterToContinue()
			return
		}

		var restoreMenuOpts []disp.MenuStringOption
		for _, f := range files {
			if strings.HasSuffix(strings.ToLower(f.Name()), ctrlMapFileExt) {
				restoreMenuOpts = append(restoreMenuOpts, disp.MenuStringOption(f.Name()))
			}
		}

		if len(restoreMenuOpts) == 0 {
			fmt.Println("No objects found for restoration")
			disp.EnterToContinue()
			return
		}

		restoreMenuOpts = append(restoreMenuOpts, optBack)
		menuOption := disp.NewStringOptionMenu("Select file to restore", restoreMenuOpts).Run()
		if menuOption == optBack {
			return
		}

		restoreFiles(backupDir, mappingsDir, string(menuOption))
	}
	disp.EnterToContinue()
}

func BackupScreenshots() {
	choice := ptuOrLiveMenu.Run()
	if choice == optBack {
		return
	}
	gameDir, err := common.FindDir(RootDir, string(choice))
	if err != nil || gameDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}

	disp.ClearTerminal()
	fmt.Printf("\nGame directory found: %s\n", gameDir)

	screenshotDir := filepath.Join(gameDir, screenshotsDir)
	backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), screenshotsBackupDir, string(choice))

	backupFiles(screenshotDir, backupDir, false, ".jpg")

	disp.EnterToContinue()
}

func Exit() {
	disp.EnterToContinue()
	os.Exit(0)
}

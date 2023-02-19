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
	twoSecondDur     time.Duration = 2 * time.Second
)

var (
	RootDir string = ""
)

func ClearAllDataExceptP4k() {
	choice := ptuOrLiveMenu.Run()
	if choice == verBack {
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
	if choice == verBack {
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
	if choice == verBack {
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
	if choice == verBack {
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

func BackupControlMappings() {
	choice := ptuOrLiveMenu.Run()
	if choice == verBack {
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

	mappingsDir := filepath.Join(gameDir, "USER", "Client", "0", "Controls", "Mappings")
	files, err := os.ReadDir(mappingsDir)
	if err != nil {
		fmt.Println("Unable to open Mappings directory")
		disp.EnterToContinue()
	}

	controlMapCount := 0
	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name()), ".xml") {
			controlMapCount++
			backupFileName := strings.TrimSuffix(f.Name(), ".xml") + "-" + time.Now().Format("2006.01.02-15.04.05") + ".xml"

			backupDir := filepath.Join(filepath.Dir(filepath.Dir(gameDir)), "BACKUPS", string(choice))
			if _, err := os.Stat(backupDir); os.IsNotExist(err) {
				if err := os.MkdirAll(backupDir, 0755); err != nil {
					fmt.Println("Unable to create backup directory")
					disp.EnterToContinue()
					return
				}
				println("Backup directory created: " + backupDir)
			}

			if err := common.CopyFile(
				filepath.Join(mappingsDir, f.Name()),     // src
				filepath.Join(backupDir, backupFileName), // dst
			); err != nil {
				fmt.Printf("Backup error: %s\n", err.Error())
			} else {
				fmt.Printf("Mapping backed up: %s\n", backupFileName)
			}
		}
	}
	if controlMapCount == 0 {
		fmt.Println("No control mappings found")
	}

	disp.EnterToContinue()
}

func Exit() {
	disp.EnterToContinue()
	os.Exit(0)
}

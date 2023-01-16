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
	"github.com/inancgumus/screen"
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

	gameVersion := PtuOrLive()
	gameDir, err := common.FindDir(RootDir, string(gameVersion))
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

	if disp.YesOrNo() {
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
	gameVersion := PtuOrLive()
	userDir, err := common.FindDir(RootDir, string(gameVersion))
	if err != nil || userDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}
	exclusion := filepath.Join(userDir, "USER", "Client", "0", "Controls")
	fmt.Printf("\nUser directory found: %s\n", userDir)

	if disp.YesOrNo() {
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
	gameVersion := PtuOrLive()
	gameDir, err := common.FindDir(RootDir, string(gameVersion))
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
	gameVersion := PtuOrLive()
	gameDir, err := common.FindDir(RootDir, string(gameVersion))
	if err != nil || gameDir == "" {
		fmt.Println("Unable to find game directory")
		disp.EnterToContinue()
		return
	}

	screen.Clear()
	screen.MoveTopLeft()
	fmt.Print("Please enter your search phrase\n-> ")

	reader := bufio.NewReader(os.Stdin)
	phrase, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Invalid input: %s\n", err.Error())
		disp.EnterToContinue()
	}
	phrase = common.CleanInput(phrase)

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
	scAppDataDir := filepath.Join(common.UserHomeDir(), "AppData", "Local", "Star Citizen")
	files, _ := common.ListAllFilesAndDirs(scAppDataDir)

	if len(files) == 0 {
		fmt.Println("Star Citizen AppData is empty")
	} else {
		if disp.YesOrNo() {
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
	}
	disp.EnterToContinue()
}

func ClearRsiLauncherAppData() {

	for _, folder := range []string{"rsilauncher", "RSI Launcher"} {
		rsiLauncherDir := filepath.Join(common.UserHomeDir(), "AppData", "Roaming", folder)
		files, _ := common.ListAllFilesAndDirs(rsiLauncherDir)

		if len(files) == 0 {
			fmt.Printf("RSI Launcher AppData folder (%s) is empty!\n", folder)
			time.Sleep(twoSecondDur)
		} else {
			if disp.YesOrNo() {
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
	disp.EnterToContinue()
}

func Exit() {
	disp.EnterToContinue()
	os.Exit(0)
}

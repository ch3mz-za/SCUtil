package scu

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/pkg/common"
	p4k "github.com/ch3mz-za/SCUtil/pkg/p4kReader"
	"github.com/inancgumus/screen"
)

func clearAllDataExceptP4k() {

	gameVersion := PtuOrLive()
	gameDir, err := common.FindDir(RootDir, string(gameVersion))
	if err != nil || gameDir == "" {
		fmt.Println("Unable to find game directory")
		EnterToContinue()
		return
	}

	fmt.Printf("\nGame directory found: %s\n", gameDir)

	files, err := common.ListAllFilesAndDirs(gameDir)
	if err != nil {
		fmt.Println("Unable to list directories and files")
		EnterToContinue()
		return
	}

	if YesOrNo() {
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
		EnterToContinue()
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func clearUserFolder() {
	gameVersion := PtuOrLive()
	userDir, err := common.FindDir(RootDir, string(gameVersion))
	if err != nil || userDir == "" {
		fmt.Println("Unable to find game directory")
		EnterToContinue()
		return
	}
	exclusion := filepath.Join(userDir, "USER", "Client", "0", "Controls")
	fmt.Printf("\nUser directory found: %s\n", userDir)

	if YesOrNo() {
		err := filepath.Walk(filepath.Join(userDir, "USER"),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Check if file
				if !info.IsDir() {

					// Check prefix for exclusions
					if strings.HasPrefix(path, exclusion) {
						fmt.Println("Excluding: " + path)
					} else {

						if err := os.Remove(path); err != nil {
							fmt.Printf("Unable to remove file: %s | error: %s\n", path, err.Error())
						}
						fmt.Println("Removing: " + path)
					}
				}
				return nil
			})
		if err != nil {
			fmt.Printf("Error removing USER directory: %s\n", err.Error())
		}

		fmt.Println("Cleared USER directory")

	}
	EnterToContinue()
}

func getP4kFilenames() {
	gameVersion := PtuOrLive()
	userDir, err := common.FindDir(RootDir, string(gameVersion))
	if err != nil || userDir == "" {
		fmt.Println("Unable to find game directory")
		EnterToContinue()
		return
	}

	fmt.Printf("\nUser directory found: %s\n", userDir)
	p4k.GetP4kFilenames(userDir, p4kFileNameDir)
}

func mainFeatSearchP4kFilenames() {
	screen.Clear()
	screen.MoveTopLeft()
	fmt.Print("Please enter your search phrase\n-> ")

	reader := bufio.NewReader(os.Stdin)
	phrase, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Invalid input: %s\n", err.Error())
		EnterToContinue()
	}
	phrase = common.CleanInput(phrase)

	results, err := p4k.FindInFiles(p4kFileNameDir, phrase)
	if err != nil {
		fmt.Printf("Unable to search files: %s\n", err.Error())
		EnterToContinue()
	}

	println()
	for _, res := range results {
		println(res)
	}

	filename := strings.ReplaceAll(phrase, "\\", "_") + ".txt"
	p4k.MakeDir(p4kSearchResults)
	p4k.WriteStringsToFile(filepath.Join(p4kSearchResults, filename), results)
	EnterToContinue()
}

func clearStarCitizenAppData() {
	scAppDataDir := filepath.Join(common.UserHomeDir(), "AppData", "Local", "Star Citizen")
	files, _ := common.ListAllFilesAndDirs(scAppDataDir)

	if len(files) == 0 {
		fmt.Println("Star Citizen AppData is empty")
	} else {
		if YesOrNo() {
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
	EnterToContinue()
}

func clearRsiLauncherAppData() {

	for _, folder := range []string{"rsilauncher", "RSI Launcher"} {
		rsiLauncherDir := filepath.Join(common.UserHomeDir(), "AppData", "Roaming", folder)
		files, _ := common.ListAllFilesAndDirs(rsiLauncherDir)

		if len(files) == 0 {
			fmt.Printf("RSI Launcher AppData folder (%s) is empty!\n", folder)
			time.Sleep(twoSecondDur)
		} else {
			if YesOrNo() {
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
	EnterToContinue()
}

package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/menus"
	p4k "github.com/ch3mz-za/SCUtil/p4kReader"
	"github.com/inancgumus/screen"
	log "github.com/sirupsen/logrus"
)

const (
	p4kFileNameDir   string        = "./p4k_filenames"
	p4kSearchResults string        = "./p4k_search_results"
	twoSecondDur     time.Duration = 2 * time.Second
)

func main() {
	var rootDir string
	var err error

	rootDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Unable to determine working directory")
	}
	rootDir = filepath.Dir(rootDir)

	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); !os.IsNotExist(err) {
			rootDir = os.Args[1]
		}
	}

	for {
		switch menus.Main() {
		case menus.MainFeatClearAllExceptP4k:
			gameVersion := menus.PtuOrLive()
			gameDir, err := findDir(rootDir, string(gameVersion))
			if err != nil || gameDir == "" {
				fmt.Println("Unable to find game directory")
				menus.EnterToContinue()
				break
			}

			fmt.Printf("\nGame directory found: %s\n", gameDir)

			files, err := listAllFilesAndDirs(gameDir)
			if err != nil {
				fmt.Println("Unable to list directories and files")
				menus.EnterToContinue()
				break
			}

			if menus.YesOrNo() {
				clearAllExceptP4k(gameDir, files)
			}

		case menus.MainFeatClearUserFolder:
			gameVersion := menus.PtuOrLive()
			userDir, err := findDir(rootDir, string(gameVersion))
			if err != nil || userDir == "" {
				fmt.Println("Unable to find game directory")
				menus.EnterToContinue()
				break
			}
			fmt.Printf("\nUser directory found: %s\n", userDir)

			if menus.YesOrNo() {
				fmt.Println("Removing USER directory")
				os.Remove(userDir)
			}

		case menus.MainFeatGetP4kFilenames:
			gameVersion := menus.PtuOrLive()
			userDir, err := findDir(rootDir, string(gameVersion))
			if err != nil || userDir == "" {
				fmt.Println("Unable to find game directory")
				menus.EnterToContinue()
				break
			}

			fmt.Printf("\nUser directory found: %s\n", userDir)
			p4k.GetP4kFilenames(userDir, p4kFileNameDir)

		case menus.MainFeatSearchP4kFilenames:
			screen.Clear()
			screen.MoveTopLeft()
			fmt.Print("Please enter your search phrase\n-> ")

			reader := bufio.NewReader(os.Stdin)
			phrase, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Invalid input: %s\n", err.Error())
				menus.EnterToContinue()
			}
			phrase = strings.Replace(phrase, "\r\n", "", -1)

			results, err := p4k.FindInFiles(p4kFileNameDir, phrase)
			if err != nil {
				fmt.Printf("Unable to search files: %s\n", err.Error())
				menus.EnterToContinue()
			}

			println()
			for _, res := range results {
				println(res)
			}

			p4k.MakeDir(p4kSearchResults)
			p4k.WriteStringsToFile(filepath.Join(p4kSearchResults, phraseFileName(phrase)), results)
			menus.EnterToContinue()

		case menus.MainFeatExit:
			menus.EnterToContinue()
			return

		default:
			fmt.Println("Invalid menu option")
			time.Sleep(twoSecondDur)
		}
	}
}

func phraseFileName(phrase string) string {
	return strings.ReplaceAll(phrase, "\\", "_") + ".txt"
}

func listAllFilesAndDirs(dir string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func clearAllExceptP4k(gameDir string, files []fs.FileInfo) {
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".p4k") {
			continue
		}
		println("Deleted: " + filepath.Join(gameDir, f.Name()))
		os.Remove(filepath.Join(gameDir, f.Name()))
	}
}

func findDir(root, target string) (string, error) {

	var gamePath string
	err := filepath.WalkDir(root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if dir.IsDir() && filepath.Base(path) == target {
			gamePath = path
			return nil
		}
		return nil
	})

	return gamePath, err
}

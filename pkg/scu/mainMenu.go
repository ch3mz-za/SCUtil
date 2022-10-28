package scu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ch3mz-za/SCUtil/pkg/common"
	"github.com/inancgumus/screen"
	log "github.com/sirupsen/logrus"
)

const (
	p4kFileNameDir   string        = "./p4k_filenames"
	p4kSearchResults string        = "./p4k_search_results"
	twoSecondDur     time.Duration = 2 * time.Second
	utilVersion      string        = "v1.2.0"
)

var (
	RootDir string = ""
)

type Feature struct {
	FeatureName string
	Execute     func()
}

type SCUtil struct {
	MenuItems []*Feature
}

func NewMenu() *SCUtil {
	return &SCUtil{
		MenuItems: []*Feature{
			{
				FeatureName: "Clear all data except p4k",
				Execute:     clearAllDataExceptP4k,
			},
			{
				FeatureName: "Clear user folder",
				Execute:     clearUserFolder,
			},
			{
				FeatureName: "Get all filenames in p4k",
				Execute:     getP4kFilenames,
			},
			{
				FeatureName: "Search filenames in p4k",
				Execute:     mainFeatSearchP4kFilenames,
			},
			{
				FeatureName: "Clear Star Citizen App Data",
				Execute:     clearStarCitizenAppData,
			},
			{
				FeatureName: "Clear RSI Launcher data",
				Execute:     clearRsiLauncherAppData,
			},
			{
				FeatureName: "Exit",
				Execute:     exit,
			},
		},
	}
}

func exit() {
	EnterToContinue()
	os.Exit(0)
}

func (m *SCUtil) Run() {
	// Display main menu
	var invalidOption bool
	reader := bufio.NewReader(os.Stdin)
	for {
		invalidOption = false
		screen.Clear()
		screen.MoveTopLeft()

		fmt.Printf("SCUtil              [%s]\n----------------------------\n", utilVersion)
		for i, item := range m.MenuItems {
			fmt.Printf("%d. %s\n", i+1, item.FeatureName)
		}
		fmt.Printf("-> ")

		menuOption, err := reader.ReadString('\n')
		if err != nil {
			log.Error("Unable to read input: " + err.Error())
		}
		menuOption = common.CleanInput(menuOption)

		if val, err := strconv.Atoi(menuOption); err == nil {
			if 1 <= val && val <= len(m.MenuItems)+1 {
				m.MenuItems[val-1].Execute()
			} else {
				invalidOption = true
			}
		} else {
			invalidOption = true
		}

		if invalidOption {
			fmt.Println("Invalid menu option. Please enter correct number")
			time.Sleep(2 * time.Second)
		}
	}
}

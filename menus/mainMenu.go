package menus

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/inancgumus/screen"
	log "github.com/sirupsen/logrus"
)

type MainFeature int

const (
	MainFeatClearAllExceptP4k MainFeature = iota
	MainFeatClearUserFolder
	MainFeatGetP4kFilenames
	MainFeatSearchP4kFilenames
	MainFeatExit
)

func Main() MainFeature {
	// Display main menu
	reader := bufio.NewReader(os.Stdin)
	for {
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Print(strings.Join([]string{
			"Main menu",
			"----------------------------",
			"1. Clear all data except p4k",
			"2. Clear user folder",
			"3. Get all filenames in p4k",
			"4. Search filenames in p4k",
			"5. Exit",
			"-> ",
		}, "\n"))

		menuOption, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err.Error())
		}

		switch strings.ToLower(strings.Replace(menuOption, "\r\n", "", -1)) {
		case "1":
			return MainFeatClearAllExceptP4k
		case "2":
			return MainFeatClearUserFolder
		case "3":
			return MainFeatGetP4kFilenames
		case "4":
			return MainFeatSearchP4kFilenames
		case "5", "q", "quit", "exit":
			return MainFeatExit
		default:
			fmt.Println("Invalid menu option. Please enter correct number")
			time.Sleep(2 * time.Second)
		}
	}
}

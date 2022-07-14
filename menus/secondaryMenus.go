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

type SubFeature string

const (
	subFeatLive SubFeature = "LIVE"
	subFeatPtu  SubFeature = "PTU"
)

func PtuOrLive() SubFeature {
	// Display sub menu
	reader := bufio.NewReader(os.Stdin)
	for {
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Print(strings.Join([]string{
			"Game Version",
			"------------",
			"1. LIVE",
			"2. PTU",
			"-> ",
		}, "\n"))

		menuOption, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err.Error())
		}

		switch strings.ToLower(strings.Replace(menuOption, "\r\n", "", -1)) {
		case "1", "live":
			return subFeatLive
		case "2", "ptu":
			return subFeatPtu
		default:
			fmt.Println("Invalid menu option. Please enter correct number/word")
			time.Sleep(2 * time.Second)
		}
	}
}

func YesOrNo() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Print("Do you want to continue (y/n)\n-> ")

		ans, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err.Error())
		}

		switch strings.ToLower(strings.Replace(ans, "\r\n", "", -1)) {
		case "y", "live":
			return true
		case "n", "ptu":
			return false
		default:
			fmt.Println("Invalid menu option. Please enter correct letter")
			time.Sleep(2 * time.Second)
		}
	}
}

func EnterToContinue() {
	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

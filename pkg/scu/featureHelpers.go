package scu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/pkg/common"
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

		switch strings.ToLower(common.CleanInput(menuOption)) {
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

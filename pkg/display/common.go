package menus

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

		switch strings.ToLower(common.CleanInput(ans)) {
		case "y":
			return true
		case "n":
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

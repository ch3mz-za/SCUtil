package display

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/pkg/common"
	"github.com/inancgumus/screen"
	log "github.com/sirupsen/logrus"
)

type StringOptionMenu struct {
	MenuOptions  []MenuStringOption
	MaxMenuWidth int
	menuString   string
}

type MenuStringOption string

func NewStringOptionMenu(title string, menuOptions []MenuStringOption) *StringOptionMenu {
	menu := &StringOptionMenu{MenuOptions: menuOptions}

	menu.MaxMenuWidth = len(title)
	for _, opt := range menu.MenuOptions {
		if itemLen := len(string(opt)); itemLen+menuNumberPad > menu.MaxMenuWidth {
			menu.MaxMenuWidth = itemLen + menuNumberPad
		}
	}
	menu.compileMenu(title)
	return menu
}

func (m *StringOptionMenu) compileMenu(title string) {

	m.menuString += fmt.Sprintf(
		"%s\n%s\n",
		title,
		strings.Repeat("-", m.MaxMenuWidth),
	)

	for i, opt := range m.MenuOptions {
		m.menuString += fmt.Sprintf("%d. %s\n", i+1, string(opt))
	}
	m.menuString += "-> "
}



func (m *StringOptionMenu) Run() MenuStringOption {
	// Display main menu
	var invalidOption bool
	reader := bufio.NewReader(os.Stdin)
	for {
		invalidOption = false
		screen.Clear()
		screen.MoveTopLeft()

		fmt.Print(m.menuString)

		menuOption, err := reader.ReadString('\n')
		if err != nil {
			log.Error("Unable to read input: " + err.Error())
		}
		menuOption = common.CleanInput(menuOption)

		if val, err := strconv.Atoi(menuOption); err == nil {
			if 1 <= val && val <= len(m.MenuOptions)+1 {
				return m.MenuOptions[val-1]
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

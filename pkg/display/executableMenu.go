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

const (
	menuNumberPad int = 3
)

type MenuItem struct {
	Title   string
	Execute func()
}

type Menu struct {
	MenuItems    []*MenuItem
	MaxMenuWidth int
	menuString   string
}

func NewMenu(title, version string, menuItems []*MenuItem) *Menu {
	menu := &Menu{MenuItems: menuItems}

	menu.MaxMenuWidth = len(title)
	for _, item := range menu.MenuItems {
		if itemLen := len(item.Title); itemLen+menuNumberPad > menu.MaxMenuWidth {
			menu.MaxMenuWidth = itemLen + menuNumberPad
		}
	}
	menu.compileMenu(title, version)
	return menu
}

func (m *Menu) compileMenu(title, version string) {

	versionLen := len(version) + 2
	m.menuString += fmt.Sprintf(
		"%s%s[%s]\n%s\n",
		title,
		strings.Repeat(" ", m.MaxMenuWidth-len(title)-versionLen),
		version,
		strings.Repeat("-", m.MaxMenuWidth),
	)

	for i, item := range m.MenuItems {
		m.menuString += fmt.Sprintf("%d. %s\n", i+1, item.Title)
	}
	m.menuString += "-> "
}

func (m *Menu) Run() {
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

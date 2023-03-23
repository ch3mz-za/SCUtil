package tabs

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func ClearGameData() *fyne.Container {

	const (
		clearAlldataExceptP4k string = "Clear all data except p4k"
		clearUserData         string = "Clear USER data"
	)

	gameVersionDropDown.Selected = string(scu.Live)

	var removeControlMappings bool
	checkRemoveControlMappings := widget.NewCheck("Remove Control Mappings", func(value bool) {
		removeControlMappings = value
	})
	checkRemoveControlMappings.Hidden = true

	var clearDataSelection string

	return container.New(
		layout.NewVBoxLayout(),
		gameVersionDropDown,
		widget.NewRadioGroup([]string{clearAlldataExceptP4k, clearUserData}, func(value string) {
			clearDataSelection = value
			checkRemoveControlMappings.Hidden = value != clearUserData
		}),
		checkRemoveControlMappings,
		widget.NewButton("clear", func() {
			switch clearDataSelection {
			case clearAlldataExceptP4k:
				if err := scu.ClearAllDataExceptP4k(GameVersion); err != nil {
					fmt.Println(err.Error())
				}
			case clearUserData:
				if err := scu.ClearUserFolder(GameVersion, removeControlMappings); err != nil {
					fmt.Println(err.Error())
				}
			}
			log.Printf("tapped '%s' for game version: %s | remove control mappings: %v\n", clearDataSelection, GameVersion, removeControlMappings)
		}),
	)
}

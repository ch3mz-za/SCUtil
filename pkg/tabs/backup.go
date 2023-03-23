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

func Backup() fyne.CanvasObject {

	const (
		backupControlMappings string = "Backup Control Mappings"
		backupScreenshots     string = "Backup Screenshots"
	)

	var clearDataSelection string
	radioBackup := widget.NewRadioGroup([]string{backupControlMappings, backupScreenshots}, func(value string) {
		clearDataSelection = value
	})
	radioBackup.Selected = backupControlMappings

	return container.New(
		layout.NewVBoxLayout(),
		gameVersionDropDown,
		radioBackup,
		widget.NewButton("backup", func() {
			log.Println("backing up control mappings")
			switch clearDataSelection {
			case backupControlMappings:
				if err := scu.BackupControlMappings(GameVersion); err != nil {
					fmt.Println(err.Error())
				}
			case backupScreenshots:
				if err := scu.BackupScreenshots(GameVersion); err != nil {
					fmt.Println(err.Error())
				}
			}
		}),
	)
}

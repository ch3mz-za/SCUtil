package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func Backup(win fyne.Window) fyne.CanvasObject {

	const (
		backupControlMappings string = "Backup Control Mappings"
		backupScreenshots     string = "Backup Screenshots"
	)

	selectionGameVersion := widget.NewSelect([]string{scu.GameVerLIVE, scu.GameVerPTU}, func(value string) {})
	selectionGameVersion.Selected = scu.GameVerLIVE

	radioBackup := widget.NewRadioGroup([]string{backupControlMappings, backupScreenshots}, func(value string) {})
	radioBackup.Selected = backupControlMappings

	btn := widget.NewButton("backup", func() {
		var err error
		switch radioBackup.Selected {
		case backupControlMappings:
			if err = scu.BackupControlMappings(selectionGameVersion.Selected); err == nil {
				doneDiaglog(win)
			}
		case backupScreenshots:
			if err = scu.BackupScreenshots(selectionGameVersion.Selected); err == nil {
				doneDiaglog(win)
			}
		}
		if err != nil {
			dialog.ShowError(err, win)
		}
	})

	// TODO:
	// - Change to Border Container
	// - Add list of backed up files

	return container.New(
		layout.NewVBoxLayout(),
		selectionGameVersion,
		radioBackup,
		layout.NewSpacer(),
		btn,
	)
}

package frontend

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/internal/scu"
)

func backup(win fyne.Window) fyne.CanvasObject {

	backupFeatures := []struct {
		label, dir string
		fn         func(string) error
		openOption int
	}{
		{label: "Backup Control Mappings", dir: scu.ControlMappingsBackupDir, fn: scu.BackupControlMappings, openOption: openExternally},
		{label: "Backup Screenshots", dir: scu.ScreenshotsBackupDir, fn: scu.BackupScreenshots, openOption: openImage},
		{label: "Backup Characters", dir: scu.CharactersBackupDir, fn: scu.BackupUserCharacters, openOption: openExternally},
		{label: "Backup USER directory", dir: scu.UserBackupDir, fn: scu.BackupUserDirectory, openOption: openExternally},
	}

	selectionGameVersion := widget.NewSelect(scu.GetGameVersions(), func(value string) {})
	selectionGameVersion.Selected = scu.GameVerLIVE

	selectedBackupItem := 0
	listBackupItems := widget.NewList(
		func() int {
			return len(backupFeatures)
		},

		func() fyne.CanvasObject {
			btnFolder := widget.NewButton("", nil)
			btnFolder.SetIcon(theme.FolderIcon())
			return container.NewBorder(nil, nil, nil, btnFolder, widget.NewLabel("backup feature title"))
		},

		func(i widget.ListItemID, o fyne.CanvasObject) {
			lbl := o.(*fyne.Container).Objects[0].(*widget.Label)
			lbl.SetText(backupFeatures[i].label)

			btn := o.(*fyne.Container).Objects[1].(*widget.Button)
			btn.OnTapped = showOpenFileDialog(
				filepath.Join(scu.AppDir, backupFeatures[i].dir, selectionGameVersion.Selected),
				win,
				backupFeatures[i].openOption)
		})

	listBackupItems.Select(selectedBackupItem)
	listBackupItems.OnSelected = func(id int) { selectedBackupItem = id }

	btnBackup := widget.NewButton("Backup", func() {
		if err := backupFeatures[selectedBackupItem].fn(selectionGameVersion.Selected); err != nil {
			dialog.ShowError(err, win)
			return
		}
		doneDiaglog(win)
	})

	return widget.NewCard("", "", container.NewBorder(
		selectionGameVersion,
		btnBackup,
		nil, nil,
		listBackupItems,
	))
}

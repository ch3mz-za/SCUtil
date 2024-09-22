package frontend

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func clearGameData(win fyne.Window) fyne.CanvasObject {

	const (
		clearAlldataExceptP4k   string = "Clear all data except p4k"
		clearUserData           string = "Clear USER data"
		clearStarCitizenAppData string = "Clear Star Citizen AppData"
		clearRsiLauncherAppData string = "Clear RSI Launcher AppData"
	)

	clearFeatures := []string{
		clearStarCitizenAppData,
		clearRsiLauncherAppData,
		clearUserData,
		clearAlldataExceptP4k,
	}

	dropDownGameVersion := widget.NewSelect(scu.GetGameVersions(), func(value string) {})
	dropDownGameVersion.Selected = scu.GameVerLIVE
	dropDownGameVersion.Hidden = true

	checkRemoveControlMappings := widget.NewCheck("Remove Control Mappings", func(value bool) {})
	checkRemoveControlMappings.Hidden = true

	selectedBackupItem := 0
	listClearItems := widget.NewList(
		func() int {
			return len(clearFeatures)
		},

		func() fyne.CanvasObject {
			return widget.NewLabel("backup feature title")
		},

		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(clearFeatures[i])
		})

	listClearItems.Select(selectedBackupItem)
	listClearItems.OnSelected = func(id int) {
		selectedBackupItem = id
		checkRemoveControlMappings.Hidden = clearFeatures[id] != clearUserData
		dropDownGameVersion.Hidden = clearFeatures[id] != clearAlldataExceptP4k && clearFeatures[id] != clearUserData

	}

	btnClear := widget.NewButton("Clear", func() {
		var err error
		switch clearFeatures[selectedBackupItem] {
		case clearStarCitizenAppData:
			if err = scu.ClearStarCitizenAppData(); err == nil {
				dialog.ShowInformation("Files deleted", "Success!", win)
			}

		case clearRsiLauncherAppData:
			removedFiles := scu.ClearRsiLauncherAppData()
			if len(*removedFiles) != 0 {
				dialog.ShowInformation("AppData deleted", strings.Join(*removedFiles, "\n"), win)
			}

		case clearAlldataExceptP4k:
			if err = scu.ClearAllDataExceptP4k(dropDownGameVersion.Selected); err == nil {
				dialog.ShowInformation("Clear All Data", "data cleared", win)
			}

		case clearUserData:
			if err = scu.ClearUserFolder(dropDownGameVersion.Selected, checkRemoveControlMappings.Checked); err == nil {
				dialog.ShowInformation("Clear USER Data", "data cleared", win)
			}

		default:
			err = errors.New("invalid option")
		}

		if err != nil {
			dialog.ShowError(err, win)
		}
	})

	bottom := container.NewVBox(btnClear)
	return widget.NewCard("", "", container.NewBorder(
		nil, bottom, nil, nil,
		container.NewGridWithRows(2,
			listClearItems,
			container.NewVBox(dropDownGameVersion, checkRemoveControlMappings),
		),
	))
}

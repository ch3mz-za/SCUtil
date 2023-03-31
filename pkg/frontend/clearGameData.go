package frontend

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func ClearGameData(win fyne.Window) fyne.CanvasObject {

	const (
		clearAlldataExceptP4k   string = "Clear all data except p4k"
		clearUserData           string = "Clear USER data"
		clearStarCitizenAppData string = "Clear Star Citizen AppData"
		clearRsiLauncherAppData string = "Clear RSI Launcher AppData"
	)

	dropDownGameVersion := widget.NewSelect([]string{scu.GameVerLIVE, scu.GameVerPTU}, func(value string) {})
	dropDownGameVersion.Selected = scu.GameVerLIVE
	dropDownGameVersion.Hidden = true

	checkRemoveControlMappings := widget.NewCheck("Remove Control Mappings", func(value bool) {})
	checkRemoveControlMappings.Hidden = true

	radioGroup := widget.NewRadioGroup([]string{clearStarCitizenAppData, clearRsiLauncherAppData, clearUserData, clearAlldataExceptP4k}, func(value string) {
		checkRemoveControlMappings.Hidden = value != clearUserData
		dropDownGameVersion.Hidden = value != clearAlldataExceptP4k && value != clearUserData
	})
	radioGroup.Selected = clearStarCitizenAppData

	top := container.New(
		layout.NewVBoxLayout(),
		radioGroup,
		dropDownGameVersion,
		checkRemoveControlMappings,
	)

	bottom := widget.NewButton("clear", func() {
		var err error
		switch radioGroup.Selected {
		case clearStarCitizenAppData:
			var removedFiles *[]string
			removedFiles, err = scu.ClearStarCitizenAppData()
			if len(*removedFiles) != 0 {
				dialog.ShowInformation("Files deleted", strings.Join(*removedFiles, "\n"), win)
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

	return container.NewBorder(top, bottom, nil, nil, layout.NewSpacer())
}

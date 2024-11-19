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
		clearStarCitizenData    string = "Clear Star Citizen Data"
		clearStarCitizenAppData string = "Clear Star Citizen AppData"
		clearRsiLauncherAppData string = "Clear RSI Launcher AppData"
	)

	clearFeatures := []string{
		clearStarCitizenAppData,
		clearRsiLauncherAppData,
		clearStarCitizenData,
	}

	dropDownGameVersion := widget.NewSelect(scu.GetGameVersions(), func(value string) {})
	dropDownGameVersion.Selected = scu.GameVerLIVE
	dropDownGameVersion.Hidden = true

	checkRemoveControlMappings := widget.NewCheck("Remove control mappings", func(value bool) {})
	checkRemoveControlMappings.Hidden = true

	checkRemoveRenderSetting := widget.NewCheck("Exclude renderer setting", func(value bool) {})
	checkRemoveRenderSetting.Hidden = false

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

	const (
		clearOptUserData  = "User data"
		clearOptP4kData   = "P4k data"
		clearOptAllButP4k = "All except p4k data"
	)

	clearSUDataOptions := []string{
		clearOptUserData,
		clearOptP4kData,
		clearOptAllButP4k,
	}
	radioClearData := widget.NewRadioGroup(
		clearSUDataOptions,
		func(s string) {
			checkRemoveControlMappings.Hidden = s != clearOptUserData
		},
	)
	radioClearData.Hidden = true
	radioClearData.Selected = clearSUDataOptions[0]

	listClearItems.Select(selectedBackupItem)
	listClearItems.OnSelected = func(id int) {
		selectedBackupItem = id
		checkRemoveControlMappings.Hidden = clearFeatures[id] != clearStarCitizenData || radioClearData.Selected != clearOptUserData
		checkRemoveRenderSetting.Hidden = clearFeatures[id] != clearStarCitizenAppData
		dropDownGameVersion.Hidden = clearFeatures[id] != clearStarCitizenData
		radioClearData.Hidden = clearFeatures[id] != clearStarCitizenData
	}

	btnClear := widget.NewButton("Clear", func() {
		var err error
		switch clearFeatures[selectedBackupItem] {
		case clearStarCitizenAppData:
			if err = scu.ClearStarCitizenAppData(checkRemoveRenderSetting.Checked); err == nil {
				dialog.ShowInformation("Files deleted", "Success!", win)
			}

		case clearRsiLauncherAppData:
			removedFiles := scu.ClearRsiLauncherAppData()
			if len(*removedFiles) != 0 {
				dialog.ShowInformation("AppData deleted", strings.Join(*removedFiles, "\n"), win)
			}

		case clearStarCitizenData:

			switch radioClearData.Selected {
			case clearOptUserData:
				if err = scu.ClearUserFolder(dropDownGameVersion.Selected, checkRemoveControlMappings.Checked); err == nil {
					dialog.ShowInformation("Clear USER Data", "data cleared", win)
				}

			case clearOptP4kData:
				if err = scu.ClearP4kData(dropDownGameVersion.Selected); err == nil {
					dialog.ShowInformation("Clear P4k Data", "data cleared", win)
				}

			case clearOptAllButP4k:
				if err = scu.ClearAllDataExceptP4k(dropDownGameVersion.Selected); err == nil {
					dialog.ShowInformation("Clear All Data", "data cleared", win)
				}
			}

		default:
			err = errors.New("invalid option")
		}

		if err != nil {
			dialog.ShowError(err, win)
		}
	})

	topCard := widget.NewCard("", "", listClearItems)
	bottomCard := widget.NewCard("", "", container.NewBorder(
		nil, btnClear, nil, nil,
		container.NewVBox(dropDownGameVersion, radioClearData, checkRemoveControlMappings, checkRemoveRenderSetting),
	))
	return container.NewVSplit(topCard, bottomCard)
}

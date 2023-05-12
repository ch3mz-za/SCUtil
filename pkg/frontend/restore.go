package frontend

import (
	"fmt"
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func restore(win fyne.Window) fyne.CanvasObject {
	restoreData := binding.BindStringList(&[]string{})

	selectionGameVersion := widget.NewSelect([]string{scu.GameVerLIVE, scu.GameVerPTU}, func(value string) {
		items, err := scu.GetFilesListFromDir(filepath.Join(scu.AppDir, scu.ControlMappingsBackupDir, value))
		if err != nil {
			dialog.ShowError(err, win)
		}
		restoreData.Set(*items)
	})

	top := container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("Restore Control Mappings"),
		selectionGameVersion,
	)

	restoreList := widget.NewListWithData(restoreData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	var err error
	var itemToBeRestored string
	restoreList.OnSelected = func(id widget.ListItemID) {
		itemToBeRestored, err = restoreData.GetValue(id)
		if err != nil {
			dialog.ShowError(err, win)
		}
	}

	btnRestore := widget.NewButton("Restore", func() {
		log.Printf("Restore: %s\n", itemToBeRestored)
		if err = scu.RestoreControlMappings(selectionGameVersion.Selected, itemToBeRestored); err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("Restore Control Mappings", fmt.Sprintf("%s restored", itemToBeRestored), win)
		}
	})

	return widget.NewCard("", "", container.NewBorder(top, btnRestore, nil, nil, restoreList))
}

package frontend

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/config"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func settings(win fyne.Window, cfg *config.AppConfig) fyne.CanvasObject {

	gameDirData := binding.BindString(&scu.GameDir)
	gameDirLabel := widget.NewEntryWithData(gameDirData)

	btnSetGameDir := widget.NewButton("", func() {
		win.Resize(fyne.NewSize(700, 500))
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			defer resetToDefaultWindowSize(win)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}

			gameDirData.Set(list.Path())
			if !scu.IsGameDirectory(list.Path()) {
				dialog.ShowError(errors.New("not a valid game directory"), win)
				return
			}

			cfg.GameDir = list.Path()
			config.WriteAppConfig(config.AppConfigPath, cfg)
		}, win)
	})
	btnSetGameDir.SetIcon(theme.FolderIcon())

	cardGameDir := widget.NewCard("", "Game directory",
		container.NewBorder(nil, nil, nil, btnSetGameDir, gameDirLabel))

	return container.New(
		layout.NewVBoxLayout(),
		cardGameDir,
	)
}

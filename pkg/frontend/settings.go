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
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	btnSetGameDir := widget.NewButton("", func() {
		win.Resize(fyne.NewSize(700, 500))
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			defer resetToDefaultWindowSize(win)
			progressBar.Show()
			defer progressBar.Hide()

			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}

			resetToDefaultWindowSize(win)
			gameDir := scu.FindGameDirectory(list.Path())
			if gameDir == "" {
				dialog.ShowError(errors.New("could not find game directory"), win)
				return
			}

			gameDirData.Set(gameDir)
			if !scu.IsGameDirectory(gameDir) {
				dialog.ShowError(errors.New("not a valid game directory"), win)
				return
			}

			cfg.GameDir = gameDir
			config.WriteAppConfig(config.AppConfigPath, cfg)
		}, win)
	})
	btnSetGameDir.SetIcon(theme.FolderIcon())

	cardGameDir := widget.NewCard("", "Game directory", container.NewVBox(
		container.NewBorder(nil, nil, nil, btnSetGameDir, gameDirLabel),
		progressBar,
	))

	return container.New(
		layout.NewVBoxLayout(),
		cardGameDir,
	)
}

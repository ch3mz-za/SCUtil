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
	"github.com/ch3mz-za/SCUtil/internal/config"
	"github.com/ch3mz-za/SCUtil/internal/scu"
)

var (
	gameDirBind         = binding.NewString()
	gameDirBindListener binding.DataListener
)

func initGameDirBinding() {
	// Only add listener once
	if gameDirBindListener == nil {
		gameDirBindListener = binding.NewDataListener(func() {
			if dir, err := gameDirBind.Get(); err == nil {
				scu.SetGameDir(dir)
			}
		})
		gameDirBind.AddListener(gameDirBindListener)
	}
}

func settings(win fyne.Window, cfg *config.AppConfig) fyne.CanvasObject {

	gameDirLabel := widget.NewEntryWithData(gameDirBind)
	progressBar := widget.NewProgressBarInfinite()
	progressBar.Hide()

	btnSetGameDir := widget.NewButton("", func() {
		enlargeWindowForDialog(win)
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			defer resetToUserWindowSize(win)
			progressBar.Show()
			defer progressBar.Hide()

			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				return
			}

			resetToUserWindowSize(win)
			gameDir := scu.FindGameDirectory(list.Path())
			if gameDir == "" {
				dialog.ShowError(errors.New("could not find game directory"), win)
				return
			}

			if err := gameDirBind.Set(gameDir); err != nil {
				dialog.ShowError(err, win)
			}

			if !scu.IsGameDirectory(gameDir) {
				dialog.ShowError(errors.New("not a valid game directory"), win)
				return
			}

			cfg.GameDir = gameDir
			if err := config.WriteAppConfig(config.AppConfigPath, cfg); err != nil {
				dialog.ShowError(errors.New("unable to write app config"), win)
			}
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

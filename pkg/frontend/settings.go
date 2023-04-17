package frontend

import (
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

// TODO:
//	- select game dir using file explorer
//	- func: confirm game directory

func Settings(win fyne.Window, cfg *config.AppConfig) fyne.CanvasObject {

	// gameDirLabel := widget.NewLabel(scu.RootDir)
	gameDirData := binding.BindString(&scu.RootDir)
	gameDirLabel := widget.NewEntryWithData(gameDirData)
	gameDirLabel.SetText(scu.RootDir)

	btnSetGameDir := widget.NewButton("", func() {
		win.Resize(fyne.NewSize(700, 500))
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if list == nil {
				resetToDefaultWindowSize(win)
				return
			}
			gameDirData.Set(list.Path())
			cfg.GameDir = list.Path()

			config.WriteAppConfig(config.AppConfigPath, cfg)
			resetToDefaultWindowSize(win)

			// TODO: Confirm that you are happy with selected game directory

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

package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/config"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

// func MakeMainMenu(a fyne.App, w fyne.Window, cfg *config.AppConfig) *fyne.MainMenu {

// 	openGameDirSettings := func() {
// 		nw := a.NewWindow("Game Directory")
// 		nw.SetContent(settingsWindow(cfg))
// 		nw.Resize(fyne.NewSize(300, 100))
// 		nw.Show()
// 	}
// 	settingsItem := fyne.NewMenuItem("Game dir", openGameDirSettings)
// 	settings := fyne.NewMenu("Settings", settingsItem)

// 	return fyne.NewMainMenu(settings)
// }

type WindowList struct {
	Main     fyne.Window
	Settings fyne.Window
}

func Settings(wl WindowList, cfg *config.AppConfig) fyne.CanvasObject {
	// gameDirBinding := binding.BindString(&scu.RootDir)
	entry := widget.NewEntry()
	entry.Text = scu.RootDir
	btnSet := widget.NewButton("  set  ", func() {
		scu.RootDir = entry.Text
		cfg.GameDir = entry.Text
		// TODO: Find and set game_directory

		config.WriteAppConfig(config.AppConfigPath, cfg)
		if wl.Settings != nil {
			wl.Main.Show()
			wl.Settings.Close()
		}
	})

	return container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("Game directory"),
		container.NewBorder(nil, nil, nil, btnSet, entry),
	)
}

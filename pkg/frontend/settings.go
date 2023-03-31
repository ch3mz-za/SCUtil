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

func Settings(cfg *config.AppConfig) fyne.CanvasObject {
	// gameDirBinding := binding.BindString(&scu.RootDir)
	entry := widget.NewEntry()
	entry.Text = scu.RootDir
	btnSet := widget.NewButton("  set  ", func() {
		scu.RootDir = entry.Text
		cfg.GameDir = entry.Text
		config.WriteAppConfig(config.AppConfigPath, cfg)
	})
	return container.New(
		layout.NewVBoxLayout(),
		widget.NewLabel("Game directory"),
		container.NewBorder(nil, nil, nil, btnSet, entry),
	)
}

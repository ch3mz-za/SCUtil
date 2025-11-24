package frontend

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/ch3mz-za/SCUtil/internal/common"
	"github.com/ch3mz-za/SCUtil/internal/config"
	"github.com/ch3mz-za/SCUtil/internal/scu"
)

func SetupMainWindowContent(w fyne.Window) fyne.CanvasObject {

	var err error
	var cfg = &config.AppConfig{}
	var gameDir string

	// Initialize the game directory binding listener
	initGameDirBinding()

	scu.AppDir, _ = os.Getwd()
	if common.Exists(config.AppConfigPath) {
		cfg, err = config.ReadAppConfig(config.AppConfigPath)
		if err != nil {
			log.Fatalf("unable to load config: %s", err.Error())
		}

		gameDirBind.Set(cfg.GameDir)
		gameDir = cfg.GameDir
	}

	tabSettings := container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settings(w, cfg))
	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Backup", theme.StorageIcon(), backup(w)),
		container.NewTabItemWithIcon("Clean", theme.DeleteIcon(), clearGameData(w)),
		container.NewTabItemWithIcon("Restore", theme.HistoryIcon(), restore(w)),
		container.NewTabItemWithIcon("Data", theme.FileApplicationIcon(), p4kData(w)),
		container.NewTabItemWithIcon("Logs", theme.InfoIcon(), logs(w)),
		tabSettings,
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)

	mainTabs.OnSelected = func(ti *container.TabItem) {
		gameDir = getGameDir(w)

		if ti.Text == "Settings" && gameDir == "" {
			dialog.ShowInformation("Empty Game Directory", "Please set your game directory", w)
		}
	}

	if gameDir == "" {
		mainTabs.Select(tabSettings)
	}
	return mainTabs
}

func getGameDir(w fyne.Window) string {
	gameDir, err := gameDirBind.Get()
	if err != nil && w != nil {
		dialog.ShowError(err, w)
	}
	return gameDir
}

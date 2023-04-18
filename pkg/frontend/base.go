package frontend

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/ch3mz-za/SCUtil/pkg/common"
	"github.com/ch3mz-za/SCUtil/pkg/config"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func SetupMainWindowContent(w fyne.Window) fyne.CanvasObject {

	var err error
	var cfg = &config.AppConfig{}

	scu.AppDir, _ = os.Getwd()
	if common.Exists(config.AppConfigPath) {
		cfg, err = config.ReadAppConfig(config.AppConfigPath)
		if err != nil {
			log.Fatalf("unable to load config: %s", err.Error())
		}
		scu.GameDir = cfg.GameDir
	}

	tabSettings := container.NewTabItem("Settings", settings(w, cfg))
	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Clean", theme.DeleteIcon(), clearGameData(w)),
		container.NewTabItemWithIcon("Backup", theme.StorageIcon(), backup(w)),
		container.NewTabItemWithIcon("Restore", theme.UploadIcon(), restore(w)),
		container.NewTabItem("Advanced", advanced(w)),
		tabSettings,
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)

	mainTabs.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Settings" && scu.GameDir == "" {
			dialog.ShowInformation("Empty Game Directry", "Please set your game directory", w)
		}
	}

	if scu.GameDir == "" {
		mainTabs.Select(tabSettings)
	}
	return mainTabs
}

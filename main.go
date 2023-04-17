// go: generate
package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"github.com/ch3mz-za/SCUtil/pkg/common"
	"github.com/ch3mz-za/SCUtil/pkg/config"
	fend "github.com/ch3mz-za/SCUtil/pkg/frontend"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

const version string = "v2.0.2"

func main() {

	var err error
	var cfg = &config.AppConfig{GameDir: ""}

	// TODO: Write a handler
	if common.Exists(config.AppConfigPath) {
		cfg, err = config.ReadAppConfig(config.AppConfigPath)
		if err != nil {
			err = fmt.Errorf("unable to load config: %s", err)
			log.Println(err.Error())
		}
		scu.RootDir = cfg.GameDir
	}

	a := app.NewWithID("SCUtil")
	w := a.NewWindow(fmt.Sprintf("SCUtil - %s", version))
	w.SetMaster()

	tabSettings := container.NewTabItem("Settings", fend.Settings(w, cfg))
	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Clean", theme.DeleteIcon(), fend.ClearGameData(w)),
		container.NewTabItemWithIcon("Backup", theme.StorageIcon(), fend.Backup(w)),
		container.NewTabItemWithIcon("Restore", theme.UploadIcon(), fend.Restore(w)),
		container.NewTabItem("Advanced", fend.Advanced(w)),
		tabSettings,
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)
	mainTabs.OnSelected = func(ti *container.TabItem) {
		if ti.Text == "Settings" {
			dialog.NewInformation("Empty Game Directry", "Please set your game directory", w)
		}
	}
	if scu.RootDir == "" {
		mainTabs.SelectTab(tabSettings)
	}

	w.SetContent(mainTabs)
	w.Resize(fyne.NewSize(400, 310))
	w.Show()

	a.Run()
}

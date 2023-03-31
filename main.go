//go:generate
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/ch3mz-za/SCUtil/pkg/common"
	"github.com/ch3mz-za/SCUtil/pkg/config"
	fend "github.com/ch3mz-za/SCUtil/pkg/frontend"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func main() {

	var err error
	var cfg *config.AppConfig

	if common.Exists(config.AppConfigPath) {
		cfg, err = config.ReadAppConfig(config.AppConfigPath)
		if err != nil {
			err = fmt.Errorf("unable to load config: %s", err)
			log.Println(err.Error())
		}
		scu.RootDir = cfg.GameDir
	}

	if scu.RootDir == "" {
		scu.RootDir, err = os.Getwd()
		if err != nil {
			log.Fatal("Unable to determine working directory")
		}
		scu.RootDir = filepath.Dir(scu.RootDir)
	}

	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); !os.IsNotExist(err) {
			scu.RootDir = os.Args[1]
		}
	}

	a := app.NewWithID("SCUtil-v2.0.2")
	w := a.NewWindow("SCUtil - v2.0.2")
	w.SetMaster()

	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Clean", theme.DeleteIcon(), fend.ClearGameData(w)),
		container.NewTabItemWithIcon("Backup", theme.StorageIcon(), fend.Backup(w)),
		container.NewTabItemWithIcon("Restore", theme.UploadIcon(), fend.Restore(w)),
		container.NewTabItem("Advanced", fend.Advanced(w)),
		container.NewTabItem("Settings", fend.Settings(cfg)),
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(mainTabs)
	w.Resize(fyne.NewSize(400, 310))
	w.Show()
	a.Run()
}

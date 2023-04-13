// go: generate
package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	// else {
	// 	file, _ := os.Create(config.AppConfigPath)
	// 	file.close()
	// }

	a := app.NewWithID(fmt.Sprintf("SCUtil-%s", version))
	w := a.NewWindow(fmt.Sprintf("SCUtil - %s", version))
	wl := fend.WindowList{Main: w, Settings: nil}

	var tw fyne.Window = nil
	if scu.RootDir == "" {
		w.Hide()
		tw = a.NewWindow("Set game directory")
		wl.Settings = tw
		tw.SetContent(fend.Settings(wl, cfg))
		tw.Resize(fyne.NewSize(400, 100))
		tw.Show()
	} else {
		w.Show()
	}
	// w.Show()

	w.SetMaster()

	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Clean", theme.DeleteIcon(), fend.ClearGameData(w)),
		container.NewTabItemWithIcon("Backup", theme.StorageIcon(), fend.Backup(w)),
		container.NewTabItemWithIcon("Restore", theme.UploadIcon(), fend.Restore(w)),
		container.NewTabItem("Advanced", fend.Advanced(w)),
		container.NewTabItem("Settings", fend.Settings(wl, cfg)),
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(mainTabs)
	w.Resize(fyne.NewSize(400, 310))

	a.Run()
}

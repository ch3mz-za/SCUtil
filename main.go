//go:generate
package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
	"github.com/ch3mz-za/SCUtil/pkg/tabs"
)

const (
	appTitle   string = "SCUtil"
	appVersion string = "v1.3.0"
)

func main() {

	var err error
	scu.RootDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Unable to determine working directory")
	}
	scu.RootDir = filepath.Dir(scu.RootDir)

	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); !os.IsNotExist(err) {
			scu.RootDir = os.Args[1]
		}
	}

	a := app.NewWithID("SCUtil-v2.0.0")
	w := a.NewWindow("SCUtil - v2.0.0")
	mainTabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Clean", theme.DeleteIcon(), tabs.ClearGameData(w)),
		container.NewTabItemWithIcon("Backup", theme.StorageIcon(), tabs.Backup(w)),
		container.NewTabItemWithIcon("Restore", theme.UploadIcon(), tabs.Restore(w)),
		container.NewTabItem("Advanced", tabs.Advanced(w)),
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(mainTabs)
	w.Resize(fyne.NewSize(400, 310))
	w.Show()
	a.Run()
}

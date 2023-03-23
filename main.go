//go:generate
package main

import (
	"log"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	// var menuItems = []*disp.MenuItem{
	// 	{
	// 		Title:   "Clear all data except p4k",
	// 		Execute: scu.ClearAllDataExceptP4k,
	// 	},
	// 	{
	// 		Title:   "Clear user folder (excluding control mappings)",
	// 		Execute: scu.ClearUserFolerWithExclusions,
	// 	},
	// 	{
	// 		Title:   "Clear user folder (including control mappings)",
	// 		Execute: scu.ClearUserFolerWithoutExclusions,
	// 	},
	// 	{
	// 		Title:   "Get all filenames in p4k",
	// 		Execute: scu.GetP4kFilenames,
	// 	},
	// 	{
	// 		Title:   "Search filenames in p4k",
	// 		Execute: scu.SearchP4kFilenames,
	// 	},
	// 	{
	// 		Title:   "Clear Star Citizen App Data (Windows AppData)",
	// 		Execute: scu.ClearStarCitizenAppData,
	// 	},
	// 	{
	// 		Title:   "Clear RSI Launcher data (Windows AppData)",
	// 		Execute: scu.ClearRsiLauncherAppData,
	// 	},
	// 	{
	// 		Title:   "Backup & restore control mappings",
	// 		Execute: scu.BackupOrRestoreControlMappings,
	// 	},
	// 	{
	// 		Title:   "Backup screenshots",
	// 		Execute: scu.BackupScreenshots,
	// 	},
	// 	{
	// 		Title:   "Exit",
	// 		Execute: scu.Exit,
	// 	},
	// }

	a := app.NewWithID("SCUtil-v2.0.0")
	w := a.NewWindow("SCUtil - v2.0.0 (alpha)")

	mainTabs := container.NewAppTabs(
		// container.NewTabItem("Clean Game Data", tabs.ClearGameData()),
		// container.NewTabItem("Clean User Data", widget.NewLabel("Hello")),
		// container.NewTabItem("Clean AppData", tabs.ClearAppData()),
		container.NewTabItem("Backup", tabs.Backup()),
		container.NewTabItem("Restore", tabs.Restore(w)),
	)
	mainTabs.SetTabLocation(container.TabLocationLeading)
	w.SetContent(mainTabs)
	// tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	w.Resize(fyne.NewSize(400, 200))
	// w.SetFixedSize(true)
	w.ShowAndRun()
}

// go: generate
package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	fend "github.com/ch3mz-za/SCUtil/pkg/frontend"
)

const version string = "v2.3.0"

func main() {

	a := app.NewWithID("SCUtil")
	a.Settings().SetTheme(theme.LightTheme())

	// TODO: Chane theme color here

	w := a.NewWindow(fmt.Sprintf("SCUtil - %s", version))
	w.SetMaster()
	w.SetContent(fend.SetupMainWindowContent(w))
	w.Resize(fend.DefaultAppWinSize)
	w.Show()

	a.Run()
}

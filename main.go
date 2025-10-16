// go: generate
package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"

	fend "github.com/ch3mz-za/SCUtil/internal/frontend"
)

const version string = "v2.7.0"

func main() {

	a := app.NewWithID("SCUtil")

	w := a.NewWindow(fmt.Sprintf("SCUtil - %s", version))
	w.SetMaster()
	w.SetContent(fend.SetupMainWindowContent(w))
	w.Resize(fend.DefaultAppWinSize)
	w.Show()

	a.Run()
}

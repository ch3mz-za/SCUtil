// go: generate
package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	fend "github.com/ch3mz-za/SCUtil/pkg/frontend"
)

const version string = "v2.1.0"

func main() {

	a := app.NewWithID("SCUtil")

	w := a.NewWindow(fmt.Sprintf("SCUtil - %s", version))
	w.SetMaster()
	w.SetContent(fend.SetupMainWindowContent(w))
	w.Resize(fyne.NewSize(400, 310))
	w.Show()

	a.Run()
}

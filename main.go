// go: generate
package main

import (
	"fmt"
	// _ "net/http/pprof"

	"fyne.io/fyne/v2/app"

	fend "github.com/ch3mz-za/SCUtil/pkg/frontend"
)

const version string = "v2.3.1"

func main() {

	// defer profile.Start(profile.MemProfile).Stop()

	// go func() {
	// 	http.ListenAndServe(":8080", nil)
	// }()

	a := app.NewWithID("SCUtil")

	w := a.NewWindow(fmt.Sprintf("SCUtil - %s", version))
	w.SetMaster()
	w.SetContent(fend.SetupMainWindowContent(w))
	w.Resize(fend.DefaultAppWinSize)
	w.Show()

	a.Run()
}

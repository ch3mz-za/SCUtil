package frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func doneDiaglog(win fyne.Window) {
	dialog.ShowInformation("Status", "Completed successfully", win)
}

func resetToDefaultWindowSize(win fyne.Window) {
	win.Resize(fyne.Size{Width: 400, Height: 310})
}

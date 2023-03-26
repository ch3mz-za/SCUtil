package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func doneDiaglog(win fyne.Window) {
	dialog.ShowInformation("Status", "Completed successfully", win)
}

package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func ClearAppData() *fyne.Container {

	const (
		clearStarCitizenAppData string = "Clear Star Citizen AppData"
		clearRsiLauncherAppData string = "Clear RSI Launcher AppData"
	)

	var clearDataSelection string
	radioClear := widget.NewRadioGroup([]string{clearStarCitizenAppData, clearRsiLauncherAppData}, func(value string) {
		clearDataSelection = value
	})
	radioClear.Selected = clearStarCitizenAppData

	return container.New(
		layout.NewVBoxLayout(),
		radioClear,
		widget.NewButton("clear", func() {
			switch clearDataSelection {
			case clearStarCitizenAppData:
				scu.ClearStarCitizenAppData()
			case clearRsiLauncherAppData:
				scu.ClearRsiLauncherAppData()
			}
		}),
	)
}

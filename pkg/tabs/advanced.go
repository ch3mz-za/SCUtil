package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
)

func Advanced(win fyne.Window) fyne.CanvasObject {

	selectionGameVersion := widget.NewSelect([]string{scu.GameVerLIVE, scu.GameVerPTU}, func(value string) {})
	selectionGameVersion.Selected = scu.GameVerLIVE

	progress := widget.NewProgressBarInfinite()
	progress.Stop()
	progress.Hide()

	entry := widget.NewEntry()
	searchField := &widget.Form{
		Items: []*widget.FormItem{{Text: "Phrase", Widget: entry}},
		OnSubmit: func() {
			progress.Show()
			progress.Start()
			if err := scu.SearchP4kFilenames(selectionGameVersion.Selected, entry.Text); err != nil {
				dialog.ShowError(err, win)
			} else {
				doneDiaglog(win)
			}
			progress.Stop()
			progress.Hide()
		},
	}

	searchCard := widget.NewCard("", "Search Data.p4k filenames", container.NewVBox(
		searchField,
	))

	return container.New(
		layout.NewVBoxLayout(),
		selectionGameVersion,
		searchCard,
		container.New(
			layout.NewGridLayoutWithColumns(2),
			widget.NewLabel("Extract Data.p4k filenames"),
			widget.NewButton("extract", func() {
				progress.Show()
				progress.Start()
				if err := scu.GetP4kFilenames(selectionGameVersion.Selected); err != nil {
					dialog.ShowError(err, win)
				} else {
					doneDiaglog(win)
				}
				progress.Stop()
				progress.Hide()
			}),
		),
		progress,
	)
}

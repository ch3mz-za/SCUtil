package frontend

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/pkg/common"
	"github.com/ch3mz-za/SCUtil/pkg/scu"
	"github.com/skratchdot/open-golang/open"
)

func getSearchResults(resultsDir string, win fyne.Window) *[]string {
	items, err := scu.GetFilesListFromDir(resultsDir)
	if err != nil {
		dialog.ShowError(err, win)
	}
	return items
}

func p4kData(win fyne.Window) fyne.CanvasObject {

	var searchResultsDir string
	btnOpenP4kFilenames := widget.NewButtonWithIcon("", theme.FileTextIcon(), func() {
		open.Run(filepath.Join(scu.AppDir))
	})
	btnOpenP4kFilenames.Disabled()

	searchData := binding.BindStringList(&[]string{})
	selectionGameVersion := widget.NewSelect([]string{scu.GameVerLIVE, scu.GameVerPTU}, func(value string) {
		// Set P4k filename open buttin state
		btnOpenP4kFilenames.Disable()
		if common.Exists(filepath.Join(scu.AppDir, scu.P4kFilenameResultsDir, value)) {
			btnOpenP4kFilenames.Enable()
		}

		// Get list of search results
		searchResultsDir = filepath.Join(scu.AppDir, scu.P4kSearchResultsDir, value)
		if !common.Exists(searchResultsDir) {
			common.MakeDir(searchResultsDir)
		}

		// Set search results
		searchData.Set(*getSearchResults(searchResultsDir, win))
	})

	// search result list
	searchList := widget.NewListWithData(searchData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		})

	// TODO: Only do this if there are items in this list
	// setup selected item in list
	var selectedSearchResult int
	searchList.OnSelected = func(id widget.ListItemID) {
		selectedSearchResult = id
	}

	// progressbar
	progress := widget.NewProgressBarInfinite()
	progress.Stop()
	progress.Hide()

	// search button
	entrySearch := widget.NewEntry()
	btnSearch := widget.NewButton("Search P4k", func() {
		toggleProgress(progress)
		defer toggleProgress(progress)
		if err := scu.SearchP4kFilenames(selectionGameVersion.Selected, entrySearch.Text); err != nil {
			dialog.ShowError(err, win)
		} else {
			doneDiaglog(win)
		}

		// Get list of search results
		items, err := scu.GetFilesListFromDir(filepath.Join(scu.AppDir, scu.P4kSearchResultsDir, selectionGameVersion.Selected))
		if err != nil {
			dialog.ShowError(err, win)
		}
		searchData.Set(*items)
	})

	// delete search result button
	btnDelete := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {

	})

	// open search result button
	btnOpenSearchResult := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		itemToBeOpened, err := searchData.GetValue(selectedSearchResult)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		if err := open.Run(filepath.Join(searchResultsDir, itemToBeOpened)); err != nil {
			dialog.ShowError(err, win)
		}
	})

	searchResLabel := widget.NewLabel("Search Results")
	searchResLabel.TextStyle.Bold = true
	top := container.NewVBox(
		selectionGameVersion,
		container.NewBorder(nil, nil, nil, btnSearch, entrySearch),
		container.NewBorder(
			widget.NewSeparator(), widget.NewSeparator(),
			searchResLabel,
			container.NewHBox(btnOpenSearchResult, btnDelete),
			layout.NewSpacer(),
		),
	)

	btnGetP4kFilenames := widget.NewButton("Get P4k Filenames", func() {
		toggleProgress(progress)
		defer toggleProgress(progress)
		if err := scu.GetP4kFilenames(selectionGameVersion.Selected); err != nil {
			dialog.ShowError(err, win)
		} else {
			doneDiaglog(win)
		}
	})

	btnOpenP4kFilenames.Disable()

	bottom := container.NewVBox(
		container.NewBorder(nil, nil, nil, btnOpenP4kFilenames, btnGetP4kFilenames),
		progress,
	)

	// TODO: replace searchList with scrollable list - list can get too long for the view
	return widget.NewCard("", "", container.NewBorder(top, bottom, nil, nil, searchList))
}

func toggleProgress(p *widget.ProgressBarInfinite) {
	if !p.Visible() {
		p.Show()
		p.Start()
	} else {
		p.Stop()
		p.Hide()
	}
}

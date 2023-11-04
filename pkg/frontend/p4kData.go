package frontend

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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
	var selectedSearchResult int = -1
	searchData := binding.BindStringList(&[]string{})

	// open p4k filenames button
	var p4kFilenamesResult string
	btnOpenP4kFilenames := widget.NewButtonWithIcon("", theme.FileTextIcon(), func() {
		if err := open.Run(p4kFilenamesResult); err != nil {
			dialog.ShowError(err, win)
			return
		}
	})
	btnOpenP4kFilenames.Disable()

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
	btnOpenSearchResult.Disable()

	// delete search result button
	btnDelete := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	btnDelete.OnTapped = func() {
		itemToBeDeleted, err := searchData.GetValue(selectedSearchResult)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		if err := os.Remove(filepath.Join(searchResultsDir, itemToBeDeleted)); err != nil {
			dialog.ShowError(err, win)
			return
		}

		if err := searchData.Set(*getSearchResults(searchResultsDir, win)); err != nil {
			dialog.ShowError(err, win)
			return
		}

		if searchData.Length() == 0 {
			btnDelete.Disable()
			btnOpenSearchResult.Disable()
		}
	}
	btnDelete.Disable()

	// search result list
	searchList := widget.NewListWithData(searchData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))

		})

	// setup selected item in list
	searchList.OnSelected = func(id widget.ListItemID) {
		selectedSearchResult = id
		btnDelete.Enable()
		btnOpenSearchResult.Enable()
	}

	selectionGameVersion := widget.NewSelect(scu.GetGameVersions(), func(value string) {
		// set P4k filename open button state
		btnOpenP4kFilenames.Disable()
		p4kFilenamesResult = filepath.Join(scu.AppDir, fmt.Sprintf(scu.P4kFilenameResultsDir, value))
		if common.Exists(p4kFilenamesResult) {
			btnOpenP4kFilenames.Enable()
		}

		// set list of search results
		searchResultsDir = filepath.Join(scu.AppDir, scu.P4kSearchResultsDir, value)
		if !common.Exists(searchResultsDir) {
			common.MakeDir(searchResultsDir)
		}

		// set search results and button states
		if err := searchData.Set(*getSearchResults(searchResultsDir, win)); err != nil {
			dialog.ShowError(err, win)
		}
		btnDelete.Disable()
		btnOpenSearchResult.Disable()

		selectedSearchResult = -1
		searchList.UnselectAll()
	})

	// progressbar
	progress := widget.NewProgressBarInfinite()
	progress.Stop()
	progress.Hide()

	// search button
	entrySearch := widget.NewEntry()
	entrySearch.SetPlaceHolder("Enter phrase here")
	btnSearch := widget.NewButton("Search P4k", func() {
		if selectionGameVersion.Selected == "" {
			dialog.ShowError(errors.New("no game version selected"), win)
			return
		}

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

		if err := searchData.Set(*items); err != nil {
			dialog.ShowError(err, win)
		}

		entrySearch.SetText("")
		runtime.GC()
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
			btnOpenP4kFilenames.Enable()
		}
		runtime.GC()
	})

	bottom := container.NewVBox(
		container.NewBorder(nil, nil, nil, btnOpenP4kFilenames, btnGetP4kFilenames),
		progress,
	)

	return widget.NewCard("", "", container.NewBorder(top, bottom, nil, nil, container.NewVScroll(searchList)))
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

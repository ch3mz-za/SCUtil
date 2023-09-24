package frontend

import (
	"errors"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/skratchdot/open-golang/open"
)

const (
	openExternally = iota
	openImage
	openText
)

var (
	DefaultAppWinSize fyne.Size = fyne.NewSize(500, 450)
	UserPreferredSize fyne.Size = fyne.NewSize(500, 450)
)

func doneDiaglog(win fyne.Window) {
	dialog.ShowInformation("Status", "Completed successfully", win)
}

func resetToUserWindowSize(win fyne.Window) {
	win.Resize(UserPreferredSize)
}

func enlargeWindowForDialog(win fyne.Window) {
	UserPreferredSize = fyne.NewSize(
		win.Canvas().Size().Width,
		win.Canvas().Size().Height,
	)

	var w float32 = 700
	var h float32 = 500

	if win.Canvas().Size().Width > w {
		w = win.Canvas().Size().Width
	}

	if win.Canvas().Size().Height > h {
		h = win.Canvas().Size().Height
	}

	win.Resize(fyne.NewSize(w, h))
}

func showOpenFileDialog(dirPath string, win fyne.Window, openOpt int) func() {
	return func() {
		enlargeWindowForDialog(win)
		folderDiag := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			defer resetToUserWindowSize(win)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}

			switch openOpt {
			case openExternally:
				if err := open.Run(reader.URI().Path()); err != nil {
					dialog.ShowError(err, win)
					return
				}
			case openImage:
				showImage(reader)
			default:
				dialog.ShowError(errors.New("invalid mechanism for opening file"), win)
			}

		}, win)

		uri, err := storage.ListerForURI(storage.NewFileURI(dirPath))
		if err != nil {
			dialog.ShowError(errors.New("No directory found. Perform a backup first."), win)
			resetToUserWindowSize(win)
			return
		}

		folderDiag.SetLocation(uri)
		folderDiag.Show()
	}
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := io.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)
	return canvas.NewImageFromResource(res)
}

func showImage(f fyne.URIReadCloser) {
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillContain

	w := fyne.CurrentApp().NewWindow(f.URI().Name())
	w.SetContent(container.NewScroll(img))
	w.Resize(fyne.NewSize(1024, 720))
	w.Show()
}

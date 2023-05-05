package frontend

import (
	"errors"
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/skratchdot/open-golang/open"
)

func doneDiaglog(win fyne.Window) {
	dialog.ShowInformation("Status", "Completed successfully", win)
}

func resetToDefaultWindowSize(win fyne.Window) {
	win.Resize(fyne.Size{Width: 400, Height: 310})
}

const (
	openExternally = iota
	openImage
	openText
)

func showOpenFileDialog(dirPath string, win fyne.Window, openOpt int) func() {
	return func() {
		win.Resize(fyne.NewSize(700, 500))
		folderDiag := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			defer resetToDefaultWindowSize(win)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}

			switch openOpt {
			case openExternally:
				open.Run(reader.URI().Path())
			case openImage:
				showImage(reader)
			default:
				dialog.ShowError(errors.New("invalid mechanism for opening file"), win)
			}

		}, win)

		uri, err := storage.ListerForURI(storage.NewFileURI(dirPath))
		if err != nil {
			dialog.ShowError(err, win)
			resetToDefaultWindowSize(win)
			return
		}

		folderDiag.SetLocation(uri)
		folderDiag.Show()
	}
}

func fileSaved(f fyne.URIWriteCloser, w fyne.Window) {
	defer f.Close()
	_, err := f.Write([]byte("Written by Fyne demo\n"))
	if err != nil {
		dialog.ShowError(err, w)
	}
	err = f.Close()
	if err != nil {
		dialog.ShowError(err, w)
	}
	log.Println("Saved to...", f.URI())
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := ioutil.ReadAll(f)
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

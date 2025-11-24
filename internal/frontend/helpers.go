package frontend

import (
	"errors"
	"io"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ch3mz-za/SCUtil/internal/scu"
	"github.com/skratchdot/open-golang/open"
)

const (
	openExternally = iota
	openImage
	openText
)

var (
	DefaultAppWinSize fyne.Size = fyne.NewSize(1000, 550)
	UserPreferredSize fyne.Size = fyne.NewSize(1000, 550)
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

var ErrFileOpenCancelled = errors.New("file open cancelled")

func fileOpenPath(dirPath string, w fyne.Window, filter storage.FileFilter) (<-chan string, <-chan error) {
	pathCh := make(chan string, 1)
	errCh := make(chan error, 1)

	d := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
		defer func() {
			if rc != nil {
				_ = rc.Close()
			}
		}()

		if err != nil {
			errCh <- err
			return
		}
		if rc == nil {
			errCh <- ErrFileOpenCancelled
			return
		}

		pathCh <- filePathFromURI(rc.URI())
	}, w)

	uri, err := storage.ListerForURI(storage.NewFileURI(dirPath))
	if err == nil {
		d.SetLocation(uri)
	}

	if filter != nil {
		d.SetFilter(filter)
	}

	d.Show()
	return pathCh, errCh
}

func folderOpenPath(dirPath string, w fyne.Window) (<-chan string, <-chan error) {
	pathCh := make(chan string, 1)
	errCh := make(chan error, 1)

	d := dialog.NewFolderOpen(func(rc fyne.ListableURI, err error) {
		if err != nil {
			errCh <- err
			return
		}

		if rc == nil {
			errCh <- ErrFileOpenCancelled
			return
		}

		pathCh <- rc.Path()
	}, w)

	uri, err := storage.ListerForURI(storage.NewFileURI(dirPath))
	if err == nil {
		d.SetLocation(uri)
	}

	d.SetLocation(uri)
	d.Show()
	return pathCh, errCh
}

func filePathFromURI(u fyne.URI) string {
	// For local files prefer a native path; otherwise return the URI string
	if u.Scheme() != "file" {
		return u.String()
	}
	p := u.Path()
	if runtime.GOOS == "windows" && strings.HasPrefix(p, "/") {
		p = p[1:] // trim leading slash on Windows
	}
	return filepath.FromSlash(p)
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
		if err == nil {
			folderDiag.SetLocation(uri)
		}

		folderDiag.Show()
	}
}

func showOpenFolderDialog(dirPath string, win fyne.Window, openOpt int) func() {
	return func() {
		enlargeWindowForDialog(win)
		folderDiag := dialog.NewFolderOpen(func(path fyne.ListableURI, err error) {
			defer resetToUserWindowSize(win)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if path == nil {
				return
			}

			switch openOpt {
			case openExternally:
				if err := open.Run(path.URI().Path()); err != nil {
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
		if err == nil {
			folderDiag.SetLocation(uri)
		}

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

// newGameVersionSelect creates a Select widget that automatically updates
// when the game directory changes. The onChange callback is optional.
func newGameVersionSelect(onChange func(string)) *widget.Select {
	sel := widget.NewSelect(scu.GetGameVersions(), onChange)

	// Add a listener to refresh version options when game directory changes
	gameDirBind.AddListener(binding.NewDataListener(func() {
		versions := scu.GetGameVersions()
		sel.Options = versions

		// If the current selection is not in the new options, clear it
		if sel.Selected != "" {
			found := false
			for _, v := range versions {
				if v == sel.Selected {
					found = true
					break
				}
			}
			if !found {
				sel.Selected = ""
			}
		}

		// Auto-select LIVE if available
		if sel.Selected == "" && len(versions) > 0 {
			for _, v := range versions {
				if v == scu.GameVerLIVE {
					sel.SetSelected(scu.GameVerLIVE)
					break
				}
			}
			// If LIVE not found, select first option
			if sel.Selected == "" {
				sel.SetSelected(versions[0])
			}
		}

		sel.Refresh()
	}))

	return sel
}

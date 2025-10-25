package frontend

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/ch3mz-za/SCUtil/internal/logmon"
	"github.com/ch3mz-za/SCUtil/internal/scu"
)

type ParseOption int

const (
	TailFile ParseOption = iota
	SingleFile
	AggregateFiles
)

func logs(win fyne.Window) fyne.CanvasObject {
	const logItemLimit = 1000
	// State
	logItems := make([]*logmon.LogItem, 0, logItemLimit) // in-memory store
	activeFilter := None                                 // or Humans/Vehicles/Search
	activeQuery := ""                                    // search string

	// --- Multi-line log view using VBox inside a VScroll (variable-height rows) ---
	logBox := container.NewVBox()
	scroll := container.NewVScroll(logBox)

	// Rebuild the entire view (call this only when filter/search changes)
	rebuild := func() {
		fyne.Do(func() {
			logBox.RemoveAll()
			for _, it := range logItems {
				if filterMatch(it, activeFilter, activeQuery) {
					logBox.Add(parseLogResult(it))
				}
			}
			logBox.Refresh()
			scroll.ScrollToBottom()
		})
	}

	// Bottom controls
	selectionGameVersion := widget.NewSelect(scu.GetGameVersions(), nil)
	selectionGameVersion.Selected = scu.GameVerLIVE

	logFilePath := binding.NewString()
	logFilePath.Set(filepath.Join(scu.GameDir, selectionGameVersion.Selected, "Game.log")) // scu.GameLogBackupDir, "Game Build(10275505) 19 Sep 25 (21 30 24).log")

	selectionGameVersion.OnChanged = func(s string) {
		logFilePath.Set(filepath.Join(scu.GameDir, selectionGameVersion.Selected, "Game.log"))
	}

	logFileEntry := widget.NewEntryWithData(logFilePath)

	var (
		btnEnabled bool
		cancelFn   context.CancelFunc
	)

	btnStartParse := widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), nil)

	runParser := func(cfg *logmon.Config, parseOption ParseOption) context.CancelFunc {
		ctx, cancel := context.WithCancel(context.Background())

		parser := logmon.GameParser(logmon.DefaultAD(), logmon.DefaultVD())
		out, errs := logmon.Run(ctx, cfg, parser)

		if parseOption == AggregateFiles {
			// TODO:
			// - Open file for writing
			// - Show loading bar
		}

		// Process in background to avoid blocking UI
		go func() {
			for out != nil || errs != nil {
				select {
				case s, ok := <-out:
					if !ok {
						out = nil
						continue
					}

					switch parseOption {
					case SingleFile, TailFile:

						// keep everything
						logItems = append(logItems, s)

						// (optional) cap memory
						if len(logItems) > logItemLimit {
							logItems = logItems[len(logItems)-logItemLimit:]
						}

						// append only if this item matches current predicate
						if filterMatch(s, activeFilter, activeQuery) {
							fyne.Do(func() {
								logBox.Add(parseLogResult(s))
								logBox.Refresh()
								scroll.ScrollToBottom()
							})
						}

					case AggregateFiles:
						// Write line to file

					}

				case e, ok := <-errs:
					if !ok {
						errs = nil
						continue
					}
					// non-fatal; log and keep going
					fmt.Fprintf(os.Stderr, "warn: %v\n", e)
				case <-ctx.Done():
					return
				}
			}

			// TODO: Check if this is still required for tailing
			// auto-reset button when stream finishes
			if parseOption == TailFile {
				fyne.Do(func() {
					btnEnabled = false
					btnStartParse.Icon = theme.MediaPlayIcon()
					btnStartParse.Text = "Start"
					btnStartParse.Refresh()
				})
			}

			if parseOption == AggregateFiles {
				// TODO: Maybe display file at this point
			}
		}()
		return cancel
	}
	btnFilterAll := widget.NewButton("All", nil)
	btnFilterHumans := widget.NewButton("Humans", nil)
	btnFilterVehicles := widget.NewButton("Vehicles", nil)

	checkBoxAggregate := widget.NewCheck("Aggregate", nil)
	btnOpenLogs := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {

		var (
			pathCh <-chan string
			errCh  <-chan error
			conf   *logmon.Config
		)
		aggragateLogs := checkBoxAggregate.Checked
		if aggragateLogs {
			pathCh, errCh = folderOpenPath(
				filepath.Join(scu.GameDir, selectionGameVersion.Selected, scu.GameLogBackupDir),
				win,
			)

			conf = &logmon.Config{
				Mode:      logmon.ModeOnce,
				FromStart: true,
				ChanSize:  1000,
			}
		} else {
			pathCh, errCh = fileOpenPath(
				filepath.Join(scu.GameDir, selectionGameVersion.Selected, scu.GameLogBackupDir),
				win,
				storage.NewExtensionFileFilter([]string{".log"}),
			)

			conf = &logmon.Config{
				Mode:      logmon.ModeOnce,
				FromStart: true,
				ChanSize:  256,
			}
		}

		// Reset state for a fresh file
		logItems = logItems[:0]
		activeQuery = ""
		activeFilter = None
		// Reset filter buttone states
		btnFilterAll.FocusGained()
		btnFilterHumans.FocusLost()
		btnFilterVehicles.FocusLost()

		rebuild() // clears view immediately

		go func() {
			select {
			case p := <-pathCh:
				fmt.Fprintf(os.Stdout, "--- PATH: %s ---\n", p)
				if aggragateLogs {
					conf.Archives = p
					_ = runParser(conf, AggregateFiles)
				} else {
					conf.ActivePath = p
					_ = runParser(conf, SingleFile)
				}

			case err := <-errCh:
				if err != nil {
					dialog.ShowError(err, win)
				}
			}
		}()
	})

	btnStartParse.OnTapped = func() {
		if btnEnabled {
			// Stop
			btnEnabled = false
			btnStartParse.Icon = theme.MediaPlayIcon()
			btnStartParse.Text = "Start"
			if cancelFn != nil {
				cancelFn()
			}
			btnOpenLogs.Enable()
			btnStartParse.Refresh()
			return
		}

		// Start
		logItems = logItems[:0]
		btnEnabled = true
		btnOpenLogs.Disable()
		btnStartParse.Icon = theme.MediaStopIcon()
		btnStartParse.Text = "Stop"
		btnStartParse.Refresh()

		logFilePathStr, err := logFilePath.Get()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		rebuild()

		cancelFn = runParser(&logmon.Config{
			ActivePath: logFilePathStr,
			Mode:       logmon.ModeTail,
			PollEvery:  5 * time.Second,
			FromStart:  true,
			ChanSize:   256,
		}, TailFile)
	}

	// top controlls
	btnFilterAll.FocusGained()
	btnFilterHumans.FocusLost()
	btnFilterVehicles.FocusLost()

	btnFilterAll.OnTapped = func() {
		activeFilter = None
		btnFilterAll.FocusGained()
		btnFilterHumans.FocusLost()
		btnFilterVehicles.FocusLost()
		rebuild()
	}

	btnFilterHumans.OnTapped = func() {
		activeFilter = Humans
		btnFilterHumans.FocusGained()
		btnFilterAll.FocusLost()
		btnFilterVehicles.FocusLost()
		rebuild()
	}

	btnFilterVehicles.OnTapped = func() {
		activeFilter = Vehicles
		btnFilterVehicles.FocusGained()
		btnFilterHumans.FocusLost()
		btnFilterAll.FocusLost()
		rebuild()
	}

	searchField := widget.NewEntry()
	searchField.PlaceHolder = "Search"
	btnSearch := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		q := strings.TrimSpace(searchField.Text)
		activeQuery = q
		if q == "" {
			// If query cleared, fall back to last non-search filter (or All)
			activeFilter = None
		} else {
			activeFilter = Search
		}
		rebuild()
	})

	top := container.NewBorder(nil, nil, container.NewHBox(btnFilterAll, btnFilterHumans, btnFilterVehicles), btnSearch, searchField)
	bottom := container.NewBorder(nil, nil, container.NewHBox(selectionGameVersion, btnStartParse), container.NewHBox(btnOpenLogs, checkBoxAggregate), logFileEntry)

	return widget.NewCard("", "", container.NewBorder(top, bottom, nil, nil, scroll))
}

type FilterOption int

const (
	None FilterOption = iota
	Humans
	Vehicles
	Search
)

func filterMatch(item *logmon.LogItem, opt FilterOption, q string) bool {
	switch opt {
	case None:
		return true
	case Humans:
		return item.Type == logmon.ActorDeath
	case Vehicles:
		return item.Type == logmon.VehicleDestruction
	case Search:
		q = strings.ToLower(strings.TrimSpace(q))
		if q == "" {
			return true
		}
		// add/remove fields as needed
		return strings.Contains(strings.ToLower(item.Attacker), q) ||
			strings.Contains(strings.ToLower(item.Victim), q) ||
			strings.Contains(strings.ToLower(item.Vehicle), q) ||
			strings.Contains(strings.ToLower(item.Weapon), q) ||
			strings.Contains(strings.ToLower(item.Location), q)
	default:
		return true
	}
}

func filterLogs(items []*logmon.LogItem, opt FilterOption, searchFields ...string) []fyne.CanvasObject {
	var filterToType = map[FilterOption]logmon.EventType{
		Humans:   logmon.ActorDeath,
		Vehicles: logmon.VehicleDestruction,
	}

	filtered := make([]fyne.CanvasObject, 0, len(items))

	switch opt {
	case None:
		for _, item := range items {
			filtered = append(filtered, parseLogResult(item))
		}
	case Humans, Vehicles:
		for _, item := range items {
			if item.Type == filterToType[opt] {
				filtered = append(filtered, parseLogResult(item))
			}
		}
	case Search:
		q := ""
		if len(searchFields) > 0 {
			q = strings.ToLower(searchFields[0])
		}
		if q == "" {
			// behave like None when query empty
			for _, item := range items {
				filtered = append(filtered, parseLogResult(item))
			}
			break
		}
		for _, item := range items {
			// Expand/adjust these fields as needed
			if strings.Contains(strings.ToLower(item.Attacker), q) ||
				strings.Contains(strings.ToLower(item.Victim), q) ||
				strings.Contains(strings.ToLower(item.Vehicle), q) ||
				strings.Contains(strings.ToLower(item.Weapon), q) ||
				strings.Contains(strings.ToLower(item.Location), q) {
				filtered = append(filtered, parseLogResult(item))
			}
		}
	}
	return filtered
}

func parseLogResult(i *logmon.LogItem) fyne.CanvasObject {
	switch i.Type {
	case logmon.ActorDeath:
		r := widget.NewRichText()
		r.Wrapping = fyne.TextWrapWord
		r.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text:  fmt.Sprintf("%s - ", i.Time),
				Style: widget.RichTextStyle{Inline: true},
			},
			&widget.TextSegment{
				Text: i.Attacker,
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					TextStyle: fyne.TextStyle{Bold: true},
					Inline:    true,
				},
			},
			&widget.TextSegment{
				Text:  " killed ",
				Style: widget.RichTextStyle{Inline: true},
			},
			&widget.TextSegment{
				Text: i.Victim,
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNamePrimary,
					TextStyle: fyne.TextStyle{Bold: true},
					Inline:    true,
				},
			},
			&widget.TextSegment{
				Text:  " using ",
				Style: widget.RichTextStyle{Inline: true},
			},
			&widget.TextSegment{
				Text: i.Weapon,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
				},
			},
		}
		return r
	case logmon.VehicleDestruction:
		r := widget.NewRichText()
		r.Wrapping = fyne.TextWrapWord
		r.Segments = []widget.RichTextSegment{
			&widget.TextSegment{
				Text:  fmt.Sprintf("%s - ", i.Time),
				Style: widget.RichTextStyle{Inline: true},
			},
			&widget.TextSegment{
				Text: i.Attacker,
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameError,
					TextStyle: fyne.TextStyle{Bold: true},
					Inline:    true,
				},
			},
			&widget.TextSegment{
				Text:  " destroyed ",
				Style: widget.RichTextStyle{Inline: true},
			},
			&widget.TextSegment{
				Text: i.Vehicle,
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNamePrimary,
					TextStyle: fyne.TextStyle{Bold: true},
					Inline:    true,
				},
			},
			&widget.TextSegment{
				Text:  " using ",
				Style: widget.RichTextStyle{Inline: true},
			},
			&widget.TextSegment{
				Text: i.Weapon,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
					Inline:    true,
				},
			},
			&widget.TextSegment{
				Text:  " at ",
				Style: widget.RichTextStyle{Inline: true},
			},
			&widget.TextSegment{
				Text: i.Location,
				Style: widget.RichTextStyle{
					TextStyle: fyne.TextStyle{Bold: true},
					Inline:    true,
				},
			},
		}
		return r
	default:
		r := widget.NewRichText()
		r.Wrapping = fyne.TextWrapWord
		r.Segments = []widget.RichTextSegment{
			&widget.TextSegment{Text: fmt.Sprintf("Unknown Event: %s", i.Type)},
		}
		return r
	}
}

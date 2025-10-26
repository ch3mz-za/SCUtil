package frontend

import (
	"bufio"
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
	"github.com/ch3mz-za/SCUtil/internal/util"
)

const (
	logItemLimit          int    = 500
	logItemAggregateLimit int    = 5000
	prefixAggregateLog    string = "Aggregated"
)

type ParseOption int

const (
	TailFile ParseOption = iota
	SingleFile
	AggregateFiles
)

type FilterOption int

const (
	None FilterOption = iota
	Humans
	Vehicles
	Search
)

// logViewerState holds the state for the log viewer
type logViewerState struct {
	logItems      []*logmon.LogItem
	filteredItems []*logmon.LogItem // Filtered subset for display
	activeFilter  FilterOption
	activeQuery   string
	cancelFn      context.CancelFunc
	btnEnabled    bool
}

// logViewCtx holds the log viewer's visual object as context for various functions
type logViewCtx struct {
	win         fyne.Window
	logList     *widget.List
	progressBar *dialog.CustomDialog
}

// startLogParser initializes and runs the log parser in the background
func startLogParser(
	lvCtx *logViewCtx,
	state *logViewerState,
	cfg *logmon.Config,
	parseOption ParseOption,
) (context.CancelFunc, string, error) {
	ctx, cancel := context.WithCancel(context.Background())

	parser := logmon.GameParser(logmon.DefaultAD(), logmon.DefaultVD())
	out, errs := logmon.Run(ctx, cfg, parser)

	var (
		f *os.File
		w *bufio.Writer
	)

	// Setup for background processing
	var filePath string
	switch parseOption {
	case SingleFile, TailFile:
		state.logItems = make([]*logmon.LogItem, 0, logItemLimit)
	case AggregateFiles:
		state.logItems = make([]*logmon.LogItem, 0, logItemAggregateLimit)

		destDir := filepath.Join(scu.AppDir, scu.AggregatedLogsDir)
		if err := util.CreateDirIfNotExist(destDir); err != nil {
			return cancel, "", err
		}

		var err error
		filePath = filepath.Join(destDir, generateLogName())
		f, err = os.OpenFile(
			filePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0644,
		)
		if err != nil {
			return cancel, "", err
		}
		w = bufio.NewWriter(f)
	}

	// Process in background to avoid blocking UI
	go func() {
		defer func() {
			if parseOption == AggregateFiles {
				if w != nil {
					_ = w.Flush()
				}
				if f != nil {
					_ = f.Close()
				}
				fyne.Do(func() {
					if lvCtx.progressBar != nil {
						lvCtx.progressBar.Hide()
					}
				})
			}
		}()

		for out != nil || errs != nil {
			select {
			case s, ok := <-out:
				if !ok {
					out = nil
					continue
				}

				switch parseOption {
				case SingleFile, TailFile:
					state.logItems = append(state.logItems, s)

					// Cap memory
					if len(state.logItems) > logItemLimit {
						state.logItems = state.logItems[len(state.logItems)-logItemLimit:]
						// Also cap filtered items if needed
						if len(state.filteredItems) > logItemLimit {
							state.filteredItems = state.filteredItems[len(state.filteredItems)-logItemLimit:]
						}
					}

					// Append to filtered list if this item matches current predicate
					if filterMatch(s, state.activeFilter, state.activeQuery) {
						state.filteredItems = append(state.filteredItems, s)
						fyne.Do(func() {
							lvCtx.logList.Refresh()
							lvCtx.logList.ScrollToBottom()
						})
					}

				case AggregateFiles:
					res := parseLogResult(s)
					resStr := res.(*widget.RichText).String()

					if w != nil {
						_, err := w.WriteString(resStr + "\n")
						if err != nil {
							fyne.Do(func() { dialog.ShowError(err, lvCtx.win) })
						}
					}
				}

			case e, ok := <-errs:
				if !ok {
					errs = nil
					continue
				}
				// Non-fatal; log and keep going
				fmt.Fprintf(os.Stderr, "warn: %v\n", e)
			case <-ctx.Done():
				return
			}
		}
	}()

	return cancel, filePath, nil
}

func logs(win fyne.Window) fyne.CanvasObject {
	// Initialize state
	state := &logViewerState{
		logItems:      make([]*logmon.LogItem, 0, logItemLimit),
		filteredItems: make([]*logmon.LogItem, 0, logItemLimit),
		activeFilter:  None,
		activeQuery:   "",
	}

	// Create virtualized list for high performance
	logList := widget.NewList(
		// Length function
		func() int {
			return len(state.filteredItems)
		},
		// CreateItem - returns a template RichText widget
		func() fyne.CanvasObject {
			r := widget.NewRichText()
			r.Wrapping = fyne.TextWrap(fyne.TextTruncateClip) // Single line, no wrapping
			return r
		},
		// UpdateItem - updates the template with data for index i
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			if i >= len(state.filteredItems) {
				return
			}
			r := obj.(*widget.RichText)
			item := state.filteredItems[i]

			// Reuse parseLogResult logic by extracting segments
			switch item.Type {
			case logmon.ActorDeath:
				r.Segments = []widget.RichTextSegment{
					newTextSegment(fmt.Sprintf("%s [AD] ", formatTimestamp(item.Time))),
					newColoredBoldSegment(item.Attacker, theme.ColorNameError),
					newTextSegment(" killed "),
					newColoredBoldSegment(item.Victim, theme.ColorNamePrimary),
					newTextSegment(" using "),
					newBoldSegment(item.Weapon),
				}

			case logmon.VehicleDestruction:
				r.Segments = []widget.RichTextSegment{
					newTextSegment(fmt.Sprintf("%s [VD] ", formatTimestamp(item.Time))),
					newColoredBoldSegment(item.Attacker, theme.ColorNameError),
					newTextSegment(" destroyed "),
					newColoredBoldSegment(item.Vehicle, theme.ColorNamePrimary),
					newTextSegment(" using "),
					newBoldSegment(item.Weapon),
					newTextSegment(" at "),
					newBoldSegment(item.Location),
				}

			default:
				r.Segments = []widget.RichTextSegment{
					newTextSegment(fmt.Sprintf("Unknown Event: %s", item.Type)),
				}
			}
			r.Refresh()
		},
	)

	// Rebuild the filtered list (call this when filter/search changes)
	rebuild := func() {
		state.filteredItems = state.filteredItems[:0]
		for _, it := range state.logItems {
			if filterMatch(it, state.activeFilter, state.activeQuery) {
				state.filteredItems = append(state.filteredItems, it)
			}
		}
		fyne.Do(func() {
			logList.Refresh()
			logList.ScrollToBottom()
		})
	}

	// Bottom controls
	selectionGameVersion := widget.NewSelect(scu.GetGameVersions(), nil)
	selectionGameVersion.Selected = scu.GameVerLIVE

	logFilePath := binding.NewString()
	logFilePath.Set(filepath.Join(scu.GameDir, selectionGameVersion.Selected, "Game.log")) // scu.GameLogBackupDir, "Game Build(10275505) 19 Sep 25 (21 30 24).log")

	selectionGameVersion.OnChanged = func(s string) {
		logFilePath.Set(getGameLogPath(selectionGameVersion))
	}

	logFileEntry := widget.NewEntryWithData(logFilePath)
	logFileEntry.Validator = nil

	btnStartParse := widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), nil)

	progressBar := dialog.NewCustomWithoutButtons("Loading...", widget.NewProgressBarInfinite(), win)

	lvCtx := &logViewCtx{
		win:         win,
		logList:     logList,
		progressBar: progressBar,
	}

	runParser := func(cfg *logmon.Config, parseOption ParseOption) (context.CancelFunc, string, error) {
		return startLogParser(lvCtx, state, cfg, parseOption)
	}
	btnFilterAll := widget.NewButton("All", nil)
	btnFilterHumans := widget.NewButton("Humans", nil)
	btnFilterVehicles := widget.NewButton("Vehicles", nil)

	btnAggregate := widget.NewButton("Aggregate Logs", func() {
		conf := &logmon.Config{
			Mode:      logmon.ModeOnce,
			FromStart: true,
			ChanSize:  1000,
			Archives:  filepath.Join(scu.GameDir, selectionGameVersion.Selected, scu.GameLogBackupDir),
		}
		_, filePath, err := runParser(conf, AggregateFiles)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		dialog.ShowConfirm("Open aggregated file", "Open aggregated file?", func(open bool) {
			if open {
				readAggregatedLog(filePath, lvCtx, state)
			}
		}, win)
	})

	// Open logs
	btnOpenLogs := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		var (
			pathCh <-chan string
			errCh  <-chan error
			conf   *logmon.Config
		)

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

		// Reset state for a fresh file
		state.logItems = state.logItems[:0]
		state.activeQuery = ""
		state.activeFilter = None
		// Reset filter button states
		btnFilterAll.FocusGained()
		btnFilterHumans.FocusLost()
		btnFilterVehicles.FocusLost()

		rebuild() // clears view immediately

		go func() {
			select {
			case p := <-pathCh:
				logFilePath.Set(p)
				f, err := os.Stat(p)
				if err != nil {
					dialog.ShowError(err, win)
				}

				if strings.HasPrefix(f.Name(), prefixAggregateLog) {
					// read aggregated logs
					readAggregatedLog(p, lvCtx, state)
				} else {
					conf.ActivePath = p
					_, _, err = runParser(conf, SingleFile)
					if err != nil {
						dialog.ShowError(err, win)
					}
				}

			case err := <-errCh:
				if err != nil {
					dialog.ShowError(err, win)
				}
			}
		}()
	})

	// Start live tailing of logs
	btnStartParse.OnTapped = func() {
		if state.btnEnabled {
			// Stop
			state.btnEnabled = false
			btnStartParse.Icon = theme.MediaPlayIcon()
			btnStartParse.Text = "Start"
			if state.cancelFn != nil {
				state.cancelFn()
			}
			btnOpenLogs.Enable()
			btnStartParse.Refresh()
			return
		}

		// Start
		state.logItems = state.logItems[:0]
		state.btnEnabled = true
		btnOpenLogs.Disable()
		btnStartParse.Icon = theme.MediaStopIcon()
		btnStartParse.Text = "Stop"
		btnStartParse.Refresh()

		logFilePath.Set(getGameLogPath(selectionGameVersion))
		logFilePathStr, err := logFilePath.Get()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		// Clear filtered items for fresh start
		state.filteredItems = state.filteredItems[:0]
		logList.Refresh()

		state.cancelFn, _, err = runParser(&logmon.Config{
			ActivePath: logFilePathStr,
			Mode:       logmon.ModeTail,
			PollEvery:  5 * time.Second,
			FromStart:  true,
			ChanSize:   256,
		}, TailFile)
		if err != nil {
			dialog.ShowError(err, win)
		}
	}

	// top controlls
	btnFilterAll.FocusGained()
	btnFilterHumans.FocusLost()
	btnFilterVehicles.FocusLost()

	btnFilterAll.OnTapped = func() {
		state.activeFilter = None
		btnFilterAll.FocusGained()
		btnFilterHumans.FocusLost()
		btnFilterVehicles.FocusLost()
		rebuild()
	}

	btnFilterHumans.OnTapped = func() {
		state.activeFilter = Humans
		btnFilterHumans.FocusGained()
		btnFilterAll.FocusLost()
		btnFilterVehicles.FocusLost()
		rebuild()
	}

	btnFilterVehicles.OnTapped = func() {
		state.activeFilter = Vehicles
		btnFilterVehicles.FocusGained()
		btnFilterHumans.FocusLost()
		btnFilterAll.FocusLost()
		rebuild()
	}

	searchField := widget.NewEntry()
	searchField.PlaceHolder = "Search"
	btnSearch := widget.NewButtonWithIcon("", theme.SearchIcon(), func() {
		q := strings.TrimSpace(searchField.Text)
		state.activeQuery = q
		if q == "" {
			// If query cleared, fall back to last non-search filter (or All)
			state.activeFilter = None
		} else {
			state.activeFilter = Search
		}
		rebuild()
	})

	top := container.NewBorder(nil, nil, container.NewHBox(btnFilterAll, btnFilterHumans, btnFilterVehicles), btnSearch, searchField)
	bottom := container.NewBorder(nil, nil, container.NewHBox(selectionGameVersion, btnStartParse), container.NewHBox(btnOpenLogs, btnAggregate), logFileEntry)

	return widget.NewCard("", "", container.NewBorder(top, bottom, nil, nil, logList))
}

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

var (
	AggregateAD logmon.IndicesAD = logmon.IndicesAD{Time: 0, Victim: 4, Attacker: 2, Weapon: 6}
	AggregateVD logmon.IndicesVD = logmon.IndicesVD{Time: 0, Vehicle: 4, Location: 8, Driver: -1, Attacker: 2, Weapon: 6}
)

func stringToLogItem(line string) *logmon.LogItem {
	const (
		LogTypeIndex int    = 1
		ActorDeath   string = "[AD]"
		VehicleDeath string = "[VD]"
	)

	if line == "" {
		return nil
	}

	var fields []string
	get := func(i int) string {
		if i < 0 || i >= len(fields) {
			return ""
		}
		return fields[i]
	}

	fields = strings.Fields(line)
	switch fields[LogTypeIndex] {
	case ActorDeath:
		return &logmon.LogItem{
			Time:     logmon.RoundTimeToSeconds(get(AggregateAD.Time)),
			Attacker: get(AggregateAD.Attacker),
			Victim:   get(AggregateAD.Victim),
			Weapon:   get(AggregateAD.Weapon),
			Type:     logmon.ActorDeath,
		}
	case VehicleDeath:
		return &logmon.LogItem{
			Time:     logmon.RoundTimeToSeconds(get(AggregateVD.Time)),
			Vehicle:  get(AggregateVD.Vehicle),
			Attacker: get(AggregateVD.Attacker),
			Location: get(AggregateVD.Location),
			Driver:   get(AggregateVD.Driver),
			Weapon:   get(AggregateVD.Weapon),
			Type:     logmon.VehicleDestruction,
		}
	}

	return nil
}

func readAggregatedLog(filePath string, lvCtx *logViewCtx, state *logViewerState) {
	f, err := os.Open(filePath)
	if err != nil {
		dialog.ShowError(err, lvCtx.win)
		return
	}
	defer f.Close()

	lvCtx.progressBar.Show()

	// Load all items into memory first (no UI updates during scan)
	s := bufio.NewScanner(f)
	for s.Scan() {
		item := stringToLogItem(s.Text())
		state.logItems = append(state.logItems, item)

		// Cap memory
		if len(state.logItems) > logItemAggregateLimit {
			state.logItems = state.logItems[len(state.logItems)-logItemAggregateLimit:]
		}
	}

	// Rebuild filtered list and update UI
	state.filteredItems = state.filteredItems[:0]
	for _, it := range state.logItems {
		if filterMatch(it, state.activeFilter, state.activeQuery) {
			state.filteredItems = append(state.filteredItems, it)
		}
	}

	fyne.Do(func() {
		lvCtx.progressBar.Dismiss()
		lvCtx.logList.Refresh()
		lvCtx.logList.ScrollToBottom()
	})
}

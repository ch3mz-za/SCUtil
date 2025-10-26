package logmon

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Mode int

const (
	ModeOnce Mode = iota
	ModeTail
	ModeBoth
)

func ParseMode(s string) (Mode, error) {
	switch strings.ToLower(s) {
	case "once":
		return ModeOnce, nil
	case "tail":
		return ModeTail, nil
	case "both":
		return ModeBoth, nil
	default:
		return 0, fmt.Errorf("invalid mode %q (use once|tail|both)", s)
	}
}

type Config struct {
	ActivePath string
	Archives   string
	Mode       Mode
	FromStart  bool
	PollEvery  time.Duration // default 2s if zero
	ChanSize   int           // buffer for out/errs; default 100 if zero
}

// Parser returns (formatted, true) if the line should be emitted.
type Parser func(raw string) (*LogItem, bool)

// Run starts processing and returns two read-only channels:
// - out: parsed/“resulting” strings
// - errs: non-fatal errors encountered while running
// Channels close on ctx cancel or completion.
// NOTE: On fatal open errors, errs gets one item and then both close.
func Run(ctx context.Context, cfg *Config, parser Parser) (<-chan *LogItem, <-chan error) {
	if cfg.PollEvery <= 0 {
		cfg.PollEvery = 2 * time.Second
	}
	if cfg.ChanSize <= 0 {
		cfg.ChanSize = 100
	}
	out := make(chan *LogItem, cfg.ChanSize)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		// Helper to safely emit parsed results
		emit := func(line string) {
			if s, ok := parser(line); ok && s != nil {
				select {
				case out <- s:
				case <-ctx.Done():
				}
			}
		}

		// 1) ONCE mode part (historical pass)
		readOnce := func() error {
			// archives (optional)
			if strings.TrimSpace(cfg.Archives) != "" {
				matches, err := filepath.Glob(cfg.Archives + "/*.log")
				if err != nil {
					return fmt.Errorf("archives glob: %w", err)
				}

				sort.Slice(matches, func(i, j int) bool {
					fi, err1 := os.Stat(matches[i])
					fj, err2 := os.Stat(matches[j])
					if err1 != nil || err2 != nil {
						// fallback to filename order if we can’t stat
						return matches[i] < matches[j]
					}
					return fi.ModTime().Before(fj.ModTime())
				})

				for _, p := range matches {
					if err := readFileOnce(ctx, p, emit); err != nil {
						return err
					}
				}
			}
			// active file once
			return readFileOnce(ctx, cfg.ActivePath, emit)
		}

		switch cfg.Mode {
		case ModeOnce:
			if err := readOnce(); err != nil {
				select {
				case errs <- err:
				default:
				}
			}
			return

		case ModeTail, ModeBoth:
			if cfg.Mode == ModeBoth {
				if err := readOnce(); err != nil {
					select {
					case errs <- err:
					default:
					}
					// continue to tail anyway
				}
			}
			if err := tail(ctx, cfg, emit); err != nil && !errors.Is(err, context.Canceled) {
				select {
				case errs <- err:
				default:
				}
			}
			return

		default:
			select {
			case errs <- fmt.Errorf("unsupported mode: %v", cfg.Mode):
			default:
			}
			return
		}
	}()

	return out, errs
}

func readFileOnce(ctx context.Context, path string, emit func(string)) error {
	f, err := os.Open(path)
	if err != nil {
		// Missing active file isn't necessarily fatal in some setups.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 10*1024*1024) // up to 10MB lines
	for sc.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			emit(sc.Text())
		}
	}
	return sc.Err()
}

func tail(ctx context.Context, cfg *Config, emit func(string)) error {
	// Open (wait if not present yet)
	f, err := os.Open(cfg.ActivePath)
	if err != nil {
		for {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if !os.IsNotExist(err) {
				return err
			}
			time.Sleep(cfg.PollEvery)
			f, err = os.Open(cfg.ActivePath)
			if err == nil {
				break
			}
		}
	}
	defer f.Close()

	// Position
	if cfg.FromStart {
		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return err
		}
	} else {
		if _, err := f.Seek(0, io.SeekEnd); err != nil {
			return err
		}
	}

	reader := bufio.NewReader(f)

	for {
		// Drain available lines
		drained := false
		for {
			line, err := reader.ReadString('\n')
			switch {
			case err == nil:
				emit(strings.TrimRight(line, "\r\n"))
				drained = true
			case errors.Is(err, io.EOF):
				goto afterDrain
			default:
				// unexpected read error → try reopen
				reader, f = reopen(cfg.ActivePath, reader, f)
				goto afterDrain
			}
		}
	afterDrain:

		// Rotation/truncation?
		if isProbablyRotated(f, cfg.ActivePath) {
			reader, f = reopen(cfg.ActivePath, reader, f)
		}

		// If we drained, loop again immediately to keep up
		if drained {
			continue
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(cfg.PollEvery):
		}
	}
}

func reopen(path string, _ *bufio.Reader, f *os.File) (*bufio.Reader, *os.File) {
	_ = f.Close()
	newF, err := os.Open(path)
	if err != nil {
		// If reopen fails, keep old handle (already closed), but create a dummy reader to avoid nil deref.
		// Caller will try again on next loop.
		return bufio.NewReader(strings.NewReader("")), f
	}
	// After rotation, seek to end to avoid replaying content.
	if _, err := newF.Seek(0, io.SeekEnd); err != nil {
		_, _ = newF.Seek(0, io.SeekStart)
	}
	return bufio.NewReader(newF), newF
}

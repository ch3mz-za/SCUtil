package logmon

import (
	"io"
	"os"
	"runtime"
	"time"
)

func isProbablyRotated(f *os.File, path string) bool {
	fi1, err := f.Stat()
	if err != nil {
		return false
	}
	fi2, err := os.Stat(path)
	if err != nil {
		return false
	}

	// truncated?
	off, _ := f.Seek(0, io.SeekCurrent)
	if fi2.Size() < off {
		return true
	}

	switch runtime.GOOS {
	case "linux", "darwin", "freebsd":
		return inodeChanged(fi1, fi2)
	default:
		// heuristic for windows/others
		if fi2.ModTime().After(fi1.ModTime().Add(2*time.Second)) && fi2.Size() < fi1.Size() {
			return true
		}
		return false
	}
}

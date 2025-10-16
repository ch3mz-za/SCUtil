//go:build windows || !unix

package logmon

import "os"

// inodeChanged is a no-op stub for non-Unix systems.
func inodeChanged(_, _ os.FileInfo) bool {
	return false
}

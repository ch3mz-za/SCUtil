//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package logmon

import (
	"os"

	"golang.org/x/sys/unix"
)

func inodeChanged(fi1, fi2 os.FileInfo) bool {
	s1 := fi1.Sys().(*unix.Stat_t)
	s2 := fi2.Sys().(*unix.Stat_t)
	return s1.Ino != s2.Ino || s1.Dev != s2.Dev
}

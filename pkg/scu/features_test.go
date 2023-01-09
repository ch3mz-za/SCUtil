package scu

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var liveBase = filepath.Join("../", "test", "Star Citizen", "LIVE")

func SetupTestFolder() {
	os.MkdirAll(filepath.Join(liveBase, "USER"), os.ModeAppend)
	os.Create(filepath.Join(liveBase, "Data.p4k"))
	os.Create(filepath.Join(liveBase, "test_file1.txt"))
	os.Create(filepath.Join(liveBase, "test_file1.txt"))
}

func TestClearAllDataExceptP4k(t *testing.T) {
	SetupTestFolder()
	clearAllDataExceptP4k()

	// List all files
	// PASS IF len(files) == 1 && file hasSuffix(".p4k")
	files, _ := os.ReadDir(liveBase)

	require.Equal(t, 1, len(files))
}

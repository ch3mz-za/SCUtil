package scu

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ch3mz-za/SCUtil/internal/common"
	p4k "github.com/ch3mz-za/SCUtil/internal/p4kReader"
)

// GetP4kFilenames - Gets all the filenames from the Data.p4k file and writes them to a specific folder
func GetP4kFilenames(version string) error {
	gameDir := filepath.Join(GameDir, version)
	resultsDir := filepath.Join(AppDir, fmt.Sprintf(P4kFilenameResultsDir, version))
	return p4k.GetP4kFilenames(gameDir, resultsDir)
}

// SearchP4kFilenames - Search for specific filenames within the Data.p4k file
func SearchP4kFilenames(version, phrase string) error {
	gameDir := filepath.Join(GameDir, version)
	filename := strings.ReplaceAll(phrase, "\\", "_") + ".txt"
	resultsDir := filepath.Join(AppDir, P4kSearchResultsDir, version)
	common.MakeDir(resultsDir)
	resultsDir = filepath.Join(resultsDir, filename)

	err := p4k.SearchP4kFilenames(gameDir, phrase, resultsDir)
	if err != nil {
		return err
	}

	return nil
}

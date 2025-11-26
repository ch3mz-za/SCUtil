package scu

import (
	"time"
)

const (
	ctrlMapFileExt string        = ".xml"
	twoSecondDur   time.Duration = 2 * time.Second

	// Directories
	CharactersBackupDir      string = "Backups/CustomCharacters"
	CharactersDir            string = UserDir + "/Client/0/CustomCharacters"
	ControlMappingsBackupDir string = "Backups/ControlMappings"
	ControlMappingsDir       string = UserDir + "/Client/0/Controls/Mappings"
	GameLogDir               string = "Game.log"
	GameLogBackupDir         string = "logbackups"
	AggregatedLogsDir        string = "Backups/AggregatedLogs"
	P4kFilenameResultsDir    string = "P4kResults/AllFileNames/%s/AllP4kFilenames.txt"
	P4kSearchResultsDir      string = "P4kResults/Searches"
	ScreenshotsBackupDir     string = "Backups/Screenshots"
	ScreenshotsDir           string = "ScreenShots"
	UserBackupDir            string = "Backups/UserFolder"
	UserDir                  string = "USER"
)

var (
	GameDir string = ""
	AppDir  string = ""
)

// SetGameDir updates the global GameDir variable
// This should be called whenever the game directory changes
func SetGameDir(dir string) {
	GameDir = dir
}

// GetGameDir returns the current game directory
func GetGameDir() string {
	return GameDir
}

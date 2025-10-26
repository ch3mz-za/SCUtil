package scu

import (
	"time"
)

const (
	ctrlMapFileExt string        = ".xml"
	twoSecondDur   time.Duration = 2 * time.Second

	// Directories
	CharactersBackupDir      string = "BACKUPS/CustomCharacters"
	CharactersDir            string = UserDir + "/Client/0/CustomCharacters"
	ControlMappingsBackupDir string = "BACKUPS/ControlMappings"
	ControlMappingsDir       string = UserDir + "/Client/0/Controls/Mappings"
	GameLogDir               string = "Game.log"
	GameLogBackupDir         string = "logbackups"
	AggregatedLogsDir        string = "BACkUPS/AggregatedLogs"
	P4kFilenameResultsDir    string = "P4kResults/AllFileNames/%s/AllP4kFilenames.txt"
	P4kSearchResultsDir      string = "P4kResults/Searches"
	ScreenshotsBackupDir     string = "BACKUPS/Screenshots"
	ScreenshotsDir           string = "ScreenShots"
	UserBackupDir            string = "BACKUPS/UserFolder"
	UserDir                  string = "USER"
)

var (
	GameDir string = ""
	AppDir  string = ""
)

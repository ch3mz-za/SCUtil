package scu

import "path/filepath"

// BackupScreenshots - Backup all screenshots for specific game version
func BackupScreenshots(version string) error {
	screenshotDir := filepath.Join(GameDir, version, ScreenshotsDir)
	backupDir := filepath.Join(AppDir, ScreenshotsBackupDir, version)
	return backupFiles(screenshotDir, backupDir, false, ".jpg")
}

// BackupUserDirectory - Backup the USER directory
func BackupUserDirectory(version string) error {
	userDir := filepath.Join(GameDir, version, UserDir)
	backupDir := filepath.Join(AppDir, UserBackupDir, version, UserDir)
	return BackupDirectory(userDir, backupDir)
}

// BackupUserCharacters - Backup the custom characters in the USER directory
func BackupUserCharacters(version string) error {
	charDir := filepath.Join(GameDir, version, CharactersDir)
	backupDir := filepath.Join(AppDir, CharactersBackupDir, version)
	return backupFiles(charDir, backupDir, false, ".chf")
}

// BackupControlMappings - Backup game control mappings
func BackupControlMappings(version string) error {
	mappingsDir := filepath.Join(GameDir, version, ControlMappingsDir)
	backupDir := filepath.Join(AppDir, ControlMappingsBackupDir, version)
	return backupFiles(mappingsDir, backupDir, true, ctrlMapFileExt)
}

// GetBackedUpControlMappings - Retrieve a list of all the backed-up control mappings
func GetBackedUpControlMappings(version string) ([]string, error) {
	return filepath.Glob(filepath.Join(AppDir, ControlMappingsBackupDir, version))
}

// RestoreControlMappings - Restores a specified control mapping for a specific game version
func RestoreControlMappings(version string, filename string) error {
	mappingsFilePath := filepath.Join(GameDir, version, ControlMappingsDir, filename)
	backupFilePath := filepath.Join(AppDir, ControlMappingsBackupDir, version, filename)
	return restoreFile(backupFilePath, mappingsFilePath, true)
}

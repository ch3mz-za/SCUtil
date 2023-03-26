package scu

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ch3mz-za/SCUtil/pkg/common"
)

const (
	GameVerLIVE string = "LIVE"
	GameVerPTU  string = "PTU"
)

func restoreFiles(sourceDir, destDir, filename string) error {

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return errors.New("unable to create control mappings directory")
		}
	}

	restoreFileName := filename
	split := strings.Split(strings.TrimSuffix(restoreFileName, ctrlMapFileExt), "-")

	// Remove backup timestamp
	if splitCnt := len(split); splitCnt >= 3 {
		restoreFileName = strings.Join(split[:splitCnt-2], "-") + ctrlMapFileExt
	}

	if err := common.CopyFile(
		filepath.Join(sourceDir, string(filename)), // src
		filepath.Join(destDir, restoreFileName),    // dst
	); err != nil {
		return fmt.Errorf("restore error:\n %s", err.Error())
	}
	return nil
}

func backupFiles(sourceDir, destDir string, addTimestamp bool, filetypes ...string) error {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("unable to open source directory %s", sourceDir)
	}

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("unable to create destination directory: %s", destDir)
		}
		println("Destination directory created: " + destDir)
	}

	var destFilename string
	fileBackupCount := 0
	for _, fType := range filetypes {
		for _, f := range files {
			if strings.HasSuffix(strings.ToLower(f.Name()), fType) {

				if addTimestamp {
					destFilename = strings.TrimSuffix(f.Name(), fType) + "-" + time.Now().Format("2006.01.02-15.04.05") + fType
				} else {
					destFilename = f.Name()
				}

				fileBackupCount++
				if err := common.CopyFile(
					filepath.Join(sourceDir, f.Name()),   // src
					filepath.Join(destDir, destFilename), // dst
				); err != nil {
					fmt.Printf("copy error: %s\n", err.Error())
				} else {
					fmt.Printf("file copied: %s\n", destFilename)
				}

			}
		}
	}

	if fileBackupCount == 0 {
		return errors.New("no files found")
	}
	return nil
}

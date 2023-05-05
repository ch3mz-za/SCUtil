package scu

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
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

// backupFiles
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

func BackupDirectory(sourceDir, destDir string) error {

	return filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// copy to this path
		outpath := filepath.Join(destDir, strings.TrimPrefix(path, sourceDir))

		if info.IsDir() {
			os.MkdirAll(outpath, info.Mode())
			return nil // means recursive
		}

		// handle irregular files
		if !info.Mode().IsRegular() {
			switch info.Mode().Type() & os.ModeType {
			case os.ModeSymlink:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}
				return os.Symlink(link, outpath)
			}
			return nil
		}

		// copy contents of regular file efficiently

		// open input
		in, _ := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		// create output
		fh, err := os.Create(outpath)
		if err != nil {
			return err
		}
		defer fh.Close()

		// make it the same
		fh.Chmod(info.Mode())

		// copy content
		_, err = io.Copy(fh, in)
		return err
	})
}

func FindGameDirectory(searchDir string) string {
	gameDir, err := common.FindDir(searchDir, filepath.Join("StarCitizen", GameVerLIVE))
	if err != nil && gameDir != "" {
		return ""
	}
	return filepath.Dir(gameDir)
}

func IsGameDirectory(gameDir string) bool {
	entries, err := os.ReadDir(gameDir)
	if err != nil {
		return false
	}

	for _, f := range entries {
		if f.IsDir() && (f.Name() == GameVerLIVE || f.Name() == GameVerPTU) {
			return true
		}
	}
	return false
}

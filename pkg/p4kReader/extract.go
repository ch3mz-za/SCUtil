package p4k

import (
	"archive/zip"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ExtractFileFromP4k(p4kFilePath, filename, outputDir string) error {
	// Open the .p4k file
	r, err := zip.OpenReader(p4kFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0755)
	}

	// Iterate through the files in the .p4k archive
	for _, f := range r.File {
		if f.Name != filename {
			continue
		}
		log.Printf("filename: %s | f.Name %s\n", filename, f.Name)

		// Open the file in the archive
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create the output file
		outputPath := filepath.Join(outputDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(outputPath, f.Mode())
			return errors.New("this is a direcotry")
		}
		log.Println("outputpath: " + outputPath)
		outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outputFile.Close()

		// Copy the file data from the archive to the output file
		_, err = io.Copy(outputFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}

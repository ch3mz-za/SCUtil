package p4k

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const outputDir string = "./extracted_p4k_files"

func searchFile(filename string, files []*zip.File, start, end int) {
	defer wg.Done()
	for _, f := range files[start:end] {
		if f.Name == filename {
			// extract file
			fmt.Printf("File found! Extracting...\n")
			fileToExtract, err := f.Open()
			if err != nil {
				fmt.Printf("Unable to open zipped file: %s\n", err.Error())
				return
			}
			defer fileToExtract.Close()

			MakeDir(outputDir)

			fpath := filepath.Join(outputDir, filepath.Base(filename))
			outputFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				fmt.Printf("Unable to open file: %s\n", err.Error())
				return
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, fileToExtract)
			if err != nil {
				fmt.Printf("Unable to copy file content: %s\n", err.Error())
				return
			}
		}
	}
}

func ExtractP4kFile(gameDir, filename string) {
	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		fmt.Printf("Unable to open p4k data file: %s\n", err.Error())
	}

	div := runtime.NumCPU()
	fileCount := len(r.File)
	if fileCount > 1000 {
		interval := fileCount / div
		for i := 0; i < div; i++ {
			wg.Add(1)
			if i == div-1 {
				go searchFile(filename, r.File, interval*i, fileCount)
			} else {
				go searchFile(filename, r.File, interval*i, interval*(i+1))
			}
		}
	} else {
		wg.Add(1)
		buildFileNamesFile(fmt.Sprintf("%s%d.txt", filename, 0), r.File, 0, fileCount)
	}
	wg.Wait()
}

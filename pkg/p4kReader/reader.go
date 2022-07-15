package p4k

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

const dataP4k = "Data.p4k"

func buildFileNamesFile(filename string, files []*zip.File, start, end int) {
	defer wg.Done()
	p4kFiles, err := os.Create(filename)
	if err != nil {
		return
	}
	defer p4kFiles.Close()

	for _, f := range files[start:end] {
		p4kFiles.WriteString(f.Name + "\n")
	}
}

func GetP4kFilenames(gameDir, outputDir string) {

	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		fmt.Printf("Unable to open p4k data file: %s", err.Error())
	}
	defer r.Close()

	MakeDir(outputDir)

	div := runtime.NumCPU()
	fileCount := len(r.File)
	filename := filepath.Join(outputDir, "P4k_filenames_")

	if fileCount > 1000 {
		interval := fileCount / div
		for i := 0; i < div; i++ {
			wg.Add(1)
			if i == div-1 {
				go buildFileNamesFile(fmt.Sprintf("%s%d.txt", filename, i+1), r.File, interval*i, fileCount)
			} else {
				go buildFileNamesFile(fmt.Sprintf("%s%d.txt", filename, i+1), r.File, interval*i, interval*(i+1))
			}
		}
		wg.Wait()
	} else {
		wg.Add(1)
		buildFileNamesFile(fmt.Sprintf("%s%d.txt", filename, 0), r.File, 0, fileCount)
	}

}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to create directory: %s - Error: %s", dir, err.Error())
		}
	}
}

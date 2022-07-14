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

func openP4k(fileDir string) []*zip.File {
	r, err := zip.OpenReader(fileDir)
	if err != nil {
		println("unable to open p4k data file")
	}
	return r.File
}

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

	MakeDir(outputDir)
	files := openP4k(filepath.Join(gameDir, "Data.p4k"))

	div := runtime.NumCPU()
	fileCount := len(files)
	interval := fileCount / div

	filename := filepath.Join(outputDir, "P4k_filenames_")
	for i := 0; i < div; i++ {
		wg.Add(1)
		if i == div-1 {
			go buildFileNamesFile(fmt.Sprintf("%s%d.txt", filename, i+1), files, interval*i, fileCount)
		} else {
			go buildFileNamesFile(fmt.Sprintf("%s%d.txt", filename, i+1), files, interval*i, interval*(i+1))
		}
	}
	wg.Wait()
}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to create directory: %s - Error: %s", dir, err.Error())
		}
	}
}

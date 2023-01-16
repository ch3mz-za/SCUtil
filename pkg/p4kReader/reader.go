package p4k

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ch3mz-za/SCUtil/pkg/common"
)

var wg sync.WaitGroup

const dataP4k = "Data.p4k"

func GetP4kFilenames(gameDir, outputDir string) {

	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		fmt.Printf("Unable to open p4k data file: %s", err.Error())
	}
	defer r.Close()

	MakeDir(outputDir)

	fileCount := len(r.File)
	filename := filepath.Join(outputDir, "P4k_filenames.txt")

	p4kFileNameFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Unable to create p4k filenames file - Error: %s", err.Error())
		return
	}
	defer p4kFileNameFile.Close()

	var wg sync.WaitGroup
	progress := make(chan int)

	wg.Add(1)
	go common.ProgressBar(int64(fileCount), progress, &wg)

	cnt := 0
	for _, f := range r.File {
		p4kFileNameFile.WriteString(f.Name + "\n")
		cnt++
		if cnt%1000 == 0 {
			progress <- 1000
			cnt = 0
		}
	}

	progress <- cnt
	close(progress)
	wg.Wait()
}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to create directory: %s - Error: %s", dir, err.Error())
		}
	}
}

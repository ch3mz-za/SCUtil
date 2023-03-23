package p4k

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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

	filename := filepath.Join(outputDir, "P4k_filenames.txt")

	p4kFileNameFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Unable to create p4k filenames file - Error: %s", err.Error())
		return
	}
	defer p4kFileNameFile.Close()

	for _, f := range r.File {
		p4kFileNameFile.WriteString(f.Name + "\n")

	}
}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to create directory: %s - Error: %s", dir, err.Error())
		}
	}
}

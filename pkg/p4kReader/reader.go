package p4k

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

const dataP4k = "Data.p4k"

func GetP4kFilenames(gameDir, outputDir string) error {

	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		return fmt.Errorf("unable to open p4k data file:\n %s", err.Error())
	}
	defer r.Close()

	MakeDir(filepath.Dir(outputDir))

	p4kFileNameFile, err := os.Create(outputDir)
	if err != nil {
		return fmt.Errorf("unable to create p4k filenames file:\n %s", err.Error())
	}
	defer p4kFileNameFile.Close()

	for _, f := range r.File {
		p4kFileNameFile.WriteString(f.Name + "\n")
	}
	return nil
}

func MakeDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			fmt.Printf("Unable to create directory: %s - Error: %s", dir, err.Error())
		}
	}
}

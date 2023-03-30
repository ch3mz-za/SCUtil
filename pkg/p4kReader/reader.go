package p4k

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ch3mz-za/SCUtil/pkg/common"
)

const dataP4k = "Data.p4k"

func GetP4kFilenames(gameDir, outputPath string) error {

	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		return fmt.Errorf("unable to open p4k data file:\n %s", err.Error())
	}
	defer r.Close()

	common.MakeDir(outputPath)
	p4kFileNameFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to create p4k filenames file:\n %s", err.Error())
	}
	defer p4kFileNameFile.Close()

	for _, f := range r.File {
		p4kFileNameFile.WriteString(f.Name + "\n")
	}
	return nil
}

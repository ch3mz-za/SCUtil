package p4k

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ch3mz-za/SCUtil/pkg/common"
)

const dataP4k = "Data.p4k"

func GetP4kFilenames(gameDir, outputDir string) error {

	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		return fmt.Errorf("unable to open p4k data file:\n %s", err.Error())
	}
	defer r.Close()

	common.MakeDir(filepath.Dir(outputDir))

	p4kFileNameFile, err := os.Create(outputDir)
	if err != nil {
		return fmt.Errorf("unable to create p4k filenames file:\n %s", err.Error())
	}
	defer p4kFileNameFile.Close()

	for _, f := range r.File {
		if _, err = p4kFileNameFile.WriteString(f.Name + "\n"); err != nil {
			return err
		}
	}
	return nil
}

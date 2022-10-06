package main

import (
	"os"
	"path/filepath"

	"github.com/ch3mz-za/SCUtil/pkg/scu"
	log "github.com/sirupsen/logrus"
)

func main() {

	var err error

	scu.RootDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Unable to determine working directory")
	}
	scu.RootDir = filepath.Dir(scu.RootDir)

	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); !os.IsNotExist(err) {
			scu.RootDir = os.Args[1]
		}
	}

	m := scu.NewMenu()
	m.Run()
}

package p4k

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type SearchBorder struct {
	Start int
	Stop  int
}

func searchFilenameWorker(phrase string, r *zip.ReadCloser, border chan SearchBorder, results chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for b := range border {
		for i := b.Start; i < b.Stop; i++ {
			if strings.Contains(r.File[i].Name, phrase) {
				results <- r.File[i].Name
			}
		}
	}
}

func SearchP4kFilenames(gameDir, phrase, resultsFile string) error {

	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		return fmt.Errorf("unable to open p4k data file: %s", err.Error())
	}
	defer r.Close()

	results := make(chan string)
	borders := make(chan SearchBorder)
	resultsDone := make(chan bool, 1)
	defer close(resultsDone)

	go WriteStringsToFile(resultsFile, results, resultsDone)

	var wg sync.WaitGroup
	div := runtime.NumCPU()
	for i := 0; i < div; i++ {
		wg.Add(1)
		go searchFilenameWorker(phrase, r, borders, results, &wg)
	}

	fileCount := len(r.File)
	if fileCount > 1000 {
		interval := fileCount / div
		for i := 0; i < div; i++ {
			if i == div-1 {
				borders <- SearchBorder{Start: interval * i, Stop: fileCount}
			} else {
				borders <- SearchBorder{Start: interval * i, Stop: interval * (i + 1)}
			}
		}
	} else {
		borders <- SearchBorder{Start: 0, Stop: fileCount}
	}
	close(borders)

	go func() {
		wg.Wait()
		close(results)
	}()

	<-resultsDone

	return nil
}

func findInFile(phrase, filePath string, resultsChan chan string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fr := bufio.NewReader(file)
	for {
		line, _, err := fr.ReadLine()
		if err == io.EOF {
			return
		}

		if strings.Contains(string(line), phrase) {
			resultsChan <- string(line)
		}
	}
}

func WriteStringsToFile(filename string, results chan string, done chan bool) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for r := range results {
		if _, err = file.WriteString(r + "\n"); err != nil {
			continue
		}
	}

	done <- true
}

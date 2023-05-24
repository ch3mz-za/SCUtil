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
		for i := range r.File[b.Start:b.Stop] {
			if strings.Contains(r.File[i].Name, phrase) {
				results <- r.File[i].Name
			}
		}
	}
}

func SearchP4kFilenames(gameDir, phrase string) (*[]string, error) {
	r, err := zip.OpenReader(filepath.Join(gameDir, dataP4k))
	if err != nil {
		return nil, fmt.Errorf("unable to open p4k data file: %s", err.Error())
	}

	results := make(chan string)
	borders := make(chan SearchBorder)

	// TODO: Check if threads aren't creating copies of `r.File`
	//  [DONE - CHECK IF THIS WORKED]
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

	searchResults := make([]string, 0, 1000)
	for res := range results {
		searchResults = append(searchResults, res)
	}

	return &searchResults, nil
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

func WriteStringsToFile(filename string, strings *[]string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, s := range *strings {
		if _, err = file.WriteString(s + "\n"); err != nil {
			continue
		}
	}
}

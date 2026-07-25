package p4k

import (
	"archive/zip"
	"fmt"
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

func searchFilenameWorker(phrase string, r *zip.ReadCloser, border <-chan SearchBorder, results chan<- string, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for b := range border {
		for i := b.Start; i < b.Stop; i++ {
			if strings.Contains(r.File[i].Name, phrase) {
				select {
				case results <- r.File[i].Name:
				case <-stop:
					return
				}
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

	div := runtime.NumCPU()
	results := make(chan string)
	borders := make(chan SearchBorder, div)
	writerDone := make(chan struct{})
	stop := make(chan struct{})

	writerErr := make(chan error, 1)
	go WriteStringsToFile(resultsFile, results, writerDone, writerErr, stop)

	var wg sync.WaitGroup
	for range div {
		wg.Add(1)
		go searchFilenameWorker(phrase, r, borders, results, stop, &wg)
	}

	fileCount := len(r.File)
	if fileCount > 1000 {
		interval := fileCount / div
		for i := range div {
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

	<-writerDone
	select {
	case err := <-writerErr:
		return err
	default:
		return nil
	}

}

func WriteStringsToFile(filename string, results <-chan string, done chan<- struct{}, errCh chan<- error, stop chan<- struct{}) {
	defer close(done)
	file, err := os.Create(filename)
	if err != nil {
		errCh <- fmt.Errorf("create search results file: %w", err)
		close(stop)
		return
	}
	defer file.Close()

	for r := range results {
		if _, err = file.WriteString(r + "\n"); err != nil {
			errCh <- fmt.Errorf("write search results file: %w", err)
			close(stop)
			return
		}
	}
}

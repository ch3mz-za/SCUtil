package p4k

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FindInFiles(filenameDir, phrase string) ([]string, error) {
	files, err := ioutil.ReadDir(filenameDir)
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan string)
	for _, f := range files {
		wg.Add(1)
		go findInFile(phrase, filepath.Join(filenameDir, f.Name()), resultsChan)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var searchResults []string
	for res := range resultsChan {
		println(res)
		searchResults = append(searchResults, res)
	}

	return searchResults, nil
}

func findInFile(phrase, filePath string, resultsChan chan string) {
	defer wg.Done()

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

func WriteStringsToFile(filename string, strings []string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	for _, s := range strings {
		file.WriteString(s + "\n")
	}
}

package reader

import (
	"bufio"
	"customerimporter/pkg/email"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"sync"
)

type CSVReader struct {
	FilePath string
}

func (cr *CSVReader) ProcessSequentially(rowsChan chan<- []string) error {
	file, err := os.Open(cr.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	_, err = reader.Read()
	if err != nil {
		return err
	}
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		if len(row) > 0 {
			rowsChan <- row
		}
	}
	return nil
}

func (cr *CSVReader) ProcessConcurrently(domainCounter *email.DomainCounter) error {
	file, err := os.Open(cr.FilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	_, err = reader.Read() // Skipping header
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		wg.Add(1)
		go func(row []string) {
			defer wg.Done()
			domainCounter.CountDomain(row)
		}(row)
	}

	wg.Wait()
	return nil
}

func (cr *CSVReader) DoesFileExist() error {
	_, error := os.Stat(cr.FilePath)

	if os.IsNotExist(error) {
		return errors.New("file under path: " + cr.FilePath + " does not exists\n")
	} else {
		return nil
	}
}

func NewCSVReader(filePath string) *CSVReader {
	return &CSVReader{FilePath: filePath}
}

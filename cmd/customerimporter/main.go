package main

import (
	"customerimporter/pkg/email"
	"customerimporter/pkg/reader"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	SequentiallyMode string = "sequentially"
	ConcurrencyMode  string = "concurrency"
)

func main() {
	csvFilePath := flag.String("path", "", "Path to the CSV file")
	outputFilePath := flag.String("output", "", "Path to the output directory (optional; default is stdout)")
	mode := flag.String("mode", ConcurrencyMode, "Mode of processing: 'sequentially' or 'concurrency' (default: concurrency)")
	flag.Parse()

	if *csvFilePath == "" {
		log.Fatal("CSV file path is required. Use -csv-path to specify the file.")
	}
	csvReader := reader.NewCSVReader(*csvFilePath)
	err := csvReader.DoesFileExist()
	if err != nil {
		log.Fatal(err)
	}
	domainCounter := email.NewDomainCounter()

	switch *mode {
	case SequentiallyMode:
		rowsChan := make(chan []string)
		go domainCounter.ListenAndCount(rowsChan)
		if err := csvReader.ProcessSequentially(rowsChan); err != nil {
			log.Fatalf("Failed to process CSV file sequentially: %v", err)
		}
		close(rowsChan)
	case ConcurrencyMode:
		if err := csvReader.ProcessConcurrently(domainCounter); err != nil {
			log.Fatalf("Failed to process CSV file concurrently: %v", err)
		}
	default:
		if err := csvReader.ProcessConcurrently(domainCounter); err != nil {
			log.Fatalf("Failed to process CSV file concurrently: %v", err)
		}
	}

	output, err := email.PrepareOutput(domainCounter.GetCounts())
	if err != nil {
		log.Fatalf("Failed to prepare output: %v", err)
	}

	if *outputFilePath != "" {
		if err := ensureOutputDir(*outputFilePath); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}

		outputFile := filepath.Join(*outputFilePath, fmt.Sprintf("output-%v.csv", time.Now().Unix()))
		if err := os.WriteFile(outputFile, output, 0644); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
		log.Printf("File saved to %s", outputFile)
	} else {
		log.Default().Print(string(output))
	}
}

func ensureOutputDir(dir string) error {
	if dir == "" {
		return nil
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

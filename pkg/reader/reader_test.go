package reader

import (
	"customerimporter/pkg/email"
	"os"
	"testing"
)

func createTestCSVFile(content string) (string, error) {
	file, err := os.CreateTemp("", "test_*.csv")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func TestCSVReader_ProcessSequentially(t *testing.T) {
	content := "name,last_name,email\ntest,test,user1@example.com\ntest,test,user2@example.com\ntest,test,user3@example.it"
	filePath, err := createTestCSVFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filePath)

	csvReader := NewCSVReader(filePath)
	rowsChan := make(chan []string) // Using an unbuffered channel

	go func() {
		if err := csvReader.ProcessSequentially(rowsChan); err != nil {
			close(rowsChan)
			t.Error(err)
		}
		close(rowsChan)
	}()

	count := 0
	for range rowsChan {
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 rows, got %d", count)
	}
}

func TestCSVReader_ProcessConcurrently(t *testing.T) {
	content := "name,last_name,email\ntest,test,user1@example.com\ntest,test,user2@example.com\ntest,test,user3@example.it"
	filePath, err := createTestCSVFile(content)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filePath)

	csvReader := NewCSVReader(filePath)
	domainCounter := email.NewDomainCounter()

	if err := csvReader.ProcessConcurrently(domainCounter); err != nil {
		t.Error(err)
	}

	// Check domain counts
	counts := domainCounter.GetCounts()
	if len(counts) != 2 {
		t.Errorf("Expected 2 domains, got %d", len(counts))
	}
}

func TestCSVReader_DoesFileExist(t *testing.T) {
	filePath := "nonexistent_file.csv"
	csvReader := NewCSVReader(filePath)

	if err := csvReader.DoesFileExist(); err == nil {
		t.Errorf("Expected error for nonexistent file, got none")
	}
}

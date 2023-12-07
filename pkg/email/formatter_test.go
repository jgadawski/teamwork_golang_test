package email

import (
	"encoding/csv"
	"strings"
	"testing"
)

func TestPrepareOutput(t *testing.T) {
	domainCounts := []DomainCount{{Domain: "example.com", Count: 2}, {Domain: "test.com", Count: 1}}
	got, err := PrepareOutput(domainCounts)
	if err != nil {
		t.Errorf("PrepareOutput() error = %v", err)
		return
	}

	reader := csv.NewReader(strings.NewReader(string(got)))
	records, err := reader.ReadAll()
	if err != nil {
		t.Errorf("CSV Read error = %v", err)
	}

	if len(records) != 3 {
		t.Errorf("PrepareOutput() got %d records, want 3", len(records))
	}

	if records[0][0] != "domain" || records[0][1] != "count" {
		t.Errorf("PrepareOutput() CSV header incorrect, got = %v", records[0])
	}
}

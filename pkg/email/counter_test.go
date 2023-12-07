package email

import (
	"reflect"
	"sync"
	"testing"
)

func TestDomainCounter_CountDomain(t *testing.T) {
	tests := []struct {
		name string
		rows [][]string
		want []DomainCount
	}{
		{
			name: "count single domain",
			rows: [][]string{{"", "", "user@example.com"}},
			want: []DomainCount{{Domain: "example.com", Count: 1}},
		},
		{
			name: "count multiple domains",
			rows: [][]string{{"", "", "user1@example.com"}, {"", "", "user2@example.com"}, {"", "", "user@example.net"}},
			want: []DomainCount{{Domain: "example.com", Count: 2}, {Domain: "example.net", Count: 1}},
		},
		{
			name: "handle invalid email format",
			rows: [][]string{{"", "", "invalidemail"}, {"", "", "user@example.com"}},
			want: []DomainCount{{Domain: "example.com", Count: 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := NewDomainCounter()
			for _, row := range tt.rows {
				dc.CountDomain(row)
			}
			if got := dc.GetCounts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainCounter.CountDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainCounter_ConcurrentAccess(t *testing.T) {
	dc := NewDomainCounter()
	var wg sync.WaitGroup
	rows := [][]string{{"", "", "user1@example.com"}, {"", "", "user2@example.com"}, {"", "", "user3@example.net"}}

	for _, row := range rows {
		wg.Add(1)
		go func(row []string) {
			defer wg.Done()
			dc.CountDomain(row)
		}(row)
	}

	wg.Wait()

	// Check if the counts are as expected
	expected := []DomainCount{{Domain: "example.com", Count: 2}, {Domain: "example.net", Count: 1}}
	got := dc.GetCounts()

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Concurrent access to DomainCounter failed, got = %v, want %v", got, expected)
	}
}

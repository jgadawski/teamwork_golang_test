package email

import (
	"sort"
	"strings"
	"sync"
)

type DomainCount struct {
	Domain string
	Count  int
}

type DomainCounter struct {
	mu     sync.Mutex
	Counts []DomainCount
}

func NewDomainCounter() *DomainCounter {
	return &DomainCounter{}
}

func (dc *DomainCounter) CountDomain(row []string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	parts := strings.Split(row[2], "@")
	if len(parts) != 2 {
		return
	}
	domain := parts[1]

	// Find the position of the domain or where it should be inserted.
	i := sort.Search(len(dc.Counts), func(i int) bool {
		return dc.Counts[i].Domain >= domain
	})

	// Exists, increment the count.
	if i < len(dc.Counts) && dc.Counts[i].Domain == domain {
		dc.Counts[i].Count++
		return
	}

	// New domain, insert it into the slice at the sorted position.
	dc.Counts = append(dc.Counts, DomainCount{})
	copy(dc.Counts[i+1:], dc.Counts[i:])
	dc.Counts[i] = DomainCount{Domain: domain, Count: 1}
}

func (dc *DomainCounter) ListenAndCount(csvRows <-chan []string) {
	for csvRow := range csvRows {
		dc.CountDomain(csvRow)
	}
}

func (dc *DomainCounter) GetCounts() []DomainCount {
	return dc.Counts
}

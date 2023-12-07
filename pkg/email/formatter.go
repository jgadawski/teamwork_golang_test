package email

import (
	"encoding/csv"
	"fmt"
	"strings"
)

func PrepareOutput(domainCounts []DomainCount) ([]byte, error) {
	builder := &strings.Builder{}
	writer := csv.NewWriter(builder)
	header := []string{"domain", "count"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}
	for _, dc := range domainCounts {
		record := []string{dc.Domain, fmt.Sprintf("%d", dc.Count)}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return []byte(builder.String()), nil
}

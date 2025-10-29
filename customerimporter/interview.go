// Package customerimporter provides functionality for processing customer data from CSV files
// and aggregating statistics by email domain.
//
// The package reads customer records from CSV files (with format: first_name, last_name, email, gender, ip_address)
// and returns a slice of domain statistics sorted alphabetically by domain name.
//
// Performance characteristics:
//   - Time complexity: O(n) for reading and aggregating, O(d log d) for sorting where d is number of unique domains
//   - Space complexity: O(d) where d is number of unique domains
//   - Memory efficient: streams CSV processing, doesn't load entire file into memory
//
// Email validation rules:
//   - Must contain exactly one '@' symbol
//   - Local part (before @) must not be empty
//   - Domain part (after @) must not be empty
//   - Whitespace is trimmed from both email and domain
package customerimporter

import (
	"cmp"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
)

const (
	// CSV column index for email field
	emailColumnIndex = 2
)

// validateEmail validates email format and extracts the domain.
// Returns the domain and an error if the email is invalid.
// Valid email format: local-part@domain
// Domain must not be empty and must not contain whitespace.
func validateEmail(email string) (domain string, err error) {
	// Trim whitespace
	email = strings.TrimSpace(email)

	if email == "" {
		return "", fmt.Errorf("email address is empty")
	}

	// Split email into local and domain parts
	local, dom, found := strings.Cut(email, "@")
	if !found {
		return "", fmt.Errorf("invalid email format: missing '@' separator")
	}

	// Validate local part is not empty
	if strings.TrimSpace(local) == "" {
		return "", fmt.Errorf("invalid email format: empty local part")
	}

	// Validate domain is not empty
	dom = strings.TrimSpace(dom)
	if dom == "" {
		return "", fmt.Errorf("invalid email format: empty domain")
	}

	// Check for multiple @ symbols (strings.Cut only finds the first one)
	if strings.Contains(dom, "@") {
		return "", fmt.Errorf("invalid email format: multiple '@' symbols")
	}

	return dom, nil
}

// DomainData represents aggregated customer statistics for a single email domain.
type DomainData struct {
	// Domain is the email domain (e.g., "example.com")
	Domain string
	// CustomerQuantity is the number of customers with email addresses at this domain
	CustomerQuantity uint64
}

// CustomerImporter processes customer CSV files and aggregates domain statistics.
type CustomerImporter struct {
	path string
}

// NewCustomerImporter creates a new CustomerImporter that will read from the specified CSV file path.
//
// The filePath should point to a valid CSV file with customer data. The file is not opened or validated
// until ImportDomainData is called.
func NewCustomerImporter(filePath string) *CustomerImporter {
	return &CustomerImporter{
		path: filePath,
	}
}

// ImportDomainData reads customer data from the CSV file and returns aggregated domain statistics.
//
// The CSV file must have a header row and at least 3 columns, with the email address in the 3rd column (index 2).
// Expected CSV format:
//
//	first_name,last_name,email,gender,ip_address
//	John,Doe,john@example.com,Male,192.168.1.1
//
// Returns a slice of DomainData sorted alphabetically by domain name, or an error if:
//   - The file cannot be opened
//   - The CSV format is invalid (wrong number of columns)
//   - Any email address fails validation (see validateEmail)
//   - Any other CSV parsing error occurs
//
// The function processes the file incrementally and does not load the entire file into memory,
// making it suitable for processing large files.
//
// When verbose logging is enabled (via slog), progress is logged every 10,000 rows.
func (ci CustomerImporter) ImportDomainData() ([]DomainData, error) {
	file, err := os.Open(ci.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	data := make(map[string]uint64)

	// skip first line with headers
	_, readErr := csvReader.Read()
	if readErr != nil {
		slog.Error("failed to read CSV header", "error", readErr)
		return nil, readErr
	}

	rowCount := uint64(0)
	const progressInterval = 10000

	for line, readErr := csvReader.Read(); readErr != io.EOF; line, readErr = csvReader.Read() {
		if readErr != nil {
			return nil, readErr
		}
		rowCount++

		// Log progress every 10k rows
		if rowCount%progressInterval == 0 {
			slog.Info("processing", "rows", rowCount, "unique_domains", len(data))
		}

		// Validate CSV has enough columns
		if len(line) <= emailColumnIndex {
			return nil, fmt.Errorf("invalid CSV format: expected at least %d columns, got %d", emailColumnIndex+1, len(line))
		}

		// Validate email and extract domain
		domain, err := validateEmail(line[emailColumnIndex])
		if err != nil {
			return nil, fmt.Errorf("invalid email in CSV: %w", err)
		}

		data[domain] += 1
	}

	slog.Info("aggregation complete", "total_rows", rowCount, "unique_domains", len(data))
	domainData := make([]DomainData, 0, len(data))
	for k, v := range data {
		domainData = append(domainData, DomainData{
			Domain:           k,
			CustomerQuantity: v,
		})
	}
	slices.SortFunc(domainData, func(l, r DomainData) int {
		return cmp.Compare(l.Domain, r.Domain)
	})
	return domainData, nil
}

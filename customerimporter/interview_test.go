package customerimporter

import (
	"os"
	"strings"
	"testing"
)

func TestImportData(t *testing.T) {
	path := "./test_data.csv"
	importer := NewCustomerImporter(path)

	_, err := importer.ImportDomainData()
	if err != nil {
		t.Error(err)
	}
}

func TestImportDataSort(t *testing.T) {
	sortedDomains := []string{"360.cn", "acquirethisname.com", "blogtalkradio.com", "chicagotribune.com", "cnet.com", "cyberchimps.com", "github.io", "hubpages.com", "rediff.com", "statcounter.com"}
	path := "./test_data.csv"
	importer := NewCustomerImporter(path)
	data, err := importer.ImportDomainData()
	if err != nil {
		t.Error(err)
	}
	for i, v := range data {
		if v.Domain != sortedDomains[i] {
			t.Errorf("data not sorted properly. mismatch:\nhave: %v\nwant: %v", v.Domain, sortedDomains[i])
		}
	}
}

func TestImportInvalidPath(t *testing.T) {
	path := ""
	importer := NewCustomerImporter(path)

	_, err := importer.ImportDomainData()
	if err == nil {
		t.Error("invalid path error not caught")
	}
}

func TestImportInvalidData(t *testing.T) {
	path := "./test_invalid_data.csv"
	importer := NewCustomerImporter(path)

	_, err := importer.ImportDomainData()
	if err == nil {
		t.Error("invalid data not caught")
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		wantDomain  string
		expectError bool
	}{
		{
			name:        "valid email",
			email:       "user@example.com",
			wantDomain:  "example.com",
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			email:       "user@mail.example.com",
			wantDomain:  "mail.example.com",
			expectError: false,
		},
		{
			name:        "valid email with numbers",
			email:       "user123@example123.com",
			wantDomain:  "example123.com",
			expectError: false,
		},
		{
			name:        "email with whitespace - trimmed",
			email:       "  user@example.com  ",
			wantDomain:  "example.com",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			expectError: true,
		},
		{
			name:        "whitespace only email",
			email:       "   ",
			expectError: true,
		},
		{
			name:        "missing @ symbol",
			email:       "userexample.com",
			expectError: true,
		},
		{
			name:        "empty local part",
			email:       "@example.com",
			expectError: true,
		},
		{
			name:        "whitespace only local part",
			email:       "   @example.com",
			expectError: true,
		},
		{
			name:        "empty domain",
			email:       "user@",
			expectError: true,
		},
		{
			name:        "whitespace only domain",
			email:       "user@   ",
			expectError: true,
		},
		{
			name:        "multiple @ symbols",
			email:       "user@domain@extra.com",
			expectError: true,
		},
		{
			name:        "only @ symbol",
			email:       "@",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, err := validateEmail(tt.email)
			if tt.expectError {
				if err == nil {
					t.Errorf("validateEmail(%q) expected error, got nil", tt.email)
				}
			} else {
				if err != nil {
					t.Errorf("validateEmail(%q) unexpected error: %v", tt.email, err)
				}
				if domain != tt.wantDomain {
					t.Errorf("validateEmail(%q) = %q, want %q", tt.email, domain, tt.wantDomain)
				}
			}
		})
	}
}

func TestImportDomainData_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		csvContent  string
		expectError bool
		errorMsg    string
	}{
		{
			name: "insufficient columns",
			csvContent: `first_name,last_name,email,gender,ip_address
John,Doe`,
			expectError: true,
			errorMsg:    "wrong number of fields",
		},
		{
			name: "invalid email - no @",
			csvContent: `first_name,last_name,email,gender,ip_address
John,Doe,invalidemailexample.com,Male,192.168.1.1`,
			expectError: true,
			errorMsg:    "invalid email",
		},
		{
			name: "invalid email - empty domain",
			csvContent: `first_name,last_name,email,gender,ip_address
John,Doe,user@,Male,192.168.1.1`,
			expectError: true,
			errorMsg:    "invalid email",
		},
		{
			name: "invalid email - multiple @ symbols",
			csvContent: `first_name,last_name,email,gender,ip_address
John,Doe,user@domain@extra.com,Male,192.168.1.1`,
			expectError: true,
			errorMsg:    "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary CSV file
			tmpDir := t.TempDir()
			csvPath := tmpDir + "/test.csv"
			if err := writeTestCSV(csvPath, tt.csvContent); err != nil {
				t.Fatalf("failed to write test CSV: %v", err)
			}

			importer := NewCustomerImporter(csvPath)
			_, err := importer.ImportDomainData()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func BenchmarkImportDomainData(b *testing.B) {
	b.StopTimer()
	path := "./benchmark10k.csv"
	importer := NewCustomerImporter(path)

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := importer.ImportDomainData(); err != nil {
			b.Error(err)
		}
	}
}

// writeTestCSV is a helper function to write test CSV content to a file
func writeTestCSV(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

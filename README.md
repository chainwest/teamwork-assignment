# Customer Domain Importer

A high-performance Go CLI application that processes customer CSV files and aggregates email domain statistics.

## Features

- Streaming CSV processing for large files (millions of records)
- Comprehensive email validation
- Optional verbose logging mode
- Export to terminal or CSV file
- 67.5% test coverage

## Installation

```bash
# Build the binary
make build
# or
go build -o customer-importer main.go
```

**Requirements**: Go 1.21+

## Usage

```bash
# Process default file (./customers.csv) and print to stdout
./customer-importer

# Custom input file
./customer-importer -path=/path/to/customers.csv

# Export to file
./customer-importer -out=output.csv

# Enable verbose logging
./customer-importer -verbose

# All options combined
./customer-importer -path=input.csv -out=output.csv -verbose
```

### Flags

- `-path` - Input CSV file path (default: `./customers.csv`)
- `-out` - Output CSV file path (default: stdout)
- `-verbose` - Enable detailed logging (default: `false`)

### Input Format

```csv
first_name,last_name,email,gender,ip_address
John,Doe,john@example.com,Male,192.168.1.1
```

### Output Format

```csv
domain,number_of_customers
example.com,42
another.com,17
```

## Development

```bash
make help           # Show all commands
make build          # Build binary
make test           # Run tests
make test-coverage  # Tests with coverage
make benchmark      # Run benchmarks
make lint           # Run golangci-lint
make fmt            # Format code
make ci             # Run all CI checks
```

## Testing

```bash
# All tests
make test

# With coverage
make test-coverage

# Benchmarks
make benchmark
```

**Coverage**: 67.5% overall (92.5% customerimporter, 85.0% exporter)

## Architecture

```
CSV File → CustomerImporter → []DomainData → Terminal/File
              │
              ├─ Validate email
              ├─ Extract domain
              ├─ Aggregate counts (map)
              └─ Sort alphabetically
```

### Complexity

- Time: O(n) for processing, O(d log d) for sorting (d = unique domains)
- Space: O(d) where d is number of unique domains

## Project Structure

```
.
├── main.go                      # CLI entry point
├── customerimporter/            # CSV import and aggregation
├── exporter/                    # CSV export
├── .github/workflows/           # CI/CD
├── .golangci.yml               # Linter config
└── Makefile                    # Development tasks
```

# Verbose Logging Feature

Structured logging system using Go 1.21+ `slog` package for production observability.

## Two Modes

### Quiet Mode (Default)
- No logging output (only errors to stderr)
- Clean CSV output for piping
- Zero logging overhead

### Verbose Mode (`-verbose` flag)
- INFO level messages with progress tracking
- Structured logging with key-value pairs
- Logs to stderr, CSV to stdout

## Usage

```bash
# Quiet mode
./customer-importer -path=input.csv

# Verbose mode
./customer-importer -path=input.csv -verbose
```

## Example Output

### Quiet Mode
```
domain,number_of_customers
example.com,42
another.com,17
```

### Verbose Mode
```
time=2025-10-29T14:30:00.000+01:00 level=INFO msg="starting customer domain import" file=input.csv
time=2025-10-29T14:30:00.100+01:00 level=INFO msg="processing" rows=10000 unique_domains=1234
time=2025-10-29T14:30:00.200+01:00 level=INFO msg="processing" rows=20000 unique_domains=1500
time=2025-10-29T14:30:00.250+01:00 level=INFO msg="aggregation complete" total_rows=25000 unique_domains=1600
time=2025-10-29T14:30:00.255+01:00 level=INFO msg="import complete" domains=1600 duration=255ms
domain,number_of_customers
example.com,15000
another.com,10000
```

## Implementation

### Logger Setup
```go
func setupLogger(verbose bool) {
    if verbose {
        handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
        slog.SetDefault(slog.New(handler))
    } else {
        slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
    }
}
```

### Progress Logging
Progress updates every 10,000 rows:
```go
if rowCount%10000 == 0 {
    slog.Info("processing", "rows", rowCount, "unique_domains", len(data))
}
```

## Benefits

1. **Separation of Concerns**: Logs to stderr, data to stdout
2. **Structured Data**: Machine-parseable format
3. **Zero Overhead**: Quiet mode has no logging cost
4. **Progress Visibility**: Track processing of large files
5. **Production Ready**: Compatible with log aggregation systems

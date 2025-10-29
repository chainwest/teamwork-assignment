// Package main provides a CLI application for processing customer CSV files
// and generating email domain statistics.
//
// Usage:
//
//	# Process default file (./customers.csv) and print to stdout
//	go run main.go
//
//	# Process custom input file
//	go run main.go -path=/path/to/customers.csv
//
//	# Export results to a file instead of stdout
//	go run main.go -out=output.csv
//
//	# Custom input and output
//	go run main.go -path=input.csv -out=output.csv
//
//	# Enable verbose logging for detailed progress
//	go run main.go -verbose
//
// The application reads customer data from a CSV file, aggregates customers by email domain,
// and outputs the results either to stdout or to a CSV file.
//
// Flags:
//   - path: Input CSV file path (default: ./customers.csv)
//   - out: Output CSV file path (default: stdout)
//   - verbose: Enable detailed logging (default: false)
//
// Exit codes:
//   - 0: Success
//   - 1: Error occurred (file not found, invalid CSV, etc.)
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"importer/customerimporter"
	"importer/exporter"
	"log/slog"
)

// Options holds command-line flags for the application
type Options struct {
	path    *string
	outFile *string
	verbose *bool
}

func readOptions() *Options {
	opts := &Options{}
	opts.path = flag.String("path", "./customers.csv", "Path to the file with customer data")
	opts.outFile = flag.String("out", "", "Optional: output file path. If empty program will output results to the terminal")
	opts.verbose = flag.Bool("verbose", false, "Enable verbose logging with detailed progress information")
	flag.Parse()
	return opts
}

// setupLogger configures the global slog logger based on verbosity setting.
// In quiet mode (verbose=false), only ERROR level messages are shown.
// In verbose mode (verbose=true), INFO and DEBUG messages are also displayed.
func setupLogger(verbose bool) {
	var level slog.Level
	if verbose {
		level = slog.LevelInfo
	} else {
		// Only show errors in quiet mode
		level = slog.LevelError
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}

func main() {
	opts := readOptions()
	setupLogger(*opts.verbose)

	startTime := time.Now()
	slog.Info("starting customer domain import", "file", *opts.path)

	importer := customerimporter.NewCustomerImporter(*opts.path)
	data, err := importer.ImportDomainData()
	if err != nil {
		slog.Error("failed to import customer data", "error", err, "file", *opts.path)
		os.Exit(1)
	}

	duration := time.Since(startTime)
	slog.Info("import complete",
		"domains", len(data),
		"duration", duration.Round(time.Millisecond).String())

	if *opts.outFile == "" {
		printData(data)
	} else {
		exporter := exporter.NewCustomerExporter(*opts.outFile)
		if saveErr := exporter.ExportData(data); saveErr != nil {
			slog.Error("failed to export domain data", "error", saveErr, "file", *opts.outFile)
			os.Exit(1)
		}
		slog.Info("export complete", "file", *opts.outFile, "records", len(data))
	}
}

func printData(data []customerimporter.DomainData) {
	fmt.Println("domain,number_of_customers")
	for _, v := range data {
		fmt.Printf("%s,%v\n", v.Domain, v.CustomerQuantity)
	}
}

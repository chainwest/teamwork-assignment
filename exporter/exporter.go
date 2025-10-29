// Package exporter provides functionality for exporting customer domain statistics to CSV files.
//
// The package writes domain aggregation data (from customerimporter package) to CSV files
// with the format:
//
//	domain,number_of_customers
//	example.com,42
//	another.com,17
//
// The exporter creates or truncates the target file and writes data incrementally,
// making it suitable for large datasets.
package exporter

import (
	"encoding/csv"
	"fmt"
	"importer/customerimporter"
	"io"
	"log/slog"
	"os"
	"strconv"
)

// CustomerExporter exports customer domain statistics to CSV files.
type CustomerExporter struct {
	outputPath string
}

// NewCustomerExporter creates a new CustomerExporter that will write to the specified file path.
//
// The outputPath should be a valid file path. The file is created when ExportData is called.
// If the file already exists, it will be truncated (all existing content will be lost).
func NewCustomerExporter(outputPath string) *CustomerExporter {
	return &CustomerExporter{
		outputPath: outputPath,
	}
}

// ExportData writes customer domain statistics to a CSV file.
//
// The output CSV format is:
//
//	domain,number_of_customers
//	example.com,42
//	another.com,17
//
// The data parameter should be a slice of DomainData, typically from customerimporter.ImportDomainData.
// The data is written in the order provided (no sorting is performed by this function).
//
// WARNING: If the output file already exists, it will be truncated and all existing content will be lost.
//
// Returns an error if:
//   - data is nil
//   - the output file cannot be created (invalid path, permissions, etc.)
//   - any error occurs during CSV writing
//
// When verbose logging is enabled (via slog), export operations are logged.
func (ex CustomerExporter) ExportData(data []customerimporter.DomainData) error {
	if data == nil {
		return fmt.Errorf("provided data is empty (nil)")
	}

	slog.Info("starting export", "file", ex.outputPath, "records", len(data))

	outputFile, err := os.Create(ex.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	if err := exportCsv(data, outputFile); err != nil {
		return err
	}

	slog.Info("export written successfully", "file", ex.outputPath)
	return nil
}

func exportCsv(data []customerimporter.DomainData, output io.Writer) error {
	headers := []string{"domain", "number_of_customers"}
	csvWriter := csv.NewWriter(output)
	defer csvWriter.Flush()

	if err := csvWriter.Write(headers); err != nil {
		return err
	}
	for _, v := range data {
		pair := []string{v.Domain, strconv.FormatUint(v.CustomerQuantity, 10)}
		if err := csvWriter.Write(pair); err != nil {
			return err
		}
	}

	// Check for any errors that occurred during flush
	if err := csvWriter.Error(); err != nil {
		return err
	}
	return nil
}

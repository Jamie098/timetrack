package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func importFromCSV(filename string, config Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header and one data row")
	}

	// First row is header
	headers := records[0]
	if len(headers) < 2 {
		return fmt.Errorf("CSV must have at least Date and one project column")
	}

	// Projects are all columns except first (Date) and last (Total)
	projectCols := headers[1:]
	if strings.Contains(strings.ToLower(headers[len(headers)-1]), "total") {
		projectCols = headers[1 : len(headers)-1]
	}

	data := loadData()
	imported := 0

	// Process each data row
	for i, record := range records[1:] {
		if len(record) == 0 {
			continue
		}

		// Parse date
		dateStr := strings.TrimSpace(record[0])
		if dateStr == "" {
			continue
		}

		// Try to parse the date in various formats
		var parsedDate time.Time
		dateFormats := []string{
			"2-Jan",
			"2-Jan-06",
			"2006-01-02",
			"01/02/2006",
			"1/2/2006",
		}

		parsed := false
		for _, format := range dateFormats {
			if t, err := time.Parse(format, dateStr); err == nil {
				parsedDate = t
				// If year is not specified, use current year
				if format == "2-Jan" {
					now := time.Now()
					parsedDate = time.Date(now.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
				}
				parsed = true
				break
			}
		}

		if !parsed {
			fmt.Printf("Warning: Skipping row %d - invalid date format: %s\n", i+2, dateStr)
			continue
		}

		dayKey := parsedDate.Format("2006-01-02")

		// Create or load day data
		day := DayData{
			Date:             dayKey,
			Projects:         make(map[string]float64),
			ExcludedMeetings: make(map[string]float64),
		}

		if existing, ok := data[dayKey]; ok {
			day = existing
		}

		// Parse project percentages
		for j, projName := range projectCols {
			colIdx := j + 1
			if colIdx >= len(record) {
				break
			}

			valueStr := strings.TrimSpace(record[colIdx])
			if valueStr == "" {
				continue
			}

			// Remove % sign if present
			valueStr = strings.TrimSuffix(valueStr, "%")

			pct, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				fmt.Printf("Warning: Invalid percentage for %s on %s: %s\n", projName, dateStr, valueStr)
				continue
			}

			if day.Projects == nil {
				day.Projects = make(map[string]float64)
			}
			day.Projects[projName] = pct
		}

		data[dayKey] = day
		imported++
	}

	saveData(data)
	fmt.Printf("Successfully imported %d days from %s\n", imported, filename)
	return nil
}

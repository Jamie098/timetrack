package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

func exportToJSON(data map[string]DayData, filename string) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("Exported data to %s\n", filename)
	return nil
}

// getAllProjects extracts all unique projects from the data in alphabetical order
func getAllProjects(data map[string]DayData) []string {
	projectSet := make(map[string]bool)
	for _, day := range data {
		for project := range day.Projects {
			projectSet[project] = true
		}
	}

	projects := make([]string, 0, len(projectSet))
	for project := range projectSet {
		projects = append(projects, project)
	}
	sort.Strings(projects) // Alphabetical order
	return projects
}

// getWeekProjects gets all projects used in the current week
func getWeekProjects(data map[string]DayData) []string {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))

	projectSet := make(map[string]bool)
	for i := range 7 {
		date := monday.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		if day, exists := data[dateStr]; exists {
			for project := range day.Projects {
				projectSet[project] = true
			}
		}
	}

	projects := make([]string, 0, len(projectSet))
	for project := range projectSet {
		projects = append(projects, project)
	}
	sort.Strings(projects) // Alphabetical order
	return projects
}

func exportWeekToCSV(data map[string]DayData, config Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))

	// Auto-discover projects from this week's data (alphabetical)
	projects := getWeekProjects(data)

	// Header
	file.WriteString("Date")
	for _, p := range projects {
		file.WriteString("," + p)
	}
	file.WriteString(",Total Time Spent\n")

	// Each day
	for i := range 7 {
		date := monday.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		displayDate := date.Format("2-Jan")

		file.WriteString(displayDate)

		day, exists := data[dateStr]
		var total float64

		for _, p := range projects {
			pct := 0.0
			if exists {
				if val, ok := day.Projects[p]; ok {
					pct = val
				}
			}
			total += pct
			if pct == 0 {
				file.WriteString(",")
			} else {
				file.WriteString(fmt.Sprintf(",%.1f%%", pct))
			}
		}

		file.WriteString("\n")
	}

	fmt.Printf("Exported current week to %s\n", filename)
	return nil
}

func exportAllToCSV(data map[string]DayData, config Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Auto-discover all projects from all data (alphabetical)
	projects := getAllProjects(data)

	// Header
	file.WriteString("Date")
	for _, p := range projects {
		file.WriteString("," + p)
	}
	file.WriteString(",Total Time Spent\n")

	// Get all dates sorted
	dates := make([]string, 0, len(data))
	for date := range data {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	// Each day
	for _, dateStr := range dates {
		day := data[dateStr]
		t, _ := time.Parse("2006-01-02", dateStr)
		displayDate := t.Format("2-Jan")

		file.WriteString(displayDate)

		var total float64
		for _, p := range projects {
			pct := 0.0
			if val, ok := day.Projects[p]; ok {
				pct = val
			}
			total += pct
			if pct == 0 {
				file.WriteString(",")
			} else {
				file.WriteString(fmt.Sprintf(",%.1f%%", pct))
			}
		}

		file.WriteString("\n")
	}

	fmt.Printf("Exported all data to %s\n", filename)
	return nil
}

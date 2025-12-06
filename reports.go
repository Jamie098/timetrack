package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func generateWeeklyReport(data map[string]DayData) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	sunday := monday.AddDate(0, 0, 6)

	fmt.Println()
	fmt.Printf("%s📊 Weekly Report%s\n", ColorBold, ColorReset)
	fmt.Printf("%s to %s\n", monday.Format("Jan 2"), sunday.Format("Jan 2, 2006"))
	fmt.Println(strings.Repeat("─", 60))

	projectTotals := make(map[string]float64)
	var totalTracked float64
	var totalAvailable float64
	daysTracked := 0

	// Collect data for the week
	for i := range 7 {
		date := monday.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		if day, exists := data[dateStr]; exists && len(day.Projects) > 0 {
			daysTracked++
			available := getAvailablePercent(day)
			tracked := getTotalTracked(day)

			totalAvailable += available
			totalTracked += tracked

			for project, pct := range day.Projects {
				projectTotals[project] += pct
			}
		}
	}

	if daysTracked == 0 {
		fmt.Println("No data for this week")
		return
	}

	fmt.Printf("\n%sSummary:%s\n", ColorBold, ColorReset)
	fmt.Printf("  Days tracked: %d/7\n", daysTracked)
	fmt.Printf("  Total available: %.1f%%\n", totalAvailable)
	fmt.Printf("  Total tracked: %s%.1f%%%s\n", ColorBlue, totalTracked, ColorReset)
	fmt.Printf("  Average per day: %.1f%%\n", totalTracked/float64(daysTracked))

	if len(projectTotals) > 0 {
		fmt.Printf("\n%sTime by Project:%s\n", ColorBold, ColorReset)

		// Sort projects by time spent
		type projectTime struct {
			name  string
			total float64
		}
		projects := make([]projectTime, 0, len(projectTotals))
		for name, total := range projectTotals {
			projects = append(projects, projectTime{name, total})
		}
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].total > projects[j].total
		})

		for _, pt := range projects {
			percentage := (pt.total / totalTracked) * 100
			bar := progressBar(percentage, 15)
			fmt.Printf("  %s %s%.1f%%%s (%s%.0f%%%s of week) %s\n",
				bar, ColorBlue, pt.total, ColorReset,
				ColorCyan, percentage, ColorReset, pt.name)
		}
	}

	fmt.Println()
}

func generateProjectReport(data map[string]DayData, projectName string, config Config) {
	// Resolve project name
	projectName = resolveProject(projectName, config)

	fmt.Println()
	fmt.Printf("%s📊 Project Report: %s%s\n", ColorBold, projectName, ColorReset)
	fmt.Println(strings.Repeat("─", 60))

	// Get all dates sorted
	dates := make([]string, 0)
	for date := range data {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	type dayEntry struct {
		date string
		pct  float64
	}

	entries := make([]dayEntry, 0)
	var totalPct float64

	for _, dateStr := range dates {
		day := data[dateStr]
		if pct, ok := day.Projects[projectName]; ok {
			entries = append(entries, dayEntry{dateStr, pct})
			totalPct += pct
		}
	}

	if len(entries) == 0 {
		fmt.Printf("No time tracked for project '%s'\n", projectName)
		return
	}

	fmt.Printf("\n%sSummary:%s\n", ColorBold, ColorReset)
	fmt.Printf("  Days worked: %d\n", len(entries))
	fmt.Printf("  Total time: %.1f%%\n", totalPct)
	fmt.Printf("  Average per day: %.1f%%\n", totalPct/float64(len(entries)))

	// Show last 10 entries
	fmt.Printf("\n%sRecent Activity:%s\n", ColorBold, ColorReset)
	start := len(entries) - 10
	if start < 0 {
		start = 0
	}

	for i := len(entries) - 1; i >= start; i-- {
		entry := entries[i]
		t, _ := time.Parse("2006-01-02", entry.date)
		fmt.Printf("  %s: %s%.1f%%%s\n",
			t.Format("Jan 2, 2006"), ColorBlue, entry.pct, ColorReset)
	}

	fmt.Println()
}

func generateStatsReport(data map[string]DayData) {
	fmt.Println()
	fmt.Printf("%s📊 Statistics%s\n", ColorBold, ColorReset)
	fmt.Println(strings.Repeat("─", 60))

	if len(data) == 0 {
		fmt.Println("No data available")
		return
	}

	// Overall stats
	var totalTracked float64
	var totalAvailable float64
	overAllocatedDays := 0
	fullyAllocatedDays := 0
	projectFrequency := make(map[string]int)

	for _, day := range data {
		available := getAvailablePercent(day)
		tracked := getTotalTracked(day)

		totalAvailable += available
		totalTracked += tracked

		if tracked > available {
			overAllocatedDays++
		} else if tracked == available {
			fullyAllocatedDays++
		}

		for project := range day.Projects {
			projectFrequency[project]++
		}
	}

	avgTracked := totalTracked / float64(len(data))
	avgAvailable := totalAvailable / float64(len(data))

	fmt.Printf("\n%sOverall:%s\n", ColorBold, ColorReset)
	fmt.Printf("  Total days tracked: %d\n", len(data))
	fmt.Printf("  Average tracked: %.1f%%/day\n", avgTracked)
	fmt.Printf("  Average available: %.1f%%/day\n", avgAvailable)
	fmt.Printf("  Fully allocated days: %d\n", fullyAllocatedDays)
	fmt.Printf("  Over-allocated days: %s%d%s\n", ColorRed, overAllocatedDays, ColorReset)

	if len(projectFrequency) > 0 {
		fmt.Printf("\n%sMost Frequent Projects:%s\n", ColorBold, ColorReset)

		type projFreq struct {
			name  string
			count int
		}
		projects := make([]projFreq, 0, len(projectFrequency))
		for name, count := range projectFrequency {
			projects = append(projects, projFreq{name, count})
		}
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].count > projects[j].count
		})

		for i, pf := range projects {
			if i >= 10 {
				break
			}
			fmt.Printf("  %d. %s (%d days)\n", i+1, pf.name, pf.count)
		}
	}

	fmt.Println()
}

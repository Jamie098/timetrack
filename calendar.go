package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// printCalendar displays a calendar view of tracked time
func printCalendar(data map[string]DayData, days int) {
	// Get all dates and sort them in reverse chronological order
	dates := make([]string, 0)
	for date := range data {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	// Limit to requested number of days
	if len(dates) > days {
		dates = dates[:days]
	}

	if len(dates) == 0 {
		fmt.Println("No tracked time found")
		return
	}

	// Collect all unique projects across selected dates (for column headers)
	allProjects := make(map[string]bool)
	for _, date := range dates {
		day := data[date]
		for project := range day.Projects {
			allProjects[project] = true
		}
	}

	// Convert to sorted slice (alphabetical)
	projects := make([]string, 0, len(allProjects))
	for project := range allProjects {
		projects = append(projects, project)
	}
	sort.Strings(projects)

	// Calculate column widths (UK format "Mon 08/12/24" is 13 chars)
	dateWidth := 14

	// Print header
	fmt.Println()
	fmt.Printf("%-*s", dateWidth, "Date")
	for _, project := range projects {
		// Truncate long project names for column header
		displayName := project
		if len(displayName) > 15 {
			displayName = displayName[:12] + "..."
		}
		fmt.Printf("  %-15s", displayName)
	}
	fmt.Printf("  %-9s  %6s", "Tracked", "Avail")
	fmt.Println()

	// Print separator
	// dateWidth (14) + projects (17 each) + "Tracked" (13) + "Avail" (8)
	totalWidth := dateWidth + len(projects)*17 + 23
	fmt.Println(strings.Repeat("─", totalWidth))

	// Print each day
	for _, date := range dates {
		day := data[date]
		available := getAvailablePercent(day)
		tracked := getTotalTracked(day)

		// Format date (e.g., "Mon 08/12/24" - DD/MM/YY UK format)
		parsedDate, err := time.Parse("2006-01-02", date)
		var dateStr string
		if err == nil {
			dateStr = parsedDate.Format("Mon 02/01/06")
		} else {
			dateStr = date
		}

		// Print date column
		fmt.Printf("%-*s", dateWidth, dateStr)

		// Print each project percentage
		for _, project := range projects {
			pct, exists := day.Projects[project]
			if exists && pct > 0 {
				// Format: "  " + color + "95.0%" (6 chars) + reset + 9 spaces = 15 visible chars total
				fmt.Printf("  %s%-6s%s         ", ColorBlue, fmt.Sprintf("%.1f%%", pct), ColorReset)
			} else {
				fmt.Printf("  %-15s", "-")
			}
		}

		// Print total with color coding based on remaining time
		remaining := available - tracked
		totalStr := fmt.Sprintf("%.1f%%", tracked)
		var colorCode string
		var statusIcon string

		if tracked > available {
			colorCode = ColorRed // Over-allocated
			statusIcon = " ⚠️"
		} else if remaining < 10 {
			colorCode = ColorYellow // Nearly full (less than 10% remaining)
			statusIcon = ""
		} else {
			colorCode = ColorGreen // Good (10%+ remaining)
			statusIcon = ""
		}

		fmt.Printf("  %s%-6s%s%-3s", colorCode, totalStr, ColorReset, statusIcon)

		// Print available percentage aligned to the right (6 chars to fit "100.0%")
		availStr := fmt.Sprintf("%.1f%%", available)
		if day.ExcludedPercent > 0 {
			fmt.Printf("  %s%6s%s", ColorGray, availStr, ColorReset)
		} else {
			fmt.Printf("  %6s", availStr)
		}

		fmt.Println()
	}

	fmt.Println()

	// Print legend
	fmt.Printf("%sLegend:%s ", ColorGray, ColorReset)
	fmt.Printf("%s●%s On track  ", ColorGreen, ColorReset)
	fmt.Printf("%s●%s Nearly full (<10%% left)  ", ColorYellow, ColorReset)
	fmt.Printf("%s●%s Over-allocated\n", ColorRed, ColorReset)
	fmt.Printf("  Tracked = time logged  |  Avail = available after excluding meetings ")
	fmt.Printf("(%sgray%s = has exclusions)\n", ColorGray, ColorReset)
	fmt.Println()
}

// printCompactCalendar shows a more condensed calendar view
func printCompactCalendar(data map[string]DayData, days int) {
	// Get all dates and sort them in reverse chronological order
	dates := make([]string, 0)
	for date := range data {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	// Limit to requested number of days
	if len(dates) > days {
		dates = dates[:days]
	}

	if len(dates) == 0 {
		fmt.Println("No tracked time found")
		return
	}

	fmt.Println()
	fmt.Println("📆 Time Tracking Calendar")
	fmt.Println(strings.Repeat("─", 80))

	for _, date := range dates {
		day := data[date]
		available := getAvailablePercent(day)
		tracked := getTotalTracked(day)
		remaining := available - tracked

		// Format date
		parsedDate, err := time.Parse("2006-01-02", date)
		var dateStr string
		if err == nil {
			dateStr = parsedDate.Format("Mon 01/02/06")
		} else {
			dateStr = date
		}

		// Status icon
		statusIcon := "⏳"
		if tracked > available {
			statusIcon = "⚠️"
		} else if remaining < 1 {
			statusIcon = "✓"
		}

		fmt.Printf("\n%s %s  ", statusIcon, dateStr)

		// Progress bar for the day
		barWidth := 30
		filledWidth := int((tracked / available) * float64(barWidth))
		if filledWidth > barWidth {
			filledWidth = barWidth
		}

		bar := "["
		for i := 0; i < barWidth; i++ {
			if i < filledWidth {
				if tracked > available {
					bar += ColorRed + "█" + ColorReset
				} else if remaining < 10 {
					bar += ColorYellow + "█" + ColorReset
				} else {
					bar += ColorGreen + "█" + ColorReset
				}
			} else {
				bar += ColorGray + "░" + ColorReset
			}
		}
		bar += "]"
		fmt.Printf("%s  %.1f%%", bar, tracked)

		// Projects list
		if len(day.Projects) > 0 {
			fmt.Println()
			projects := sortedKeys(day.Projects)
			for _, name := range projects {
				pct := day.Projects[name]
				fmt.Printf("    • %s%-5.1f%%%s  %s\n", ColorBlue, pct, ColorReset, name)
			}
		} else {
			fmt.Println("  (no projects)")
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("─", 80))
	fmt.Println()
}

package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func printStatus(day DayData) {
	available := getAvailablePercent(day)
	tracked := getTotalTracked(day)
	remaining := available - tracked

	fmt.Println()
	fmt.Printf("%s📅 %s%s\n", ColorBold, day.Date, ColorReset)
	fmt.Println(strings.Repeat("─", 45))

	if day.ExcludedPercent > 0 {
		fmt.Printf("🚫 Excluded (ceremonies): %s%.1f%%%s\n", ColorGray, day.ExcludedPercent, ColorReset)
		if len(day.ExcludedMeetings) > 0 {
			meetings := sortedKeys(day.ExcludedMeetings)
			for _, name := range meetings {
				fmt.Printf("   • %s: %s%.1f%%%s\n", name, ColorGray, day.ExcludedMeetings[name], ColorReset)
			}
		}
		fmt.Println()
	}

	fmt.Printf("📊 Available to track: %s%.1f%%%s\n", ColorCyan, available, ColorReset)
	fmt.Printf("✅ Tracked: %s%.1f%%%s\n", ColorBlue, tracked, ColorReset)

	statusColor := getStatusColor(remaining)
	fmt.Printf("⏳ Remaining: %s%.1f%%%s\n", statusColor, remaining, ColorReset)
	fmt.Println()

	if len(day.Projects) > 0 {
		fmt.Println("Projects:")
		projects := sortedKeys(day.Projects)
		for _, name := range projects {
			pct := day.Projects[name]
			bar := progressBar(pct, 20)
			fmt.Printf("   %s %s%5.1f%%%s %s\n", bar, ColorBlue, pct, ColorReset, name)
		}
		fmt.Println()
	}

	if remaining < 0 {
		fmt.Printf("%s⚠️  Over-allocated by %.1f%%!%s\n\n", ColorRed, -remaining, ColorReset)
	} else if remaining == 0 {
		fmt.Printf("%s✨ Day fully allocated!%s\n", ColorGreen, ColorReset)
	} else if remaining < 10 {
		fmt.Printf("%s💡 Only %.1f%% remaining - almost done!%s\n", ColorYellow, remaining, ColorReset)
	}
}

func progressBar(pct float64, width int) string {
	filled := int(pct / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}

func sortedKeys(m map[string]float64) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func printWeekExport(data map[string]DayData, config Config) {
	if len(config.Projects) == 0 {
		fmt.Println("No projects configured. Set them first:")
		fmt.Println("  timetrack projects set \"Proj1,Proj2,Proj3,...\"")
		return
	}

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))

	// Header
	fmt.Print("Date")
	for _, p := range config.Projects {
		fmt.Print("," + p)
	}
	fmt.Println(",Total Time Spent")

	// Each day
	for i := range 7 {
		date := monday.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		displayDate := date.Format("2-Jan")

		fmt.Print(displayDate)

		day, exists := data[dateStr]
		var total float64

		for _, p := range config.Projects {
			pct := 0.0
			if exists {
				if val, ok := day.Projects[p]; ok {
					pct = val
				}
			}
			total += pct
			if pct == 0 {
				fmt.Print(",")
			} else {
				fmt.Printf(",%.1f%%", pct)
			}
		}

		fmt.Println()
	}
}

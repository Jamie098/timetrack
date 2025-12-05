package main

import (
	"fmt"
	"sort"
	"strings"
)

func printStatus(day DayData) {
	available := getAvailablePercent(day)
	tracked := getTotalTracked(day)
	remaining := available - tracked

	fmt.Println()
	fmt.Printf("📅 %s\n", day.Date)
	fmt.Println(strings.Repeat("─", 45))

	if day.ExcludedPercent > 0 {
		fmt.Printf("🚫 Excluded (ceremonies): %.1f%%\n", day.ExcludedPercent)
		if len(day.ExcludedMeetings) > 0 {
			meetings := sortedKeys(day.ExcludedMeetings)
			for _, name := range meetings {
				fmt.Printf("   • %s: %.1f%%\n", name, day.ExcludedMeetings[name])
			}
		}
		fmt.Println()
	}

	fmt.Printf("📊 Available to track: %.1f%%\n", available)
	fmt.Printf("✅ Tracked: %.1f%%\n", tracked)
	fmt.Printf("⏳ Remaining: %.1f%%\n", remaining)
	fmt.Println()

	if len(day.Projects) > 0 {
		fmt.Println("Projects:")
		projects := sortedKeys(day.Projects)
		for _, name := range projects {
			pct := day.Projects[name]
			bar := progressBar(pct, 20)
			fmt.Printf("   %s %5.1f%% %s\n", bar, pct, name)
		}
		fmt.Println()
	}

	if remaining < 0 {
		fmt.Printf("⚠️  Over-allocated by %.1f%%!\n\n", -remaining)
	} else if remaining == 0 {
		fmt.Println("✨ Day fully allocated!\n")
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

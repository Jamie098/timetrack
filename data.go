package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func loadData() map[string]DayData {
	data := make(map[string]DayData)
	bytes, err := os.ReadFile(getDataPath())
	if err != nil {
		return data
	}
	json.Unmarshal(bytes, &data)
	return data
}

func saveData(data map[string]DayData) {
	bytes, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(getDataPath(), bytes, 0644)
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func parseDate(dateStr string) (string, error) {
	// If empty, return today
	if dateStr == "" {
		return today(), nil
	}

	// Try various date formats
	formats := []string{
		"2006-01-02",     // YYYY-MM-DD (ISO format)
		"02-01-2006",     // DD-MM-YYYY (UK format) - try this BEFORE US format
		"2006/01/02",     // YYYY/MM/DD
		"02/01/2006",     // DD/MM/YYYY (UK format)
		"01/02/2006",     // MM/DD/YYYY (US format)
		"01-02-2006",     // MM-DD-YYYY (US format)
		"2-Jan-06",       // Excel format
		"2-Jan",          // Short format (current year)
	}

	for _, format := range formats {
		t, err := time.Parse(format, dateStr)
		if err == nil {
			// If year is not specified (format without year), use current year
			if format == "2-Jan" {
				year := time.Now().Year()
				t = time.Date(year, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			}
			return t.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("invalid date format: %s (try DD-MM-YYYY or YYYY-MM-DD)", dateStr)
}

func getTargetDate(args []string, flagName string) (string, []string, error) {
	// Look for --date or -d flag
	targetDate := today()
	remainingArgs := []string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--date" || arg == "-d" {
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("--date flag requires a value")
			}
			var err error
			targetDate, err = parseDate(args[i+1])
			if err != nil {
				return "", nil, err
			}
			i++ // Skip the date value
		} else {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	return targetDate, remainingArgs, nil
}

func todayWeekday() string {
	return strings.ToLower(time.Now().Weekday().String()[:3])
}

func isWeekday() bool {
	day := time.Now().Weekday()
	return day >= time.Monday && day <= time.Friday
}

func getTodayData(data map[string]DayData, config Config) DayData {
	return getDateData(data, config, today())
}

func getDateData(data map[string]DayData, config Config, date string) DayData {
	if d, ok := data[date]; ok {
		return d
	}

	// New day - apply recurring meetings
	day := DayData{
		Date:             date,
		ExcludedPercent:  0,
		Projects:         make(map[string]float64),
		ExcludedMeetings: make(map[string]float64),
	}

	// Parse the date to get weekday
	t, err := time.Parse("2006-01-02", date)
	if err == nil {
		weekday := strings.ToLower(t.Weekday().String()[:3])
		isWeekdayDate := t.Weekday() >= time.Monday && t.Weekday() <= time.Friday

		for _, meeting := range config.RecurringMeetings {
			if shouldApplyMeetingForDate(meeting, weekday, isWeekdayDate) {
				day.ExcludedMeetings[meeting.Name] = meeting.Percent
				day.ExcludedPercent += meeting.Percent
			}
		}
	}

	return day
}

func shouldApplyMeetingForDate(meeting RecurringMeeting, weekday string, isWeekday bool) bool {
	for _, d := range meeting.Days {
		d = strings.ToLower(d)
		if d == weekday {
			return true
		}
		if d == "daily" {
			return true
		}
		if d == "weekdays" && isWeekday {
			return true
		}
	}
	return false
}

func shouldApplyMeeting(meeting RecurringMeeting, weekday string) bool {
	for _, d := range meeting.Days {
		d = strings.ToLower(d)
		if d == weekday {
			return true
		}
		if d == "daily" {
			return true
		}
		if d == "weekdays" && isWeekday() {
			return true
		}
	}
	return false
}

func getAvailablePercent(day DayData) float64 {
	return 100.0 - day.ExcludedPercent
}

func getTotalTracked(day DayData) float64 {
	var total float64
	for _, pct := range day.Projects {
		total += pct
	}
	return total
}

func resolveProject(name string, config Config) string {
	if fullName, ok := config.Aliases[strings.ToLower(name)]; ok {
		return fullName
	}
	return name
}

func hoursToPercent(hours float64) float64 {
	// Assuming an 8-hour workday
	return (hours / 8.0) * 100.0
}

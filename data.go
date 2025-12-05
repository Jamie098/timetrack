package main

import (
	"encoding/json"
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

func todayWeekday() string {
	return strings.ToLower(time.Now().Weekday().String()[:3])
}

func isWeekday() bool {
	day := time.Now().Weekday()
	return day >= time.Monday && day <= time.Friday
}

func getTodayData(data map[string]DayData, config Config) DayData {
	if d, ok := data[today()]; ok {
		return d
	}

	// New day - apply recurring meetings
	day := DayData{
		Date:             today(),
		ExcludedPercent:  0,
		Projects:         make(map[string]float64),
		ExcludedMeetings: make(map[string]float64),
	}

	weekday := todayWeekday()
	for _, meeting := range config.RecurringMeetings {
		if shouldApplyMeeting(meeting, weekday) {
			day.ExcludedMeetings[meeting.Name] = meeting.Percent
			day.ExcludedPercent += meeting.Percent
		}
	}

	return day
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

package main

import (
	"encoding/json"
	"os"
)

func loadConfig() Config {
	config := Config{
		ReminderTimes:     []string{"09:00", "12:00", "15:00"},
		RecurringMeetings: []RecurringMeeting{},
		Projects:          []string{},
	}
	bytes, err := os.ReadFile(getConfigPath())
	if err != nil {
		// Create default config
		saveConfig(config)
		return config
	}
	json.Unmarshal(bytes, &config)
	return config
}

func saveConfig(config Config) {
	bytes, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(getConfigPath(), bytes, 0644)
}

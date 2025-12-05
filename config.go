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
		Aliases:           make(map[string]string),
	}
	bytes, err := os.ReadFile(getConfigPath())
	if err != nil {
		saveConfig(config)
		return config
	}
	json.Unmarshal(bytes, &config)
	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}
	return config
}

func saveConfig(config Config) {
	bytes, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(getConfigPath(), bytes, 0644)
}

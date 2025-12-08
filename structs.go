package main

type DayData struct {
	Date             string             `json:"date"`
	ExcludedPercent  float64            `json:"excluded_percent"`
	Projects         map[string]float64 `json:"projects"`
	ExcludedMeetings map[string]float64 `json:"excluded_meetings"`
	LastModified     string             `json:"last_modified,omitempty"` // Track last added/edited project for undo
}

type RecurringMeeting struct {
	Name    string   `json:"name"`
	Percent float64  `json:"percent"`
	Days    []string `json:"days"` // "mon", "tue", "wed", "thu", "fri", "sat", "sun", or "daily", "weekdays"
}

type Config struct {
	ReminderTimes     []string           `json:"reminder_times"`
	RecurringMeetings []RecurringMeeting `json:"recurring_meetings"`
	Projects          []string           `json:"projects"`
	Aliases           map[string]string  `json:"aliases"`
	TimesheetURL      string             `json:"timesheet_url,omitempty"`
}

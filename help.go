package main

import (
	"fmt"
	"sort"
	"strings"
)

func printHelp() {
	fmt.Println(`
timetrack - Track your day in hours

Usage:
  timetrack                        Show today's status
  timetrack interactive            Interactive menu mode
  timetrack add <project> <hours>  Add/update time to a project
  timetrack fill <project>         Fill remaining time with project
  timetrack edit <project> <hours> Update existing project time
  timetrack exclude <name> <hours> Exclude ceremony time (one-off)
  timetrack rm <project>           Remove a project entry
  timetrack rmex <name>            Remove an excluded meeting
  timetrack undo                   Remove last added project
  timetrack clear                  Clear today's data

Viewing:
  timetrack                        Show today's status
  timetrack show [days]            Calendar view (default: 7 days)

Date Flag (for add, fill, edit, rm):
  --date <date> or -d <date>   Work with a specific date

  Supported formats:
    DD-MM-YYYY (UK format, preferred): 05-12-2024
    YYYY-MM-DD (ISO format):           2024-12-05
    DD/MM/YYYY (UK format):            05/12/2024
    MM/DD/YYYY (US format):            12/05/2024

  Examples:
    timetrack add bugs 2 --date 05-12-2024
    timetrack fill "Main Project" -d 2024-12-05
    timetrack edit automation 3 --date 05/12/2024
    timetrack rm bugs -d 05-12-2024

Export & Import:
  timetrack export csv [file]      Export week to CSV (auto-discovered projects)
  timetrack export all [file]      Export all data to CSV
  timetrack export json [file]     Export as JSON
  timetrack import <csv-file>      Import data from CSV

Projects & Aliases:
  timetrack projects list          List all projects (auto-discovered from time)
  timetrack alias <short> <full>   Create/update alias
  timetrack alias rm <short>       Remove alias
  timetrack alias list             List all aliases

Config:
  timetrack config                 Show current config
  timetrack config edit            Open config file in editor
  timetrack meeting add <name> <hours> <days>   Add recurring meeting
  timetrack meeting rm <name>      Remove recurring meeting
  timetrack reminder <times>       Set reminder times (e.g., "09:00,12:00,15:00")
  timetrack url set <url>          Set online timesheet URL
  timetrack url open               Open timesheet URL in browser
  timetrack url                    Show current timesheet URL
  timetrack url rm                 Clear timesheet URL

Reminders:
  timetrack start                  Start reminder service (foreground)
  timetrack start-bg               Start reminder service (background)
  timetrack stop                   Stop reminder service
  timetrack status                 Check if reminder service is running

Days: mon, tue, wed, thu, fri, sat, sun, daily, weekdays

Note: Based on 8-hour workday. All input is in hours, converted to percentages internally.`)
}

func printConfig(config Config) {
	fmt.Println()
	fmt.Println("⚙️  Configuration")
	fmt.Println(strings.Repeat("─", 45))
	fmt.Printf("Config file: %s\n\n", getConfigPath())

	fmt.Println("Reminder times:")
	if len(config.ReminderTimes) == 0 {
		fmt.Println("   (none)")
	} else {
		for _, t := range config.ReminderTimes {
			fmt.Printf("   • %s\n", t)
		}
	}

	fmt.Println("\nRecurring meetings:")
	if len(config.RecurringMeetings) == 0 {
		fmt.Println("   (none)")
	} else {
		for _, m := range config.RecurringMeetings {
			fmt.Printf("   • %s: %.1f%% on %s\n", m.Name, m.Percent, strings.Join(m.Days, ", "))
		}
	}

	fmt.Println("\nProjects:")
	fmt.Println("   Auto-discovered from your time entries")
	fmt.Println("   Use 'timetrack projects list' to see all tracked projects")

	fmt.Println("\nAliases:")
	if len(config.Aliases) == 0 {
		fmt.Println("   (none)")
	} else {
		aliases := make([]string, 0, len(config.Aliases))
		for k := range config.Aliases {
			aliases = append(aliases, k)
		}
		sort.Strings(aliases)
		for _, k := range aliases {
			fmt.Printf("   • %s → %s\n", k, config.Aliases[k])
		}
	}

	fmt.Println("\nTimesheet URL:")
	if config.TimesheetURL == "" {
		fmt.Println("   (not set)")
	} else {
		fmt.Printf("   %s\n", config.TimesheetURL)
	}
	fmt.Println()
}

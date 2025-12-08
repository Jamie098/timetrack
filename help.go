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
  timetrack                        Interactive mode (no args)
  timetrack add <project> <hours>  Add/update time to a project
  timetrack fill <project>         Fill remaining time with project
  timetrack edit <project> <hours> Update existing project time
  timetrack exclude <name> <hours> Exclude ceremony time (one-off)
  timetrack rm <project>           Remove a project entry
  timetrack rmex <name>            Remove an excluded meeting
  timetrack undo                   Remove last added project
  timetrack clear                  Clear today's data
  timetrack summary                Compact one-line status
  timetrack history [days]         Show history (default: 7 days)

Export & Import:
  timetrack week                   Output current week as CSV (stdout)
  timetrack export [format] [file] Export data (formats: csv, json, all)
  timetrack import <csv-file>      Import data from CSV

Reports & Analytics:
  timetrack report week            Weekly summary report
  timetrack report project <name>  Project-specific report
  timetrack report stats           Overall statistics

Projects & Aliases:
  timetrack projects parse "Date,P1,P2,...,Total"   Parse Excel header (auto-aliases)
  timetrack projects set "P1,P2,P3"                 Set project columns manually
  timetrack projects list                           List project columns
  timetrack alias <short> <full>                    Create/update alias
  timetrack alias rm <short>                        Remove alias
  timetrack alias list                              List all aliases

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

Features:
  • Interactive mode when no command specified
  • Fuzzy matching for project names
  • Color-coded display (green/yellow/red based on allocation)
  • Automatic warnings for over-allocation or unusual values
  • CSV import/export with multiple formats
  • Comprehensive reports and analytics

Quick Start:
  1. Parse your Excel header:
     timetrack projects parse "Date,CT.GOV Automation,Bugs,...,Total"

  2. Add recurring meetings:
     timetrack meeting add standup 0.5 weekdays

  3. Track time (supports fuzzy matching):
     timetrack add ctgo 2
     timetrack add bugs 1.5

  4. View reports:
     timetrack report week

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

	fmt.Println("\nProjects (column order for export):")
	if len(config.Projects) == 0 {
		fmt.Println("   (none - use 'timetrack projects parse \"...\"')")
	} else {
		for i, p := range config.Projects {
			fmt.Printf("   %d. %s\n", i+1, p)
		}
	}

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

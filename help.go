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
  timetrack add <project> <hours>  Add/update time to a project (use alias or full name)
  timetrack exclude <name> <hours> Exclude ceremony time (one-off)
  timetrack rm <project>           Remove a project entry
  timetrack rmex <name>            Remove an excluded meeting
  timetrack clear                  Clear today's data
  timetrack history [days]         Show history (default: 7 days)

Export:
  timetrack week                   Output current week as CSV (for Excel)

Projects & Aliases:
  timetrack projects parse "Date,P1,P2,...,Total"   Parse Excel header (auto-generates aliases)
  timetrack projects set "P1,P2,P3"                 Set project columns manually
  timetrack projects list                           List project columns
  timetrack alias <short> <full>                    Create/update alias
  timetrack alias rm <short>                        Remove alias
  timetrack alias list                              List all aliases

Config:
  timetrack config               Show current config
  timetrack config edit          Open config file in editor
  timetrack meeting add <name> <hours> <days>   Add recurring meeting
  timetrack meeting rm <name>    Remove recurring meeting
  timetrack reminder <times>     Set reminder times (e.g., "09:00,12:00,15:00")

Reminders:
  timetrack start                Start reminder service (foreground)
  timetrack start-bg             Start reminder service (background)
  timetrack stop                 Stop reminder service
  timetrack status               Check if reminder service is running

Setup:
  1. Parse your Excel header (copy the header row from CSV):
     timetrack projects parse "Date,CT.GOV Automation,Bugs,...,Total Time Spent"

  2. Check your aliases:
     timetrack alias list

  3. Add recurring meetings:
     timetrack meeting add standup 0.5 weekdays

  4. Start tracking:
     timetrack add ctgo 2
     timetrack add bugs 1

  5. Export at end of week:
     timetrack week

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
	fmt.Println()
}

package main

import (
	"fmt"
	"sort"
	"strings"
)

func printHelp() {
	fmt.Println(`
timetrack - Track your day as percentages

Usage:
  timetrack                      Show today's status
  timetrack add <project> <%>    Add/update time to a project
  timetrack exclude <name> <%>   Exclude ceremony time (one-off)
  timetrack rm <project>         Remove a project entry
  timetrack rmex <name>          Remove an excluded meeting
  timetrack clear                Clear today's data
  timetrack history [days]       Show history (default: 7 days)

Export:
  timetrack week                 Output current week as tab-separated (for Excel)

Projects & Aliases:
  timetrack projects set "P1,P2,P3"   Set project columns (for week export)
  timetrack projects list             List project columns
  timetrack alias <short> <full>      Create alias (e.g., ctgov -> CT.Gov Automation)
  timetrack alias rm <short>          Remove alias
  timetrack alias list                List all aliases

Config:
  timetrack config               Show current config
  timetrack config edit          Open config file in editor
  timetrack meeting add <name> <percent> <days>   Add recurring meeting
  timetrack meeting rm <name>    Remove recurring meeting
  timetrack reminder <times>     Set reminder times (e.g., "09:00,12:00,15:00")

Reminders:
  timetrack start                Start reminder service (foreground)
  timetrack start-bg             Start reminder service (background)
  timetrack stop                 Stop reminder service
  timetrack status               Check if reminder service is running

Examples:
  timetrack projects set "CT.Gov Automation,Bugs,Tech Debt"
  timetrack alias ctgov CT.Gov Automation
  timetrack add ctgov 25
  timetrack week

Days: mon, tue, wed, thu, fri, sat, sun, daily, weekdays

Quick reference (8hr day):
  15min = 3.125%    30min = 6.25%    45min = 9.375%
  1hr   = 12.5%     1.5hr = 18.75%   2hr   = 25%`)
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
		fmt.Println("   (none - use 'timetrack projects set \"Proj1,Proj2,...\"')")
	} else {
		for i, p := range config.Projects {
			fmt.Printf("   %d. %s\n", i+1, p)
		}
	}

	fmt.Println("\nAliases:")
	if len(config.Aliases) == 0 {
		fmt.Println("   (none - use 'timetrack alias <short> <full name>')")
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

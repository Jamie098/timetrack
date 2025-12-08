package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

func main() {
	config := loadConfig()
	data := loadData()
	day := getTodayData(data, config)

	// Save day data if it's new (has recurring meetings applied)
	if _, exists := data[today()]; !exists && day.ExcludedPercent > 0 {
		data[today()] = day
		saveData(data)
	}

	if len(os.Args) < 2 {
		// Show today's status by default
		printStatus(day)
		return
	}

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "help", "-h", "--help":
		printHelp()

	case "interactive", "i":
		// Explicit interactive mode
		runInteractive(data, config, day)

	case "add":
		targetDate, args, err := getTargetDate(os.Args[2:], "add")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(args) < 2 {
			fmt.Println("Usage: timetrack add <project> <hours> [--date YYYY-MM-DD]")
			return
		}
		project := resolveProjectWithSuggestions(args[0], config, true)
		hours, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Println("Invalid hours:", args[1])
			return
		}

		// Get data for target date
		targetDay := getDateData(data, config, targetDate)

		// Validate time
		pct := hoursToPercent(hours)
		if pct > 100 {
			fmt.Printf("⚠️  Warning: %.2f hours is %.1f%% of an 8-hour day (>100%%). Did you mean %.2f hours?\n",
				hours, pct, hours/10)
		}

		if targetDay.Projects == nil {
			targetDay.Projects = make(map[string]float64)
		}
		targetDay.Projects[project] = pct
		targetDay.LastModified = project // Track for undo
		data[targetDate] = targetDay
		saveData(data)

		// Check total allocation
		total := getTotalTracked(targetDay)
		available := getAvailablePercent(targetDay)
		if total > available {
			fmt.Printf("Added %.1f%% to %s", pct, project)
			if targetDate != today() {
				fmt.Printf(" on %s", targetDate)
			}
			fmt.Println()
			fmt.Printf("⚠️  Warning: Over-allocated by %.1f%%!\n", total-available)
		} else {
			fmt.Printf("Added %.1f%% to %s", pct, project)
			if targetDate != today() {
				fmt.Printf(" on %s", targetDate)
			}
			fmt.Println()
		}
		printStatus(targetDay)

	case "fill":
		targetDate, args, err := getTargetDate(os.Args[2:], "fill")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(args) < 1 {
			fmt.Println("Usage: timetrack fill <project> [--date YYYY-MM-DD]")
			return
		}
		project := resolveProjectWithSuggestions(args[0], config, true)

		// Get data for target date
		targetDay := getDateData(data, config, targetDate)

		available := getAvailablePercent(targetDay)
		tracked := getTotalTracked(targetDay)
		remaining := available - tracked

		if remaining <= 0 {
			fmt.Printf("⚠️  No remaining time to fill (%.1f%% available, %.1f%% already tracked)", available, tracked)
			if targetDate != today() {
				fmt.Printf(" on %s", targetDate)
			}
			fmt.Println()
			return
		}

		if targetDay.Projects == nil {
			targetDay.Projects = make(map[string]float64)
		}
		targetDay.Projects[project] = remaining
		targetDay.LastModified = project // Track for undo
		data[targetDate] = targetDay
		saveData(data)

		remainingHours := remaining / 100.0 * 8.0
		fmt.Printf("Filled remaining %.1f%% (%.2f hours) to %s", remaining, remainingHours, project)
		if targetDate != today() {
			fmt.Printf(" on %s", targetDate)
		}
		fmt.Println()
		printStatus(targetDay)

	case "exclude", "ex":
		if len(os.Args) < 4 {
			fmt.Println("Usage: timetrack exclude <meeting-name> <hours>")
			return
		}
		name := os.Args[2]
		hours, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Invalid hours:", os.Args[3])
			return
		}
		pct := hoursToPercent(hours)
		if day.ExcludedMeetings == nil {
			day.ExcludedMeetings = make(map[string]float64)
		}

		oldPct := day.ExcludedMeetings[name]
		day.ExcludedPercent = day.ExcludedPercent - oldPct + pct
		day.ExcludedMeetings[name] = pct

		data[today()] = day
		saveData(data)
		fmt.Printf("Excluded %.1f%% for %s\n", pct, name)
		printStatus(day)

	case "rm", "remove":
		targetDate, args, err := getTargetDate(os.Args[2:], "rm")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(args) < 1 {
			fmt.Println("Usage: timetrack rm <project> [--date YYYY-MM-DD]")
			return
		}
		project := resolveProject(args[0], config)

		// Get data for target date
		targetDay := getDateData(data, config, targetDate)

		if _, ok := targetDay.Projects[project]; ok {
			delete(targetDay.Projects, project)
			data[targetDate] = targetDay
			saveData(data)
			fmt.Printf("Removed %s", project)
			if targetDate != today() {
				fmt.Printf(" from %s", targetDate)
			}
			fmt.Println()
			printStatus(targetDay)
		} else {
			fmt.Printf("Project '%s' not found", project)
			if targetDate != today() {
				fmt.Printf(" on %s", targetDate)
			}
			fmt.Println()
		}

	case "rmex":
		if len(os.Args) < 3 {
			fmt.Println("Usage: timetrack rmex <meeting-name>")
			return
		}
		name := os.Args[2]
		if pct, ok := day.ExcludedMeetings[name]; ok {
			day.ExcludedPercent -= pct
			delete(day.ExcludedMeetings, name)
			data[today()] = day
			saveData(data)
			fmt.Printf("Removed excluded meeting: %s\n", name)
			printStatus(day)
		} else {
			fmt.Printf("Excluded meeting '%s' not found\n", name)
		}

	case "clear":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Clear today's data? (y/n): ")
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(input)) == "y" {
			delete(data, today())
			saveData(data)
			fmt.Println("Cleared today's data")
		}

	case "show", "cal", "calendar":
		days := 7
		if len(os.Args) >= 3 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				days = d
			}
		}
		printCalendar(data, days)

	case "history", "hist":
		// Legacy command - redirect to show
		days := 7
		if len(os.Args) >= 3 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				days = d
			}
		}
		printCalendar(data, days)

	case "config":
		if len(os.Args) >= 3 && os.Args[2] == "edit" {
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.Command("notepad", getConfigPath())
			} else {
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "nano"
				}
				cmd = exec.Command(editor, getConfigPath())
			}
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			printConfig(config)
		}

	case "meeting":
		if len(os.Args) < 3 {
			fmt.Println("Usage: timetrack meeting add <name> <hours> <days>")
			fmt.Println("       timetrack meeting rm <name>")
			return
		}
		subcmd := os.Args[2]
		switch subcmd {
		case "add":
			if len(os.Args) < 6 {
				fmt.Println("Usage: timetrack meeting add <name> <hours> <days>")
				fmt.Println("Days: mon,tue,wed,thu,fri,sat,sun,daily,weekdays")
				return
			}
			name := os.Args[3]
			hours, err := strconv.ParseFloat(os.Args[4], 64)
			if err != nil {
				fmt.Println("Invalid hours:", os.Args[4])
				return
			}
			pct := hoursToPercent(hours)
			days := strings.Split(strings.ToLower(os.Args[5]), ",")

			// Check if meeting already exists, update it
			found := false
			for i, m := range config.RecurringMeetings {
				if m.Name == name {
					config.RecurringMeetings[i].Percent = pct
					config.RecurringMeetings[i].Days = days
					found = true
					break
				}
			}
			if !found {
				config.RecurringMeetings = append(config.RecurringMeetings, RecurringMeeting{
					Name:    name,
					Percent: pct,
					Days:    days,
				})
			}
			saveConfig(config)
			fmt.Printf("Added recurring meeting: %s (%.1f%% on %s)\n", name, pct, strings.Join(days, ", "))

		case "rm", "remove":
			if len(os.Args) < 4 {
				fmt.Println("Usage: timetrack meeting rm <name>")
				return
			}
			name := os.Args[3]
			found := false
			for i, m := range config.RecurringMeetings {
				if m.Name == name {
					config.RecurringMeetings = append(config.RecurringMeetings[:i], config.RecurringMeetings[i+1:]...)
					found = true
					break
				}
			}
			if found {
				saveConfig(config)
				fmt.Printf("Removed recurring meeting: %s\n", name)
			} else {
				fmt.Printf("Meeting '%s' not found\n", name)
			}
		}

	case "reminder":
		if len(os.Args) < 3 {
			fmt.Println("Current reminder times:", strings.Join(config.ReminderTimes, ", "))
			fmt.Println("Usage: timetrack reminder 09:00,12:00,15:00")
			return
		}
		times := strings.Split(os.Args[2], ",")
		config.ReminderTimes = times
		saveConfig(config)
		fmt.Println("Reminder times set to:", strings.Join(times, ", "))

	case "start":
		if isDaemonRunning() {
			fmt.Println("Daemon is already running")
			return
		}
		runDaemon()

	case "start-bg":
		if isDaemonRunning() {
			fmt.Println("Daemon is already running")
			return
		}
		exe, _ := os.Executable()
		cmd := exec.Command(exe, "start")
		cmd.Start()
		fmt.Println("Reminder service started in background")
		fmt.Printf("PID: %d\n", cmd.Process.Pid)

	case "stop":
		stopDaemon()

	case "status":
		if isDaemonRunning() {
			pidBytes, _ := os.ReadFile(getPidPath())
			fmt.Printf("Reminder service is running (PID: %s)\n", strings.TrimSpace(string(pidBytes)))
		} else {
			fmt.Println("Reminder service is not running")
		}

	case "test-notify":
		sendNotification("⏰ TimeTrack Test", "Notifications are working!")
		fmt.Println("Test notification sent")

	case "projects":
		if len(os.Args) < 3 {
			fmt.Println("Usage: timetrack projects set \"Proj1,Proj2,Proj3,...\"")
			fmt.Println("       timetrack projects parse \"Date,Proj1,Proj2,...,Total\"")
			fmt.Println("       timetrack projects list")
			return
		}
		subcmd := os.Args[2]
		switch subcmd {
		case "set":
			if len(os.Args) < 4 {
				fmt.Println("Usage: timetrack projects set \"Proj1,Proj2,Proj3,...\"")
				return
			}
			projects := strings.Split(os.Args[3], ",")
			for i, p := range projects {
				projects[i] = strings.TrimSpace(p)
			}
			config.Projects = projects
			saveConfig(config)
			fmt.Printf("Set %d project columns\n", len(projects))
		case "parse":
			if len(os.Args) < 4 {
				fmt.Println("Usage: timetrack projects parse \"Date,Proj1,Proj2,...,Total Time Spent\"")
				fmt.Println("Paste the full header row from Excel - first and last columns are stripped automatically")
				return
			}
			cols := strings.Split(os.Args[3], ",")
			for i, p := range cols {
				cols[i] = strings.TrimSpace(p)
			}
			// Strip first and last columns
			if len(cols) > 2 {
				cols = cols[1 : len(cols)-1]
			}
			config.Projects = cols

			// Auto-generate aliases
			if config.Aliases == nil {
				config.Aliases = make(map[string]string)
			}
			for _, project := range cols {
				alias := generateAlias(project, config.Aliases)
				config.Aliases[alias] = project
			}

			saveConfig(config)
			fmt.Printf("Set %d project columns with aliases:\n", len(cols))
			for _, p := range cols {
				// Find alias for this project
				for alias, name := range config.Aliases {
					if name == p {
						fmt.Printf("  %s → %s\n", alias, p)
						break
					}
				}
			}
		case "list":
			// Auto-discover projects from all tracked time
			projects := getAllProjects(data)
			if len(projects) == 0 {
				fmt.Println("No projects found in tracked time")
			} else {
				fmt.Println("Projects (auto-discovered from your time tracking):")
				for i, p := range projects {
					fmt.Printf("%d. %s\n", i+1, p)
				}
			}
		}

	case "alias":
		if len(os.Args) < 3 {
			fmt.Println("Usage: timetrack alias <short> <full project name>")
			fmt.Println("       timetrack alias rm <short>")
			fmt.Println("       timetrack alias list")
			return
		}
		if os.Args[2] == "list" {
			if len(config.Aliases) == 0 {
				fmt.Println("No aliases configured")
			} else {
				aliases := make([]string, 0, len(config.Aliases))
				for k := range config.Aliases {
					aliases = append(aliases, k)
				}
				sort.Strings(aliases)
				for _, k := range aliases {
					fmt.Printf("%s → %s\n", k, config.Aliases[k])
				}
			}
			return
		}
		if os.Args[2] == "rm" && len(os.Args) >= 4 {
			short := strings.ToLower(os.Args[3])
			if _, ok := config.Aliases[short]; ok {
				delete(config.Aliases, short)
				saveConfig(config)
				fmt.Printf("Removed alias: %s\n", short)
			} else {
				fmt.Printf("Alias '%s' not found\n", short)
			}
			return
		}
		if len(os.Args) < 4 {
			fmt.Println("Usage: timetrack alias <short> <full project name>")
			return
		}
		short := strings.ToLower(os.Args[2])
		full := strings.Join(os.Args[3:], " ")
		config.Aliases[short] = full
		saveConfig(config)
		fmt.Printf("Alias set: %s → %s\n", short, full)

	case "import":
		if len(os.Args) < 3 {
			fmt.Println("Usage: timetrack import <csv-file>")
			return
		}
		if err := importFromCSV(os.Args[2], config); err != nil {
			fmt.Println("Import failed:", err)
		}

	case "export":
		format := "csv"
		filename := ""

		if len(os.Args) >= 3 {
			format = strings.ToLower(os.Args[2])
		}

		switch format {
		case "json":
			if len(os.Args) >= 4 {
				filename = os.Args[3]
			} else {
				filename = "timetrack-export.json"
			}
			if err := exportToJSON(data, filename); err != nil {
				fmt.Println("Export failed:", err)
			}

		case "csv", "week":
			if len(os.Args) >= 4 {
				filename = os.Args[3]
			} else {
				filename = "timetrack-week.csv"
			}
			if err := exportWeekToCSV(data, config, filename); err != nil {
				fmt.Println("Export failed:", err)
			}

		case "all":
			if len(os.Args) >= 4 {
				filename = os.Args[3]
			} else {
				filename = "timetrack-all.csv"
			}
			if err := exportAllToCSV(data, config, filename); err != nil {
				fmt.Println("Export failed:", err)
			}

		default:
			fmt.Println("Unknown export format:", format)
			fmt.Println("Supported formats: json, csv, week, all")
		}

	case "undo":
		handleUndo(data, &day)

	case "edit":
		targetDate, args, err := getTargetDate(os.Args[2:], "edit")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(args) < 2 {
			fmt.Println("Usage: timetrack edit <project> <hours> [--date YYYY-MM-DD]")
			return
		}
		project := resolveProjectWithSuggestions(args[0], config, true)
		hours, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			fmt.Println("Invalid hours:", args[1])
			return
		}

		// Get data for target date
		targetDay := getDateData(data, config, targetDate)

		if _, ok := targetDay.Projects[project]; !ok {
			fmt.Printf("Project '%s' not found in tracking", project)
			if targetDate != today() {
				fmt.Printf(" for %s", targetDate)
			}
			fmt.Println()
			fmt.Println("Use 'add' to create a new entry")
			return
		}

		pct := hoursToPercent(hours)
		targetDay.Projects[project] = pct
		targetDay.LastModified = project // Track for undo
		data[targetDate] = targetDay
		saveData(data)
		fmt.Printf("Updated %s to %.1f%%", project, pct)
		if targetDate != today() {
			fmt.Printf(" on %s", targetDate)
		}
		fmt.Println()
		printStatus(targetDay)

	case "summary", "sum":
		// Legacy command - show today's status
		printStatus(day)

	case "report":
		fmt.Println("Note: 'report' commands are deprecated. Use 'show' to view entries and 'export' for CSV reports.")
		fmt.Println()

		if len(os.Args) < 3 {
			fmt.Println("Usage:")
			fmt.Println("  timetrack show [days]          - View calendar (recommended)")
			fmt.Println("  timetrack export csv           - Export week to CSV")
			fmt.Println("  timetrack export all           - Export all data")
			return
		}

		reportType := strings.ToLower(os.Args[2])
		switch reportType {
		case "week", "weekly":
			generateWeeklyReport(data)
		case "project", "proj":
			if len(os.Args) < 4 {
				fmt.Println("Usage: timetrack report project <name>")
				return
			}
			generateProjectReport(data, os.Args[3], config)
		case "stats", "statistics":
			generateStatsReport(data)
		default:
			fmt.Println("Unknown report type:", reportType)
			fmt.Println("Try: timetrack show [days]")
		}

	case "url":
		if len(os.Args) < 3 {
			if config.TimesheetURL == "" {
				fmt.Println("No timesheet URL configured")
				fmt.Println("Usage: timetrack url set <url>")
			} else {
				fmt.Println("Timesheet URL:", config.TimesheetURL)
			}
			return
		}

		subcmd := strings.ToLower(os.Args[2])
		switch subcmd {
		case "set":
			if len(os.Args) < 4 {
				fmt.Println("Usage: timetrack url set <url>")
				return
			}
			config.TimesheetURL = os.Args[3]
			saveConfig(config)
			fmt.Println("Timesheet URL set to:", config.TimesheetURL)

		case "open":
			if config.TimesheetURL == "" {
				fmt.Println("No timesheet URL configured")
				fmt.Println("Use: timetrack url set <url>")
				return
			}
			openURL(config.TimesheetURL)

		case "rm", "remove", "clear":
			config.TimesheetURL = ""
			saveConfig(config)
			fmt.Println("Timesheet URL cleared")

		default:
			fmt.Println("Unknown url command:", subcmd)
			fmt.Println("Available: set, open, rm")
		}

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printHelp()
	}
}

func handleUndo(data map[string]DayData, day *DayData) {
	if len(day.Projects) == 0 {
		fmt.Println("No entries to undo")
		return
	}

	// Use the tracked last modified project
	lastProject := day.LastModified

	// If no tracking info (legacy data), fall back to asking user
	if lastProject == "" || day.Projects[lastProject] == 0 {
		if len(day.Projects) == 1 {
			// Only one project, remove it
			for p := range day.Projects {
				lastProject = p
				break
			}
		} else {
			fmt.Println("Cannot determine last added project. Current projects:")
			i := 1
			for name, pct := range day.Projects {
				fmt.Printf("  %d. %s (%.1f%%)\n", i, name, pct)
				i++
			}
			fmt.Println("\nUse 'timetrack rm <project>' to remove a specific project")
			return
		}
	}

	if lastProject != "" {
		pct := day.Projects[lastProject]
		delete(day.Projects, lastProject)
		day.LastModified = "" // Clear tracking
		data[today()] = *day
		saveData(data)
		fmt.Printf("Removed %s (%.1f%%)\n", lastProject, pct)
		printStatus(*day)
	}
}

func printSummary(day DayData) {
	available := getAvailablePercent(day)
	tracked := getTotalTracked(day)
	remaining := available - tracked

	status := "✓"
	if remaining < 0 {
		status = "⚠️"
	} else if remaining > 20 {
		status = "⏳"
	}

	fmt.Printf("%s %s: %.1f%% tracked, %.1f%% remaining", status, day.Date, tracked, remaining)
	if len(day.Projects) > 0 {
		fmt.Print(" (")
		i := 0
		for name, pct := range day.Projects {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s:%.0f%%", name, pct)
			i++
			if i >= 3 {
				fmt.Printf(" +%d more", len(day.Projects)-3)
				break
			}
		}
		fmt.Print(")")
	}
	fmt.Println()
}

func printHistory(data map[string]DayData, days int) {
	fmt.Println()
	fmt.Println("📆 History")
	fmt.Println(strings.Repeat("─", 60))

	dates := make([]string, 0)
	for date := range data {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))

	count := 0
	for _, date := range dates {
		if count >= days {
			break
		}
		day := data[date]
		available := getAvailablePercent(day)
		tracked := getTotalTracked(day)

		fmt.Printf("\n%s  Available: %.1f%%  Tracked: %.1f%%\n", date, available, tracked)

		if len(day.Projects) > 0 {
			projects := sortedKeys(day.Projects)
			for _, name := range projects {
				fmt.Printf("   • %s: %.1f%%\n", name, day.Projects[name])
			}
		}
		count++
	}

	if count == 0 {
		fmt.Println("No history found")
	}
	fmt.Println()
}

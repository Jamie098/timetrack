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
		printStatus(day)
		return
	}

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "help", "-h", "--help":
		printHelp()

	case "add":
		if len(os.Args) < 4 {
			fmt.Println("Usage: timetrack add <project> <percent>")
			return
		}
		project := resolveProject(os.Args[2], config)
		pct, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Invalid percentage:", os.Args[3])
			return
		}
		if day.Projects == nil {
			day.Projects = make(map[string]float64)
		}
		day.Projects[project] = pct
		data[today()] = day
		saveData(data)
		fmt.Printf("Added %.1f%% to %s\n", pct, project)
		printStatus(day)

	case "exclude", "ex":
		if len(os.Args) < 4 {
			fmt.Println("Usage: timetrack exclude <meeting-name> <percent>")
			return
		}
		name := os.Args[2]
		pct, err := strconv.ParseFloat(os.Args[3], 64)
		if err != nil {
			fmt.Println("Invalid percentage:", os.Args[3])
			return
		}
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
		if len(os.Args) < 3 {
			fmt.Println("Usage: timetrack rm <project>")
			return
		}
		project := resolveProject(os.Args[2], config)
		if _, ok := day.Projects[project]; ok {
			delete(day.Projects, project)
			data[today()] = day
			saveData(data)
			fmt.Printf("Removed %s\n", project)
			printStatus(day)
		} else {
			fmt.Printf("Project '%s' not found\n", project)
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

	case "history", "hist":
		days := 7
		if len(os.Args) >= 3 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				days = d
			}
		}
		printHistory(data, days)

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
			fmt.Println("Usage: timetrack meeting add <name> <percent> <days>")
			fmt.Println("       timetrack meeting rm <name>")
			return
		}
		subcmd := os.Args[2]
		switch subcmd {
		case "add":
			if len(os.Args) < 6 {
				fmt.Println("Usage: timetrack meeting add <name> <percent> <days>")
				fmt.Println("Days: mon,tue,wed,thu,fri,sat,sun,daily,weekdays")
				return
			}
			name := os.Args[3]
			pct, err := strconv.ParseFloat(os.Args[4], 64)
			if err != nil {
				fmt.Println("Invalid percentage:", os.Args[4])
				return
			}
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
		case "list":
			if len(config.Projects) == 0 {
				fmt.Println("No projects configured")
			} else {
				for i, p := range config.Projects {
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

	case "week":
		printWeekExport(data, config)

	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printHelp()
	}
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

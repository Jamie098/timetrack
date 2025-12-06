package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func runInteractive(data map[string]DayData, config Config, day DayData) {
	reader := bufio.NewReader(os.Stdin)

	for {
		printStatus(day)
		fmt.Println()
		fmt.Println("What would you like to do?")
		fmt.Println("  1. Add time to project")
		fmt.Println("  2. Exclude meeting time")
		fmt.Println("  3. Remove project")
		fmt.Println("  4. View history")
		fmt.Println("  5. View config")
		fmt.Println("  6. Exit")
		fmt.Print("\nChoice: ")

		input, _ := reader.ReadString('\n')
		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			handleInteractiveAdd(reader, data, config, &day)
		case "2":
			handleInteractiveExclude(reader, data, &day)
		case "3":
			handleInteractiveRemove(reader, data, config, &day)
		case "4":
			handleInteractiveHistory(reader, data)
		case "5":
			printConfig(config)
			fmt.Print("\nPress Enter to continue...")
			reader.ReadString('\n')
		case "6", "q", "quit", "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
			fmt.Print("Press Enter to continue...")
			reader.ReadString('\n')
		}
	}
}

func handleInteractiveAdd(reader *bufio.Reader, data map[string]DayData, config Config, day *DayData) {
	fmt.Print("\nProject name (or alias): ")
	projectInput, _ := reader.ReadString('\n')
	projectName := strings.TrimSpace(projectInput)
	if projectName == "" {
		return
	}

	project := resolveProjectWithSuggestions(projectName, config, true)

	fmt.Print("Hours: ")
	hoursInput, _ := reader.ReadString('\n')
	hours, err := strconv.ParseFloat(strings.TrimSpace(hoursInput), 64)
	if err != nil {
		fmt.Println("Invalid hours")
		fmt.Print("Press Enter to continue...")
		reader.ReadString('\n')
		return
	}

	pct := hoursToPercent(hours)

	// Validate time
	if pct > 100 {
		fmt.Printf("⚠️  Warning: %.2f hours is %.1f%% of an 8-hour day (>100%%). Did you mean %.2f hours?\n",
			hours, pct, hours/10)
	}

	if day.Projects == nil {
		day.Projects = make(map[string]float64)
	}
	day.Projects[project] = pct
	day.LastModified = project // Track for undo
	data[today()] = *day
	saveData(data)

	// Check total allocation
	total := getTotalTracked(*day)
	available := getAvailablePercent(*day)
	if total > available {
		fmt.Printf("\n✓ Added %.1f%% to %s\n", pct, project)
		fmt.Printf("⚠️  Warning: Over-allocated by %.1f%%!\n", total-available)
	} else {
		fmt.Printf("\n✓ Added %.1f%% to %s\n", pct, project)
	}
	fmt.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

func handleInteractiveExclude(reader *bufio.Reader, data map[string]DayData, day *DayData) {
	fmt.Print("\nMeeting name: ")
	nameInput, _ := reader.ReadString('\n')
	name := strings.TrimSpace(nameInput)
	if name == "" {
		return
	}

	fmt.Print("Hours: ")
	hoursInput, _ := reader.ReadString('\n')
	hours, err := strconv.ParseFloat(strings.TrimSpace(hoursInput), 64)
	if err != nil {
		fmt.Println("Invalid hours")
		fmt.Print("Press Enter to continue...")
		reader.ReadString('\n')
		return
	}

	pct := hoursToPercent(hours)

	if day.ExcludedMeetings == nil {
		day.ExcludedMeetings = make(map[string]float64)
	}

	oldPct := day.ExcludedMeetings[name]
	day.ExcludedPercent = day.ExcludedPercent - oldPct + pct
	day.ExcludedMeetings[name] = pct

	data[today()] = *day
	saveData(data)

	fmt.Printf("\n✓ Excluded %.1f%% for %s\n", pct, name)
	fmt.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

func handleInteractiveRemove(reader *bufio.Reader, data map[string]DayData, config Config, day *DayData) {
	if len(day.Projects) == 0 {
		fmt.Println("\nNo projects to remove")
		fmt.Print("Press Enter to continue...")
		reader.ReadString('\n')
		return
	}

	fmt.Println("\nCurrent projects:")
	i := 1
	projectList := make([]string, 0)
	for name := range day.Projects {
		fmt.Printf("  %d. %s\n", i, name)
		projectList = append(projectList, name)
		i++
	}

	fmt.Print("\nProject number to remove (or name): ")
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)

	var projectToRemove string
	if num, err := strconv.Atoi(choice); err == nil && num > 0 && num <= len(projectList) {
		projectToRemove = projectList[num-1]
	} else {
		projectToRemove = resolveProject(choice, config)
	}

	if _, ok := day.Projects[projectToRemove]; ok {
		delete(day.Projects, projectToRemove)
		data[today()] = *day
		saveData(data)
		fmt.Printf("\n✓ Removed %s\n", projectToRemove)
	} else {
		fmt.Printf("\nProject '%s' not found\n", projectToRemove)
	}

	fmt.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

func handleInteractiveHistory(reader *bufio.Reader, data map[string]DayData) {
	fmt.Print("\nHow many days? (default 7): ")
	input, _ := reader.ReadString('\n')
	days := 7
	if d, err := strconv.Atoi(strings.TrimSpace(input)); err == nil && d > 0 {
		days = d
	}

	printHistory(data, days)
	fmt.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

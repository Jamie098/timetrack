package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func runDaemon() {
	config := loadConfig()

	fmt.Println("TimeTrack reminder service started")
	fmt.Printf("Reminder times: %v\n", config.ReminderTimes)
	fmt.Println("Press Ctrl+C to stop")

	// Save PID
	os.WriteFile(getPidPath(), []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
	defer os.Remove(getPidPath())

	notifiedToday := make(map[string]bool)
	lastDate := today()

	for {
		// Reset notifications on new day
		if today() != lastDate {
			notifiedToday = make(map[string]bool)
			lastDate = today()
		}

		now := time.Now().Format("15:04")

		for _, reminderTime := range config.ReminderTimes {
			if now == reminderTime && !notifiedToday[reminderTime] {
				notifiedToday[reminderTime] = true

				data := loadData()
				day := getTodayData(data, config)
				available := getAvailablePercent(day)
				tracked := getTotalTracked(day)
				remaining := available - tracked

				var message string
				if remaining > 0 {
					message = fmt.Sprintf("%.1f%% remaining to track today", remaining)
				} else if remaining == 0 {
					message = "Day fully tracked! ✨"
				} else {
					message = fmt.Sprintf("Over-allocated by %.1f%%", -remaining)
				}

				sendNotification("⏰ TimeTrack Reminder", message)
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func isDaemonRunning() bool {
	pidBytes, err := os.ReadFile(getPidPath())
	if err != nil {
		return false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil {
		return false
	}

	// Check if process exists
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist", "/FI", fmt.Sprintf("PID eq %d", pid))
		output, _ := cmd.Output()
		return strings.Contains(string(output), fmt.Sprintf("%d", pid))
	} else {
		process, err := os.FindProcess(pid)
		if err != nil {
			return false
		}
		err = process.Signal(os.Signal(nil))
		return err == nil
	}
}

func stopDaemon() {
	pidBytes, err := os.ReadFile(getPidPath())
	if err != nil {
		fmt.Println("Daemon not running")
		return
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(pidBytes)))
	if err != nil {
		fmt.Println("Invalid PID file")
		return
	}

	if runtime.GOOS == "windows" {
		exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid)).Run()
	} else {
		process, err := os.FindProcess(pid)
		if err == nil {
			process.Kill()
		}
	}

	os.Remove(getPidPath())
	fmt.Println("Daemon stopped")
}

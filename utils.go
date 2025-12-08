package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func generateAlias(project string, existing map[string]string) string {
	// Clean up the project name
	name := strings.ToLower(project)

	// Common replacements
	name = strings.ReplaceAll(name, "ct.gov", "ctgov")
	name = strings.ReplaceAll(name, "sso/aug", "sso")

	// Remove common suffixes/noise
	noise := []string{"-", "_", "and", "the", " "}

	// Split into words
	for _, n := range noise {
		name = strings.ReplaceAll(name, n, " ")
	}
	words := strings.Fields(name)

	if len(words) == 0 {
		return "proj"
	}

	var alias string

	// Single word: take first 4 chars
	if len(words) == 1 {
		alias = words[0]
		if len(alias) > 4 {
			alias = alias[:4]
		}
	} else if len(words) == 2 {
		// Two words: first 2 chars of each
		alias = safeSubstring(words[0], 2) + safeSubstring(words[1], 2)
	} else {
		// 3+ words: first char of first 3-4 words
		for i, w := range words {
			if i >= 4 {
				break
			}
			if len(w) > 0 {
				alias += string(w[0])
			}
		}
		// If too short, pad with more chars from first word
		if len(alias) < 3 && len(words[0]) > 1 {
			alias = safeSubstring(words[0], 3)
		}
	}

	// Ensure uniqueness
	baseAlias := alias
	counter := 1
	for {
		if _, exists := existing[alias]; !exists {
			break
		}
		counter++
		alias = fmt.Sprintf("%s%d", baseAlias, counter)
	}

	return alias
}

func safeSubstring(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length]
}

func openURL(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux, freebsd, openbsd, etc.
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error opening URL:", err)
		fmt.Println("Please open manually:", url)
	} else {
		fmt.Println("Opening timesheet in browser...")
	}
}

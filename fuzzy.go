package main

import (
	"fmt"
	"strings"
)

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}
	m := nums[0]
	for _, n := range nums[1:] {
		if n < m {
			m = n
		}
	}
	return m
}

// fuzzyMatchProject finds the best matching project from the config
func fuzzyMatchProject(input string, config Config) (string, []string) {
	// First, check exact alias match (existing behavior)
	if fullName, ok := config.Aliases[strings.ToLower(input)]; ok {
		return fullName, nil
	}

	// Check exact match in projects
	for _, proj := range config.Projects {
		if strings.EqualFold(input, proj) {
			return proj, nil
		}
	}

	// If no exact match, try fuzzy matching
	type match struct {
		name     string
		distance int
	}

	matches := make([]match, 0)
	inputLower := strings.ToLower(input)

	// Check against all projects
	for _, proj := range config.Projects {
		projLower := strings.ToLower(proj)

		// Check if input is a substring
		if strings.Contains(projLower, inputLower) {
			matches = append(matches, match{proj, 0}) // High priority
			continue
		}

		// Calculate edit distance
		dist := levenshteinDistance(input, proj)
		// Only consider if distance is reasonable (less than half the length)
		if dist <= len(proj)/2 && dist <= 3 {
			matches = append(matches, match{proj, dist})
		}
	}

	// Check against aliases
	for alias, fullName := range config.Aliases {
		if strings.Contains(alias, inputLower) {
			// Check if not already in matches
			found := false
			for _, m := range matches {
				if m.name == fullName {
					found = true
					break
				}
			}
			if !found {
				matches = append(matches, match{fullName, 0})
			}
		}
	}

	if len(matches) == 0 {
		return input, nil // No match found, return original
	}

	// Sort by distance (0 distance first - substring matches)
	var best match
	var suggestions []string

	for _, m := range matches {
		if m.distance == 0 {
			if best.name == "" {
				best = m
			} else {
				suggestions = append(suggestions, m.name)
			}
		}
	}

	// If no substring matches, use closest edit distance
	if best.name == "" {
		best = matches[0]
		for _, m := range matches[1:] {
			if m.distance < best.distance {
				suggestions = append(suggestions, best.name)
				best = m
			} else {
				suggestions = append(suggestions, m.name)
			}
		}
	}

	return best.name, suggestions
}

// resolveProjectWithSuggestions attempts to resolve a project and shows suggestions if ambiguous
func resolveProjectWithSuggestions(input string, config Config, interactive bool) string {
	match, suggestions := fuzzyMatchProject(input, config)

	// If we found a perfect match (through alias or exact name), return it
	if len(suggestions) == 0 && match != input {
		return match
	}

	// If there are suggestions, ask user
	if len(suggestions) > 0 && interactive {
		fmt.Printf("\nDid you mean '%s'?\n", match)
		fmt.Println("Other suggestions:")
		for i, sug := range suggestions {
			fmt.Printf("  %d. %s\n", i+1, sug)
		}
		fmt.Println("\nUsing:", match)
	} else if match != input && interactive {
		fmt.Printf("No exact match found. Using fuzzy match: %s\n", match)
	}

	return match
}

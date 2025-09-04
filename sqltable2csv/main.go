package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Skip separator lines (containing only -, +, and spaces)
		if isSeparatorLine(line) {
			continue
		}

		// Split by pipe and trim spaces
		fields := splitAndTrim(line)
		
		// Write to CSV
		if err := writer.Write(fields); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV: %v\n", err)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}

func isSeparatorLine(line string) bool {
	// Check if line contains only -, +, =, and spaces
	for _, ch := range line {
		if ch != '-' && ch != '+' && ch != '=' && ch != ' ' {
			return false
		}
	}
	return len(strings.TrimSpace(line)) > 0
}

func splitAndTrim(line string) []string {
	// Split by pipe character
	parts := strings.Split(line, "|")
	
	// Trim spaces from each part
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		result = append(result, trimmed)
	}
	
	return result
}
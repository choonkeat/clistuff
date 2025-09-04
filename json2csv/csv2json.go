package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	// Read CSV data from stdin
	reader := csv.NewReader(os.Stdin)

	// Read the header row
	header, err := reader.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading CSV header: %v\n", err)
		os.Exit(1)
	}

	// Prepare the result slice for JSON objects
	var result []map[string]string

	// Read each data row and convert to JSON object
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading CSV row: %v\n", err)
			os.Exit(1)
		}

		// Create a map for each row
		rowMap := make(map[string]string)
		for i, value := range row {
			if i < len(header) {
				rowMap[header[i]] = value
			}
		}

		result = append(result, rowMap)
	}

	// Convert to JSON and print to stdout
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}

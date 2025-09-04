package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

func main() {
	var fixedHeaderString string
	var appendAllHeaders bool
	var excludeHeaderString string
	var mergeHeaderString string
	flag.StringVar(&fixedHeaderString, "headers", "", "comma-separated list of required headers (optional)")
	flag.BoolVar(&appendAllHeaders, "append-all-headers", false, "append all unique keys from all rows to headers")
	flag.StringVar(&excludeHeaderString, "exclude-headers", "", "comma-separated list of headers to exclude")
	flag.StringVar(&mergeHeaderString, "merge-headers", "", "comma-separated assignments like err=error to merge values")
	flag.Parse()
	var headers []string
	if fixedHeaderString != "" {
		headers = strings.Split(fixedHeaderString, ",")
	}
	
	// Parse exclude headers
	var excludeHeaders []string
	if excludeHeaderString != "" {
		excludeHeaders = strings.Split(excludeHeaderString, ",")
	}
	
	// Parse merge headers into a map
	mergeMap := make(map[string][]string) // header -> list of source keys
	mergeSourceKeys := make(map[string]bool) // track all source keys that are used in merges
	if mergeHeaderString != "" {
		for _, assignment := range strings.Split(mergeHeaderString, ",") {
			parts := strings.Split(assignment, "=")
			if len(parts) == 2 {
				header := strings.TrimSpace(parts[0])
				sourceKey := strings.TrimSpace(parts[1])
				mergeMap[header] = append(mergeMap[header], sourceKey)
				// Track source keys to exclude them from being added as separate headers
				if sourceKey != header {
					mergeSourceKeys[sourceKey] = true
				}
			}
		}
	}

	// Read JSON data from stdin
	jsonData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading JSON data: %v\n", err)
		os.Exit(1)
	}

	// Parse JSON into slice of maps
	var data []map[string]any
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		jsonString := "["+strings.Join(strings.Split(strings.TrimSpace(string(jsonData)), "\n"), ",")+"]"
		err = json.Unmarshal([]byte(jsonString), &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing JSON %#v: %v\n", jsonString, err)
			os.Exit(1)
		}
	}

	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "No data found in JSON\n")
		os.Exit(1)
	}

	// Map to count non-empty values for each key
	keyPopularity := make(map[string]int)
	
	if appendAllHeaders {
		// Collect all unique keys from all rows and count non-empty values
		for _, row := range data {
			for key, value := range row {
				// Only track keys that have at least one non-empty value and aren't excluded
				if value != nil && value != "" {
					if !slices.Contains(headers, key) && !slices.Contains(excludeHeaders, key) && !mergeSourceKeys[key] {
						keyPopularity[key]++
					}
				}
			}
		}
		// Add keys sorted by popularity
		type keyCount struct {
			key   string
			count int
		}
		var sortedKeys []keyCount
		for k, c := range keyPopularity {
			sortedKeys = append(sortedKeys, keyCount{k, c})
		}
		// Sort by count (descending) then by key name (ascending) for stability
		slices.SortFunc(sortedKeys, func(a, b keyCount) int {
			if a.count != b.count {
				return b.count - a.count // Higher count first
			}
			return strings.Compare(a.key, b.key)
		})
		for _, kc := range sortedKeys {
			headers = append(headers, kc.key)
		}
	} else {
		// Extract headers from the first object only and count popularity
		for _, row := range data {
			for key, value := range row {
				if value != nil && value != "" {
					keyPopularity[key]++
				}
			}
		}
		
		// Collect keys from first row that aren't in fixed headers
		var additionalKeys []string
		for key := range data[0] {
			if !slices.Contains(headers, key) && !slices.Contains(excludeHeaders, key) && !mergeSourceKeys[key] {
				additionalKeys = append(additionalKeys, key)
			}
		}
		
		// Sort additional keys by popularity
		slices.SortFunc(additionalKeys, func(a, b string) int {
			countA := keyPopularity[a]
			countB := keyPopularity[b]
			if countA != countB {
				return countB - countA // Higher count first
			}
			return strings.Compare(a, b)
		})
		
		// Append sorted additional keys to headers
		headers = append(headers, additionalKeys...)
	}

	// Remove excluded headers from the final list
	var filteredHeaders []string
	for _, h := range headers {
		if !slices.Contains(excludeHeaders, h) {
			filteredHeaders = append(filteredHeaders, h)
		}
	}
	headers = filteredHeaders

	// Create CSV writer
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Write header row
	err = writer.Write(headers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing CSV header: %v\n", err)
		os.Exit(1)
	}

	// Write data rows
	for rowIndex, row := range data {
		// Check for keys in row that aren't in headers (only warn if not using append-all-headers)
		if !appendAllHeaders {
			for key := range row {
				if !slices.Contains(headers, key) {
					fmt.Fprintf(os.Stderr, "Warning: row %d contains key %q not in headers\n", rowIndex+1, key)
				}
			}
		}

		var values []string
		for _, header := range headers {
			var value any
			
			// Check if this is a merge header
			if sourceKeys, isMerge := mergeMap[header]; isMerge {
				// Try the header itself first
				if v, exists := row[header]; exists && v != nil && v != "" {
					value = v
				} else {
					// Try each source key in order
					for _, sourceKey := range sourceKeys {
						if v, exists := row[sourceKey]; exists && v != nil && v != "" {
							value = v
							break
						}
					}
				}
			} else {
				value = row[header]
			}
			
			switch v := value.(type) {
			case float64:
				// JSON numbers are parsed as float64, check if it's actually an integer
				if v == float64(int64(v)) {
					values = append(values, fmt.Sprintf("%d", int64(v)))
				} else {
					values = append(values, fmt.Sprintf("%g", v))
				}
			case float32:
				if v == float32(int32(v)) {
					values = append(values, fmt.Sprintf("%d", int32(v)))
				} else {
					values = append(values, fmt.Sprintf("%g", v))
				}
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				values = append(values, fmt.Sprintf("%d", v))
			case string:
				values = append(values, v)
			case nil:
				values = append(values, "")
			default:
				values = append(values, fmt.Sprintf("%#v", v))
			}
		}
		err = writer.Write(values)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing CSV row: %v\n", err)
			os.Exit(1)
		}
	}
}

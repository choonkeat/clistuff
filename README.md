# CLI Tools

A collection of command-line utilities for data transformation and processing.

## json2csv / csv2json

Bidirectional converters between JSON and CSV formats.

### json2csv

Converts JSON data to CSV format.

#### Usage
```bash
json2csv < input.json > output.csv
```

#### Features
- Handles nested JSON structures
- Automatically detects headers from JSON keys
- Supports arrays of objects
- Pipes JSON from stdin and outputs CSV to stdout

#### Example
```bash
echo '[{"name":"Alice","age":30},{"name":"Bob","age":25}]' | json2csv
```

Output:
```csv
name,age
Alice,30
Bob,25
```

### csv2json

Converts CSV data to JSON format.

#### Usage
```bash
csv2json < input.csv > output.json
```

#### Features
- Parses CSV with headers
- Outputs array of JSON objects
- Handles quoted values and commas in fields
- Pipes CSV from stdin and outputs JSON to stdout

#### Example
```bash
echo -e "name,age\nAlice,30\nBob,25" | csv2json
```

Output:
```json
[{"name":"Alice","age":"30"},{"name":"Bob","age":"25"}]
```

## oneline

Converts multi-line input into a single line format.

### Usage
```bash
oneline < input.txt
```

### Features
- Removes line breaks and combines text into one line
- Useful for log processing and text manipulation
- Preserves spacing between words
- Reads from stdin and outputs to stdout

### Example
```bash
echo -e "This is\na multi-line\ntext" | oneline
```

Output:
```text
This is a multi-line text
```

## sqltable2csv

Converts SQL table output to CSV format.

### Usage
```bash
sqltable2csv < sql_output.txt > output.csv
```

### Features
- Parses ASCII-formatted SQL query results
- Handles table borders and separators
- Extracts headers and data rows
- Converts to clean CSV format

### Example

Input (SQL table format):
```text
+------+-----+
| name | age |
+------+-----+
| Alice| 30  |
| Bob  | 25  |
+------+-----+
```

Command:
```bash
cat sql_output.txt | sqltable2csv
```

Output:
```csv
name,age
Alice,30
Bob,25
```

## Installation

Place the scripts in your PATH or run them directly from the current directory.

## Requirements

- Standard Unix/Linux utilities
- Shell environment (bash/sh)

## License

See individual script files for licensing information.
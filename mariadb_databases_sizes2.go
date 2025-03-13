package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Get database credentials from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	if user == "" || password == "" {
		log.Fatalf("Database credentials are not set in environment variables")
	}

	// Define the MySQL query to execute
	query := `SELECT table_schema AS "Database", ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS "Size (MB)" FROM information_schema.tables GROUP BY table_schema;`

	// Define the command to run
	cmd := exec.Command("mysql", "-u", user, fmt.Sprintf("-p%s", password), "-e", query)

	// Run the command and capture the overall output
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run command: %s\n", err)
	}

	// Process the output (split by lines)
	lines := strings.Split(string(output), "\n")

	// Print header for formatted output
	fmt.Printf("%-30s %s\n", "Database", "Size (MB)")
	fmt.Println(strings.Repeat("-", 45))

	// Skip the first line (header) and loop through the rest
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split line by whitespace to extract database name and size
		columns := strings.Fields(line)
		if len(columns) >= 2 {
			database := columns[0]
			size := columns[1]
			fmt.Printf("%-30s %s\n", database, size)
		}
	}
}

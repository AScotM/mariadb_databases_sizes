package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type DatabaseSize struct {
	Name string  `json:"database"`
	MB   float64 `json:"size_mb"`
}

func main() {
	// Get configuration from environment
	config := getConfig()

	// Execute MySQL query
	output, err := runMySQLQuery(config)
	if err != nil {
		log.Fatalf("MySQL query failed: %v", err)
	}

	// Parse and display results
	results := parseResults(output)
	printResults(results, config.outputFormat)
}

type Config struct {
	user        string
	password    string
	host        string
	port        string
	outputFormat string
}

func getConfig() Config {
	cfg := Config{
		user:        os.Getenv("DB_USER"),
		password:    os.Getenv("DB_PASSWORD"),
		host:        os.Getenv("DB_HOST"),
		port:        os.Getenv("DB_PORT"),
		outputFormat: os.Getenv("OUTPUT_FORMAT"),
	}

	if cfg.user == "" || cfg.password == "" {
		log.Fatal("DB_USER and DB_PASSWORD environment variables must be set")
	}

	if cfg.host == "" {
		cfg.host = "localhost"
	}
	if cfg.port == "" {
		cfg.port = "3306"
	}
	if cfg.outputFormat == "" {
		cfg.outputFormat = "table"
	}

	return cfg
}

func runMySQLQuery(cfg Config) (string, error) {
	query := `SELECT table_schema AS "Database", 
              ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS "Size (MB)" 
              FROM information_schema.tables 
              GROUP BY table_schema;`

	args := []string{
		"-u", cfg.user,
		"--password=" + cfg.password,
		"-h", cfg.host,
		"-P", cfg.port,
		"-e", query,
	}

	cmd := exec.Command("mysql", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func parseResults(output string) []DatabaseSize {
	var results []DatabaseSize

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Database") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			size, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				continue
			}
			results = append(results, DatabaseSize{
				Name: fields[0],
				MB:   size,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].MB > results[j].MB
	})

	return results
}

func printResults(results []DatabaseSize, format string) {
	switch format {
	case "json":
		jsonOutput, _ := json.MarshalIndent(results, "", "  ")
		fmt.Println(string(jsonOutput))
	default:
		fmt.Printf("%-30s %s\n", "Database", "Size (MB)")
		fmt.Println(strings.Repeat("-", 45))
		for _, db := range results {
			fmt.Printf("%-30s %.2f\n", db.Name, db.MB)
		}
	}
}

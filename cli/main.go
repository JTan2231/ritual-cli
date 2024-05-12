package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type UserEntry struct {
	Text string `json:"text"`
	Date string `json:"date"`
}

func main() {
	var ROOT string
	if dir, err := os.UserHomeDir(); err != nil {
		fmt.Println("Error retrieving home directory:", err)
		os.Exit(1)
	} else {
		ROOT = dir + "/.ritual/"
	}

	MEMORY := ROOT + "entries.json"
	if _, ok := os.LookupEnv("RITUAL_CLI_KEY"); !ok {
		fmt.Println("Error: The required environment variable `RITUAL_CLI_KEY` is not set. Generate a new one at <link>.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("usage: ritual \"your entry\"")
		os.Exit(1)
	}

	// load/create local memory
	if _, err := os.Stat(ROOT); os.IsNotExist(err) {
		err := os.MkdirAll(ROOT, 0755)
		if err != nil {
			fmt.Println("Error creating ROOT directory:", err)
			os.Exit(1)
		}
	}

	var entries []UserEntry
	if _, err := os.Stat(MEMORY); err == nil {
		data, err := os.ReadFile(MEMORY)
		if err != nil {
			fmt.Println("Error reading file:", err)
			os.Exit(1)
		}

		if err := json.Unmarshal(data, &entries); err != nil {
			fmt.Println("Error reading Ritual memory:", err)
			os.Exit(1)
		}

		fmt.Println(entries)
	}

	newEntry := UserEntry{
		Text: os.Args[1],
		Date: time.Now().Format("2006-01-02"),
	}

	entries = append(entries, newEntry)
	jsonData, err := json.Marshal(entries)
	if err != nil {
		fmt.Println("Error marshaling entries to JSON:", err)
		os.Exit(1)
	}

	if err = os.WriteFile(MEMORY, jsonData, 0644); err != nil {
		fmt.Println("Error updating memory JSON:", err)
		os.Exit(1)
	}

	fmt.Println("Ritual memory updated with entry", newEntry)
}

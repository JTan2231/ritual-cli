package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type UserEntry struct {
	Text string `json:"text"`
	Date string `json:"date"`
}

func printWrapped(s string, limit int) {
	for len(s) > limit {
		spaceIndex := strings.LastIndex(s[:limit], " ")
		if spaceIndex == -1 {
			spaceIndex = limit
		}
		fmt.Println("    " + s[:spaceIndex])
		s = s[spaceIndex:]
		if len(s) > 0 && s[0] == ' ' {
			s = s[1:]
		}
	}
	fmt.Println("    " + s)
}

func list(memory []UserEntry) {
	for _, entry := range memory {
		fmt.Println(entry.Date)
		printWrapped(entry.Text, 80)
		fmt.Println()
	}
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

	if len(os.Args) < 2 {
		fmt.Println("usage: ritual [--list] \"your entry\"")
		os.Exit(1)
	}

	listFlag := flag.Bool("list", false, "List all entries")
	flag.Parse()

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
	}

	if len(os.Args) > 2 || !*listFlag {
		var text string
		for _, arg := range os.Args[1:] {
			if arg != "--list" {
				text = arg
				break
			}
		}

		newEntry := UserEntry{
			Text: text,
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

		fmt.Println("Adding entry:", newEntry.Text)
	}

	if *listFlag {
		list(entries)
	}
}

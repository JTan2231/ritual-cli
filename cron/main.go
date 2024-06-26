package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type UserEntry struct {
	Text string `json:"text"`
	Date string `json:"date"`
}

func errorCheck(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
		os.Exit(1)
	}
}

func main() {
	var ROOT string

	dir, err := os.UserHomeDir()
	errorCheck("Error retrieving home directory: ", err)

	ROOT = dir + "/.ritual/"

	logPath := ROOT + "logs/" + time.Now().Format("2006-01-02") + ".log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	errorCheck("Error opening log file: ", err)

	defer file.Close()

	log.SetOutput(file)

	if len(os.Args) != 2 {
		log.Fatal("Error: Invalid number of arguments. Usage: ./cron SECRET")
		os.Exit(1)
	}

	authToken := os.Args[1]
	MEMORY := ROOT + "entries.json"

	var entryData []byte
	var userEntries []UserEntry
	if _, err := os.Stat(MEMORY); err == nil {
		entryData, err = os.ReadFile(MEMORY)
		errorCheck("Error reading file: ", err)

		err = json.Unmarshal(entryData, &userEntries)
		errorCheck("", err)
	}

	// get only those entries whose dates are within the past 7 days
	var recentEntries []UserEntry
	for _, entry := range userEntries {
		entryDate, err := time.Parse("2006-01-02", entry.Date)
		errorCheck("", err)
		if time.Since(entryDate).Hours() < 168 {
			recentEntries = append(recentEntries, entry)
		}
	}

	log.Println(recentEntries)

	postBody, err := json.Marshal(struct {
		Entries []UserEntry `json:"entries"`
	}{Entries: recentEntries})
	errorCheck("", err)

	req, err := http.NewRequest("POST", "https://ritual-api-production.up.railway.app/cli-newsletters", bytes.NewBuffer(postBody))
	errorCheck("", err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	client := &http.Client{}

	response, err := client.Do(req)
	errorCheck("", err)

	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	errorCheck("", err)

	log.Println(string(responseBody))
}

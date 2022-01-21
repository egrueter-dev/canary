package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LogProcessStart(event ProcessStartEvent, filename string) {
	logFile := LogFile{}

	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &logFile)

	// Append new data to Process starts
	logFile.ProcessStarts = append(logFile.ProcessStarts, event)

	marshalledJsonFile, _ := json.MarshalIndent(logFile, "", " ")
	_ = ioutil.WriteFile(filename, marshalledJsonFile, 0644)
}

func LogFileChange(event FileChangeEvent, filename string) {
	logFile := LogFile{}

	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &logFile)

	// Append new data to file changes
	logFile.FileChangeEvents = append(logFile.FileChangeEvents, event)

	marshalledJsonFile, _ := json.MarshalIndent(logFile, "", " ")
	_ = ioutil.WriteFile(filename, marshalledJsonFile, 0644)
}

func LogNetworkRequest(event NetworkRequestEvent, filename string) {
	logFile := LogFile{}

	jsonFile, _ := os.Open(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &logFile)

	// Append new data to network requests
	logFile.NetworkRequestEvents = append(logFile.NetworkRequestEvents, event)

	marshalledJsonFile, _ := json.MarshalIndent(logFile, "", " ")
	_ = ioutil.WriteFile(filename, marshalledJsonFile, 0644)
}

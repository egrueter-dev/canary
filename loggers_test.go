package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestLogNetworkRequest(t *testing.T) {
	fileName := "log_test.json"

	GenerateLogFile(fileName)

	data := NetworkRequestEvent{
		UserName:           "Rico",
		ProcessName:        "NetworkRequest",
		CommandLine:        "NetworkRequest",
		Protocol:           "HTTP",
		DestinationAddress: "server.com",
		DestinationPort:    "80",
		SourceAddress:      "localhost",
		SourcePort:         "8080",
		DataAmount:         10, // MB
		Timestamp:          time.Now(),
	}
	// TODO: you could check for the presence of the log in the JSON file..
	// This applies to other logging functions as well
	LogNetworkRequest(data, fileName)
	testCleanup()
}

// TODO: Should we check for presence of file before writing?
// Tests the user can write a new process
// to the test_log.json file
func TestLogProcessStart(t *testing.T) {
	fileName := "log_test.json"

	GenerateLogFile(fileName)

	data := ProcessStartEvent{
		UserName:    "Rico",
		ProcessName: "ProcessStarted",
		CommandLine: "--arg",
		Timestamp:   time.Now(),
	}

	LogProcessStart(data, fileName)

	jsonFile, err := os.Open(fileName)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	var logs LogFile

	json.Unmarshal(byteValue, &logs)

	// check if events are present
	if len(logs.ProcessStarts) == 1 {
	} else {
		t.Error("Error, processes not logged properly")
	}

	// TODO: Check if file is there to complete the test
	testCleanup()
}

func TestLoggingMultipleStartProcesses(t *testing.T) {
	fileName := "log_test.json"

	GenerateLogFile(fileName)

	data := ProcessStartEvent{
		UserName:    "Rico",
		ProcessName: "ProcessStarted",
		CommandLine: "--arg",
		Timestamp:   time.Now(),
	}

	LogProcessStart(data, fileName)

	data2 := ProcessStartEvent{
		UserName:    "Rico2",
		ProcessName: "ProcessStarted",
		CommandLine: "--arg",
		Timestamp:   time.Now(),
	}

	LogProcessStart(data2, fileName)

	jsonFile, err := os.Open(fileName)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	var logs LogFile

	json.Unmarshal(byteValue, &logs)

	// check if events are present
	if len(logs.ProcessStarts) == 2 {
	} else {
		t.Error("Error, processes not logged properly", err)
	}

	testCleanup()
}

func TestLogFileChange(t *testing.T) {
	fileName := "log_test.json"

	GenerateLogFile(fileName)

	data := FileChangeEvent{
		UserName:    "Rico2",
		FilePath:    "users/egrueter/exec",
		Descriptor:  "create",
		ProcessName: "FileCreated",
		CommandLine: "--create",
		ProcessId:   123,
		Timestamp:   time.Now(),
	}

	LogFileChange(data, fileName)

	data2 := FileChangeEvent{
		UserName:    "Rico3",
		FilePath:    "users/egrueter/exec",
		Descriptor:  "delete",
		ProcessName: "FileCreated",
		CommandLine: "--create",
		ProcessId:   123,
		Timestamp:   time.Now(),
	}

	LogFileChange(data, fileName)

	jsonFile, err := os.Open(fileName)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	var logs LogFile

	json.Unmarshal(byteValue, &logs)

	// check if events are present
	if len(logs.FileChangeEvents) == 2 {
	} else {
		t.Error("Error, events not logged properly")
	}

	LogFileChange(data2, fileName)
	testCleanup()
}

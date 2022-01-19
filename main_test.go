package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestGenerateLogFile(t *testing.T) {
	GenerateLogFile("log_test")

	if _, err := os.Stat("log_test.json"); err == nil {
	} else {
		t.Error("Error, unable to create logfile", err)
	}

	// Cleanup todo: abstract to helper method
	e := os.Remove("log_test.json")

	if e != nil {
		log.Fatal(e)
	}
}

// TODO: Should we check for presence of file before writing?

// Tests the user can write a new process
// to the log.json file
func TestLogProcessStart(t *testing.T) {
	GenerateLogFile("log_test")

	data := ProcessStartEvent{
		UserName:    "Rico",
		ProcessName: "ProcessStarted",
		CommandLine: "--arg",
		Timestamp:   time.Now(),
	}

	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile("log_test.json", file, 0644)

	// test conditions

	testCleanup()
}

func testCleanup() {
	e := os.Remove("log_test.json")

	if e != nil {
		log.Fatal(e)
	}

}

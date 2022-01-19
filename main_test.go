package main

import (
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

	testCleanup()
}

// TODO: Should we check for presence of file before writing?

// Tests the user can write a new process
// to the test_log.json file
func TestLogProcessStart(t *testing.T) {
	fileName := "log_test"

	GenerateLogFile(fileName)

	data := ProcessStartEvent{
		UserName:    "Rico",
		ProcessName: "ProcessStarted",
		CommandLine: "--arg",
		Timestamp:   time.Now(),
	}

	LogProcessStart(data, fileName)

	// 	testCleanup()
}

// func TestLoggingMultipleProcesses(t *testing.T) {
// 	GenerateLogFile("log_test")

// 	data := ProcessStartEvent{
// 		UserName:    "Rico",
// 		ProcessName: "ProcessStarted",
// 		CommandLine: "--arg",
// 		Timestamp:   time.Now(),
// 	}

// 	data2 := ProcessStartEvent{
// 		UserName:    "Rico2",
// 		ProcessName: "ProcessStarted2",
// 		CommandLine: "--arg",
// 		ProcessId:   0,
// 		Timestamp:   time.Now(),
// 	}
// 	// testCleanup()
// }

// Clean up log_test.json file
// after every test run
func testCleanup() {
	e := os.Remove("log_test.json")

	if e != nil {
		log.Fatal(e)
	}
}

// func Count(r io.Reader) (int, error) {
// 	dec := json.NewDecoder(r)

// 	count := 0

// 	for {
// 		t, err := dec.Token()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return -1, err
// 		}
// 		switch t {
// 		case json.Delim('{'):
// 			count++
// 		}
// 	}
// 	return count, nil
// }

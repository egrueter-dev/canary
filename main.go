package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
)

const LogFileName = "log.json"

// TODO:
// Handle wrong # of arguments per request

func main() {
	// Pull command Line Arguments
	osArgs := os.Args
	firstArg := osArgs[1]

	switch firstArg {
	case "-list":
		pterm.DefaultTable.WithHasHeader().WithData(pterm.TableData{
			{"Command", "Parameters", "Description"},
			{"-list", "", "List available commands"},
			{"-setup", "", "Generate Logfile"},
			{"-create", "[filepath]", "Create specific file"},
			{"-delete", "[path]", "Delete specific file"},
			{"-send-data", "[destination]", "Send log data to remote server"},
			{"-start-process", "[filepath, args]", "Execute binary"},
			{"-modify", "[filepath, text]", "Modify (add text) to a file"},
		}).Render()
	case "-setup":
		pterm.Success.Println("Generated Log File")
		GenerateLogFile(LogFileName)

		pterm.Success.Println("Generated example.txt File")
		GenerateExampleFiles()
	case "-start-process":
		path := osArgs[2]

		data := ProcessStartEvent{
			UserName:    fetchUserName(),
			ProcessName: "StartProcess",
			ProcessId:   os.Getpid(),
			CommandLine: "-start-process",
			Timestamp:   time.Now(),
		}

		LogProcessStart(data, LogFileName)
		ProcessStart(path, osArgs[2:])
	case "-create":
		filePath := osArgs[2]

		data := FileChangeEvent{
			UserName:    fetchUserName(),
			ProcessName: "FileCreated",
			ProcessId:   os.Getpid(),
			CommandLine: "--create",
			FilePath:    filePath,
			Descriptor:  "create",
			Timestamp:   time.Now(),
		}

		LogFileChange(data, LogFileName)
		CreateFile(filePath)
	case "-delete":
		filePath := osArgs[1]

		data := FileChangeEvent{
			UserName:    fetchUserName(),
			ProcessName: "FileDeleted",
			ProcessId:   os.Getpid(),
			CommandLine: "-delete",
			FilePath:    filePath,
			Descriptor:  "delete",
			Timestamp:   time.Now(),
		}

		LogFileChange(data, LogFileName)
		DeleteFile((filePath))
	case "-send-data":
		destination := "https://private-anon-6f9facff1e-restapi3.apiary-mock.com/notes"

		// This should be updated after the request is
		// Actually made

		NetworkRequest(destination)
	case "-modify":
		filePath := osArgs[1]
		text := osArgs[2]

		data := FileChangeEvent{
			UserName:    fetchUserName(),
			ProcessName: "FileModified", // TODO: Change process name
			ProcessId:   os.Getpid(),
			CommandLine: "-modify",
			FilePath:    filePath,
			Descriptor:  "modify",
			Timestamp:   time.Now(),
		}

		LogFileChange(data, LogFileName)
		ModifyFile(filePath, text)
	}
}

func GenerateExampleFiles() {
	d1 := []byte("hello ")
	err := os.WriteFile("example.txt", d1, 0644)

	if err != nil {
		panic(err)
	}
}

func GenerateLogFile(file string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	// Create empty struct
	processStart := make([]ProcessStartEvent, 0)
	fileChange := make([]FileChangeEvent, 0)
	networkRequest := make([]NetworkRequestEvent, 0)

	data := LogFile{
		processStart,
		fileChange,
		networkRequest,
	}

	jsonFile, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(file, jsonFile, 0644)

	f.Close()
}

func ProcessStart(path string, args []string) {
	output, err := exec.Command(path, args...).Output()

	if err != nil {
		panic(err)
	}

	fmt.Println("Process Output:")
	fmt.Println(string(output))
}

func CreateFile(path string) {
	// Create directories for path to file with permissions
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		panic(err)
	}
	os.Create(path)
}

func DeleteFile(path string) {
	if err := os.RemoveAll(filepath.Dir(path)); err != nil {
		panic(err)
	}
}

// Modify files by adding text. works for .txt files only
func ModifyFile(path string, text string) {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend,
	)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	defer file.Close()

	_, err = file.WriteString("world")

	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

// Send Log Data in Network Request
func NetworkRequest(url string) {
	jsonFile, err := os.Open(LogFileName)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	// Sent byte value of JSON
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteValue))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "none")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data := NetworkRequestEvent{
		UserName:           fetchUserName(),
		ProcessName:        "NetworkRequest",
		CommandLine:        "-send-data", // executable name and arguments together
		Protocol:           "HTTP",
		DestinationAddress: getRemoteIP(url),
		DestinationPort:    "?",
		SourceAddress:      getLocalIP(),
		SourcePort:         "8080",                 // ?>
		DataAmount:         getFileSize(*jsonFile), // get the size of the JSON file
		Timestamp:          time.Now(),
	}
	LogNetworkRequest(data, LogFileName)
}

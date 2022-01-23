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
	"strings"
	"time"

	"github.com/pterm/pterm"
)

const LogFileName = "log.json"

// TODO:
// Compile final binaries
func main() {
	// Pull command Line Arguments
	osArgs := os.Args
	firstArg := osArgs[1]
	commandLine := strings.Join(osArgs[:], ",")

	switch firstArg {
	case "-list":
		pterm.DefaultTable.WithHasHeader().WithData(pterm.TableData{
			{"Command", "Parameters", "Description"},
			{"-list", "", "List available commands"},
			{"-setup", "", "Generate Logfile"},
			{"-start-process", "[filepath, args]", "Execute binary"},
			{"-create", "[filepath]", "Create specific file"},
			{"-delete", "[path]", "Delete specific file"},
			{"-send-data", "[destination]", "Send log data to remote server"},
			{"-modify", "[filepath, text]", "Modify (add text) to a file"},
		}).Render()
	case "-setup":
		pterm.Success.Println("Generated Log File")
		GenerateLogFile(LogFileName)

		pterm.Success.Println("Generated example.txt File")
		GenerateExampleFiles()
	case "-start-process":
		checkArgumentPresent(os.Args, 3, "[filepath, args]")

		path := osArgs[2]

		data := ProcessStartEvent{
			UserName:    fetchUserName(),
			ProcessName: osArgs[0],
			ProcessId:   os.Getpid(),
			CommandLine: commandLine,
			Timestamp:   time.Now(),
		}

		LogProcessStart(data, LogFileName)
		ProcessStart(path, osArgs[2:])
		pterm.Success.Printf("Process %s executed", path)
	case "-create":
		checkArgumentPresent(os.Args, 3, "[filepath]")

		filePath := osArgs[2]

		data := FileChangeEvent{
			UserName:    fetchUserName(),
			ProcessName: osArgs[0],
			ProcessId:   os.Getpid(),
			CommandLine: commandLine,
			FilePath:    filePath,
			Descriptor:  "create",
			Timestamp:   time.Now(),
		}

		LogFileChange(data, LogFileName)
		CreateFile(filePath)

		pterm.Success.Printf("Created %s\n", filePath)
	case "-delete":
		checkArgumentPresent(os.Args, 3, "[filepath]")

		filePath := osArgs[3]

		data := FileChangeEvent{
			UserName:    fetchUserName(),
			ProcessName: osArgs[0],
			ProcessId:   os.Getpid(),
			CommandLine: commandLine,
			FilePath:    filePath,
			Descriptor:  "delete",
			Timestamp:   time.Now(),
		}

		LogFileChange(data, LogFileName)
		DeleteFile((filePath))

		pterm.Success.Printf("Deleted %s\n", filePath)
	case "-send-data":
		destination := osArgs[2]
		processName := osArgs[0]
		NetworkRequest(commandLine, processName, destination)
		pterm.Success.Printf("Network Request Complete\n")
	case "-modify":
		filePath := osArgs[2]
		text := osArgs[3]

		data := FileChangeEvent{
			UserName:    fetchUserName(),
			ProcessName: osArgs[0],
			ProcessId:   os.Getpid(),
			CommandLine: commandLine,
			FilePath:    filePath,
			Descriptor:  "modify",
			Timestamp:   time.Now(),
		}

		LogFileChange(data, LogFileName)
		ModifyFile(filePath, text)
		pterm.Success.Printf("Modified %s successfully\n", filePath)
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
	// Remove just file
	if filepath.Dir(path) == "." {
		os.RemoveAll(path)
		return
	}

	// Remove file & directory
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

	_, err = file.WriteString(text)

	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

// Send Log Data in Network Request
func NetworkRequest(commandLine string, processName string, url string) {
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

	// Animate Request
	introSpinner, _ := pterm.DefaultSpinner.WithRemoveWhenDone(false).Start("Sending Network Request")
	introSpinner.Start()

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	ip, port := getLocalIP()

	data := NetworkRequestEvent{
		UserName:           fetchUserName(),
		ProcessName:        processName,
		CommandLine:        commandLine,
		Protocol:           "HTTP",
		DestinationAddress: getRemoteIP(url),
		DestinationPort:    getRemotePort(url),
		SourceAddress:      ip,
		SourcePort:         port,
		DataAmount:         getFileSize(*jsonFile),
		Timestamp:          time.Now(),
	}

	LogNetworkRequest(data, LogFileName)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// Commands
// go run canary.exe
//  -- list ( List all available commands )
//  -- setup ( First step. set up log file )
//  -- start-process (executable path & args)
//  -- modify
//  -- delete [filepath] ( supply filepath )
//  -- transmit (transmit the )

func main() {
	// Pull command Line Arguments
	argsWithoutProg := os.Args[1:]

	firstArg := argsWithoutProg[0]

	// TODO: handle case where no args are present..
	switch firstArg {
	case "-list":
		const list = `
		Available Commands:
		-list - List available commands
		-setup - Generate logfile
		-create [filetype, path] - Create specific file
		-delete [path] - Delete specific file
		-send-data []
		-start-process [filepath, args] - Execute binary
		-modify [filepath, text]
		`

		fmt.Println(list)
	case "-setup":
		fmt.Println("Generated Log File")
		GenerateLogFile("log.json")

		fmt.Println("Generated Example File")
		GenerateExampleFiles()
	case "-start-process":
		user, err := user.Current()

		if err != nil {
			log.Fatalf(err.Error())
		}

		data := ProcessStartEvent{
			UserName:    user.Name,
			ProcessName: "StartProcess",
			ProcessId:   os.Getpid(),
			CommandLine: "-start-process",
			Timestamp:   time.Now(),
		}

		LogProcessStart(data, "log.json")
		// Actually start process here
	case "-create":
		// TODO: Internal Methods here?
		user, err := user.Current()

		if err != nil {
			log.Fatalf(err.Error())
		}

		filePath := argsWithoutProg[1]

		data2 := FileChangeEvent{
			UserName:    user.Name,
			ProcessName: "FileCreated",
			ProcessId:   os.Getpid(),
			CommandLine: "--create",
			FilePath:    filePath,
			Descriptor:  "create",
			Timestamp:   time.Now(),
		}

		// TODO: LOG.json is a constant value in produciton, can
		// make it a constant
		LogFileChange(data2, "log.json")
		CreateFile(filePath)
	case "-delete":
		fmt.Println("Args:")
		fmt.Println(argsWithoutProg)

		user, err := user.Current()

		if err != nil {
			log.Fatalf(err.Error())
		}

		filePath := argsWithoutProg[1]

		data2 := FileChangeEvent{
			UserName:    user.Name,
			ProcessName: "FileDeleted",
			ProcessId:   os.Getpid(),
			CommandLine: "-delete",
			FilePath:    filePath,
			Descriptor:  "delete",
			Timestamp:   time.Now(),
		}

		LogFileChange(data2, "log.json")
		DeleteFile((filePath))
	case "":

	}
}

func GenerateExampleFiles() {
	d1 := []byte("hello\ngo\n")
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

	// Append new data to file changes
	logFile.NetworkRequestEvents = append(logFile.NetworkRequestEvents, event)

	marshalledJsonFile, _ := json.MarshalIndent(logFile, "", " ")
	_ = ioutil.WriteFile(filename, marshalledJsonFile, 0644)
}

// TODO: Should this be Run Process
func ProcessStart(path string, arguments []string) {
	cmd := exec.Command(
		path, arguments...,
	)

	fmt.Println("path:")
	fmt.Println(path)

	fmt.Println("strings ")

	fmt.Println(
		strings.Join(arguments[:], " "),
	)

	fmt.Println(
		cmd.Start(),
	) // and wait

	log.Printf("Waiting for command to finish...")
	log.Printf("Process id is %v", cmd.Process.Pid)
	cmd.Wait()
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

//// CORE FUNCTIONALITY

/// LOGGING

// ● Network connection and data transmission
//      Timestamp of activity
//      Username that started the process that initiated the network activity
//      Destination address and port
//      Source address and port
//      Amount of data sent
//      Protocol of data sent
//      Process name
//      Process command line
//      Process ID

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

// fmt.Println(`
//     usage:
//           dig [@local-server] host [options]
//     options:
//           +trace
//     Example:
//           dig google.com
//           dig @8.8.8.8 yahoo.com
//           dig google.com +trace
//           dig google.com MX
// 	`)

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
		-send-data [destination]
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
	case "-send-data":

		destination := "https://private-anon-6f9facff1e-restapi3.apiary-mock.com/notes"

		// This should be updated after the request is
		// Actually made

		NetworkRequest(destination)
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

// Send Log Data in Network Request
func NetworkRequest(url string) {
	jsonFile, err := os.Open("log.json")
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

	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	data := NetworkRequestEvent{
		UserName:           user.Name,
		ProcessName:        "NetworkRequest",
		CommandLine:        "-send-data",
		Protocol:           "HTTP",
		DestinationAddress: getRemoteIP(url),
		DestinationPort:    "?",
		SourceAddress:      getLocalIP(),
		SourcePort:         "8080",                 // ?>
		DataAmount:         getFileSize(*jsonFile), // get the size of the JSON file
		Timestamp:          time.Now(),
	}
	LogNetworkRequest(data, "log.json")
}

// Gets the local IP request was made from
func getLocalIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).String()
}

// TODO: get port
// https://github.com/golang/go/issues/16142
func getRemoteIP(url string) string {
	ips, _ := net.LookupIP(url)
	// Maybe try
	// conn.RemoteAddr
	// port, _ := net.LookupPort("tcp", "google.com")
	// fmt.Println("port", port)

	fmt.Println("IPs: ", ips[1])
	return ips[1].String()
}

func getFileSize(file os.File) int64 {
	stat, err := file.Stat()

	fileSize := stat.Size()

	fileKb := (fileSize / 1024)

	if err != nil {
		panic(err)
	}

	return fileKb
}

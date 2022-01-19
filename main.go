package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Commands
// go run canary.exe
//  -- list ( List all available commands )
//  -- setup ( First step. set up log file )
//  -- start-process (executable path & args)
//  -- modify
//  -- delete [filepath] ( supply filepath )
//  -- transmit (transmit the )

// TODO: Testing should not create another log.json
func main() {
	// Start by generating the empty log file
	GenerateLogFile("logfile")

	// Log Process Start

	// Setup Flag
	// GenerateJsonFile

	// Log Process start
	// log process ID
	// Log username

	// user, err := user.Current()

	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }

	// username := user.Username
	// fmt.Printf("Username: %s\n", username)

	// fmt.Println(os.Getpid())
}

func GenerateLogFile(filename string) {
	// If the file doesn't exist, create it, or append to the file
	file := filename + ".json"

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}

	// Create empty struct
	f1 := make([]ProcessStartEvent, 0)

	data := LogFile{
		f1,
	}

	jsonFile, _ := json.MarshalIndent(data, "", " ")

	_ = ioutil.WriteFile(file, jsonFile, 0644)

	f.Close()
}

func LogProcessStart(event ProcessStartEvent, filename string) {
	// todo: make process start even external

}

//// CORE FUNCTIONALITY

/// Start Process
//  startProcess(path_to_file, args)

/// createFile(location)

/// modifyFile() // look up path to modify

/// deleteFile() // path to delete

// transmitLogs() // Establish a network connection and transmit data

/// LOGGING

// * Process start
//      Timestamp of start time
//      Username that started the process
//  	Process name
//  	Process command line
//  	Process ID
// ● File creation, modification, deletion
//    	Timestamp of activity
//    	Full path to the file
//    	Activity descriptor - e.g. create, modified, delete
//    	Username that started the process that created/modified/deleted the file
//    	Process name that created/modified/deleted the file
//    	Process command line
//    	Process ID
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

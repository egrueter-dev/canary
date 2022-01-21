package main

import "time"

// TODO: finish testing for types
type ProcessStartEvent struct {
	UserName, ProcessName, CommandLine string
	ProcessId                          int
	Timestamp                          time.Time
}

type FileChangeEvent struct {
	UserName, FilePath, Descriptor string
	ProcessName, CommandLine       string
	ProcessId                      int
	Timestamp                      time.Time
}

type NetworkRequestEvent struct {
	UserName                            string
	ProcessName, CommandLine            string
	Protocol                            string
	DestinationAddress, DestinationPort string
	SourceAddress, SourcePort           string
	ProcessId                           int
	DataAmount                          int64
	Timestamp                           time.Time
}

type LogFile struct {
	ProcessStarts        []ProcessStartEvent
	FileChangeEvents     []FileChangeEvent
	NetworkRequestEvents []NetworkRequestEvent
}

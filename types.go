package main

import "time"

type ProcessStartEvent struct {
	UserName, ProcessName, CommandLine string
	ProcessId                          int
	Timestamp                          time.Time
}

type LogFile struct {
	ProcessStarts []ProcessStartEvent
}

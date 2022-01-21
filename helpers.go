package main

import (
	"fmt"
	"net"
	"os"
)

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

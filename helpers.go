package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"os"
	"os/user"
)

func UnmarshallFile(fileName string, logs *LogFile) {
	jsonFile, err := os.Open(fileName)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(byteValue, &logs)
}

func fetchUserName() string {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	return user.Name
}

// Gets the local IP request was made from
func getLocalIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).String()
}

// TODO: Get port
// func getLocalPort() string {
// 	conn, _ := net.Dial("udp", "8.8.8.8:80")
// 	defer conn.Close()
// 	// return conn.LocalAddr().(*net.UDPAddr).String()
// 	return conn.LocalAddr().Network()
// }

// TODO: get port
// https://github.com/golang/go/issues/16142
func getRemoteIP(url string) string {
	ips, _ := net.LookupIP("google.com")
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.To4().String()
		}
	}
	return "not found"
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

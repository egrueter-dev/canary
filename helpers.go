package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
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

func getLocalIP() (string, string) {
	// Check local IP address through Google DNS resolver
	// unsure why this needs to be UDP?
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	udp := conn.LocalAddr()

	udpAddr := udp.(*net.UDPAddr)
	addr := udpAddr.IP
	port := udpAddr.Port

	return addr.String(), strconv.Itoa(port)
}

func getRemotePort(remoteUrl string) string {
	u, err := url.Parse(remoteUrl)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(u.Hostname(), "www")
	truncUrl := parts[1][1:]

	conn, _ := net.Dial("tcp", truncUrl+":http")
	defer conn.Close()

	remoteAddr := conn.RemoteAddr()
	addr := remoteAddr.(*net.TCPAddr)

	return strconv.Itoa(addr.Port)
}

func getRemoteIP(remoteUrl string) string {
	u, err := url.Parse(remoteUrl)
	if err != nil {
		log.Fatal(err)
	}

	// https://www.google.com -> google.com
	parts := strings.Split(u.Hostname(), "www")
	truncUrl := parts[1][1:]

	ips, _ := net.LookupIP(truncUrl)

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

func checkArgumentPresent(osArgs []string, expected int, arguments string) {
	if len(osArgs) < expected {
		pterm.Error.Printf("Supply Arguments: %s\n", arguments)
		os.Exit(1)
	}
}

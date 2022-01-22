package main

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateLogFile(t *testing.T) {
	GenerateLogFile("log_test.json")

	if _, err := os.Stat("log_test.json"); err == nil {
	} else {
		t.Error("Error, unable to create logfile", err)
	}

	testCleanup()
}

func TestCreateExampleFiles(t *testing.T) {
	GenerateExampleFiles()

	if _, err := os.Stat("example.txt"); err == nil {
	} else {
		t.Error("Error, unable to create logfile", err)
	}

	testCleanup()
}

func TestCreateFile(t *testing.T) {
	filepath := "./test_path/file.json"

	CreateFile(filepath)

	if _, err := os.Stat("./test_path/file.json"); err == nil {
	} else {
		t.Error("Error, file not being created", err)
	}

	testCleanup()
}

func TestDeleteFile(t *testing.T) {
	filepath := "./test_path/file.json"

	// Create file so it can be deleted
	CreateFile(filepath)

	DeleteFile(filepath)

	_, err := os.Open(filepath)

	if errors.Is(err, os.ErrNotExist) {
	} else {
		t.Error("File was not successfully deleted")
	}

	testCleanup()
}

func TestModifyFile(t *testing.T) {
	GenerateExampleFiles()

	filepath := "./example.txt"
	text := "world"

	ModifyFile(filepath, text)

	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	if string(content) == "hello world" {
	} else {
		t.Error("File was not successfully deleted")
	}
}

func TestNetworkRequest(t *testing.T) {
	destination := "https://www.google.com"
	NetworkRequest(destination)
}

func TestProcessStart(t *testing.T) {
	args := []string{"--a", "--b", "--c"}
	path := "./example_executable"
	ProcessStart(path, args)
}

// Clean up log_test.json file
// after every test run
func testCleanup() {
	os.Remove("log_test.json")
	os.Remove("example.txt")
	os.RemoveAll("./test_path/")
}

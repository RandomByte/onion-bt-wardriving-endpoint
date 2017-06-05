package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleDataRequest(t *testing.T) {
	// data := []byte(`[{"Name":"Device 1", "Mac": "12:34:56:78:90:42", "Count": 15, "LastSeen": 1496695900}, {"Name":"//Device $('2 ", "Mac": "13:37:13:37:13:37", "Count": 11, "LastSeen": 1496695900}]`)

	device1 := device{Mac: "12:34:56:78:90:42", Name: "Device 1", Count: 15, LastSeen: 1496695900}
	device2 := device{Mac: "13:37:13:37:13:37", Name: "//Device $('2 ", Count: 11, LastSeen: 1496695900}

	devices := []device{device1, device2}
	data, err := json.Marshal(devices)
	if err != nil {
		t.Errorf("Error during json marshal in test %v", err)
	}
	req := httptest.NewRequest("POST", "http://example.com/data", bytes.NewBuffer(data))

	handleDataRequest(httptest.NewRecorder(), req)
	if len(deviceCollection) != 2 {
		t.Errorf("Expected deviceCollection to hane two entries, but got %v", len(deviceCollection))
	}
}

func TestWriteOutCollection(t *testing.T) {
	path := "test/results.csv"
	writeOutCollection(path)
	defer os.Remove(path)

	readSomething := false

	file, err := os.Open(path)
	if err != nil {
		t.Errorf("Can't open file in test %v", err)
	}
	reader := csv.NewReader(bufio.NewReader(file))
	for {
		_, err := reader.Read()
		if err == io.EOF && readSomething != true {
			t.Error("File is empty")
		}
		if err == io.EOF {
			break
		}
		readSomething = true
	}
}

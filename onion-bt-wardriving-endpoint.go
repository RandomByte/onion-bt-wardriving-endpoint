package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

type device struct {
	Name     string
	Count    int
	LastSeen int64
	Mac      string
}

var deviceCollection []device

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	http.HandleFunc("/data", handleDataRequest)
	http.HandleFunc("/done", handleDoneRequest)
	portStr := strconv.Itoa(*port)
	log.Printf("Listening on %s...\n", portStr)
	log.Fatal(http.ListenAndServe(":"+portStr, nil))
}

func handleDataRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		var devices []device
		err := decoder.Decode(&devices)
		if err != nil {
			panic(err)
		}
		collectDevices(devices)
	} else {
		log.Printf("Invalid method: %s", req.Method)
	}
}

func collectDevices(devices []device) {
	deviceCollection = append(deviceCollection, devices...)
}

func handleDoneRequest(w http.ResponseWriter, req *http.Request) {
	writeOutCollection("out/devices.csv")
}

func writeOutCollection(filename string) {
	if len(deviceCollection) == 0 {
		log.Println("Device collection empty")
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	t := reflect.Indirect(reflect.ValueOf(deviceCollection[0])).Type()
	numOfFields := t.NumField()
	header := make([]string, numOfFields)
	for i := 0; i < numOfFields; i++ {
		header[i] = t.Field(i).Name
	}

	body := make([][]string, len(deviceCollection))
	for j, entry := range deviceCollection {
		line := make([]string, numOfFields)
		v := reflect.ValueOf(entry)
		for k := 0; k < numOfFields; k++ {
			line[k] = v.Field(k).String()
		}
		body[j] = line
	}

	csvData := [][]string{header}
	csvData = append(csvData, body...)

	for _, value := range csvData {
		err := writer.Write(value)
		if err != nil {
			log.Panicf("Cannot write to file: %v", err)
		}
	}
}

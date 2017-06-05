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
		log.Println("Incoming data POST...")
		decoder := json.NewDecoder(req.Body)
		var devices []device
		err := decoder.Decode(&devices)
		if err != nil {
			panic(err)
		}
		collectDevices(devices)
	} else {
		log.Printf("Invalid method: %s\n", req.Method)
	}
}

func collectDevices(devices []device) {
	log.Printf("Collected %v devices\n", len(devices))
	deviceCollection = append(deviceCollection, devices...)
}

func handleDoneRequest(w http.ResponseWriter, req *http.Request) {
	log.Println("Incoming done signal...")
	writeOutCollection("out/devices.csv")
}

func writeOutCollection(path string) {
	if len(deviceCollection) == 0 {
		log.Println("Device collection empty")
		return
	}

	file, err := os.Create(path)
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
			field := v.Field(k)
			if field.Type().Name() == "string" {
				line[k] = field.String()
			} else {
				line[k] = strconv.FormatInt(field.Int(), 10)
			}
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
	log.Printf("CSV written to %s\n", path)
}

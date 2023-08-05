package main

import (
	"encoding/json"
	"log"
	"os"
)

var pathMap map[string][]Point

func setupPath() {
	pathMap = make(map[string][]Point)
	file, err := os.Open("station_details.json")
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	var stationDetailMap map[string]StationDetail
	err = json.NewDecoder(file).Decode(&stationDetailMap)
	if err != nil {
		log.Fatal("Error reading from file")
	}
}

func getPathMap() *map[string][]Point {
	return &pathMap
}

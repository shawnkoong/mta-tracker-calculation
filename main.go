package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

func main() {

	setupPath()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	doTasks()
	for {
		select {
		case <-ticker.C:
			doTasks()
		}
	}
}

func fetchData() *MtaResponse {
	response, err := http.Get("https://www.goodservice.io/api/routes/?detailed=1")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer response.Body.Close()

	var data MtaResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Printf("Error Reading Response: %s\n", err)
	}
	log.Println("fetched data")
	return &data
}

func printMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Allocated Memory: %d bytes\n", m.Alloc)
	fmt.Printf("Total Memory Allocated (including freed): %d bytes\n", m.TotalAlloc)
	fmt.Printf("Memory Obtained from the OS: %d bytes\n", m.Sys)
	fmt.Printf("Number of Garbage Collection Cycles: %d\n", m.NumGC)
}

func doTasks() {
	data := fetchData()
	go calculateArrivalTimes(data)
	go calculateTrainPositions(data)
}

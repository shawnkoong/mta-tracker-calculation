package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {

	setupPath()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data := fetchData()
			go calculateArrivalTimes(data)
			go calculateTrainPositions(data)
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

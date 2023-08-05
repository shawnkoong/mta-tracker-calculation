package main

import (
	"log"
	"time"
)

func calculateArrivalTimes(data *MtaResponse) {
	log.Printf("Starting arrival time calculation %s\n", time.Now())
	log.Printf("Finished arrival time calculation %s\n", time.Now())
}

func calculateTrainPositions(data *MtaResponse) {
	log.Printf("Starting train position calculation %s\n", time.Now())
	log.Printf("Finished train position calculation %s\n", time.Now())
}

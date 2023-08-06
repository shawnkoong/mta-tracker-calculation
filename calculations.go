package main

import (
	"log"
	"time"
)

const ArrivalTimeKey string = "time"
const TrainPositionKey string = "train"

func calculateArrivalTimes(data *MtaResponse) {
	log.Printf("Starting arrival time calculation %s\n", time.Now())
	routes := data.Routes
	allTimes := make(map[string]map[string]map[string][]ArrivalTime) // line to stopId to direction to times
	for line, route := range routes {
		if !route.Scheduled || route.Status == "No Service" {
			continue
		}
		stopTimes := make(map[string]map[string][]ArrivalTime)
		allTimes[line] = stopTimes
		trips := route.Trips
		for direction, tripArr := range trips {
			for _, trip := range tripArr {
				lastStopId := trip.LastStopMade
				if lastStopId == "" {
					continue
				}
				stops := trip.Stops
				for stopId, arriveTime := range stops {
					eta := int64(arriveTime) - time.Now().Unix()
					if stopTimes[stopId] == nil {
						// since there is no entry for this stop yet, create a new map with two empty ArrivalTime slices
						directionMap := make(map[string][]ArrivalTime)
						directionMap["north"] = make([]ArrivalTime, 0)
						directionMap["south"] = make([]ArrivalTime, 0)
						stopTimes[stopId] = directionMap
					}
					arrivalTime := ArrivalTime{Time: eta, IsDelayed: trip.IsDelayed, IsAssigned: trip.IsAssigned}
					stopTimes[stopId][direction] = append(stopTimes[stopId][direction], arrivalTime)
				}
			}
		}
	}
	client := getRedisClient()
	err := save(client, ArrivalTimeKey, allTimes)
	if err != nil {
		log.Printf("Error saving to redis: %s\n", err)
	}
	log.Printf("Finished arrival time calculation %s\n", time.Now())
}

func calculateTrainPositions(data *MtaResponse) {
	log.Printf("Starting train position calculation %s\n", time.Now())
	log.Printf("Finished train position calculation %s\n", time.Now())
}

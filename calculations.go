package main

import (
	"log"
	"sort"
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
		for direction, tripArr := range route.Trips {
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
	routes := data.Routes
	allPositions := make(map[string]map[string][]CoordinateBearing)
	for line, route := range routes {
		if !route.Scheduled || route.Status == "No Service" {
			continue
		}
		directionMap := make(map[string][]CoordinateBearing)
		for direction, tripArr := range route.Trips {
			trainPositions := make([]CoordinateBearing, 0)
			for _, trip := range tripArr {
				lastStopId := trip.LastStopMade
				if lastStopId == "" {
					continue
				}
				stops := trip.Stops
				nextStopId := findNextStopId(stops, lastStopId)
				if nextStopId == "" {
					continue
				}
				var pathName string
				if direction == "north" {
					pathName = lastStopId + "-" + nextStopId
				} else {
					pathName = nextStopId + "-" + lastStopId
				}
				path := getPath(pathName)
				if path.length() == 0 {
					continue
				}
				lastTimestamp := stops[lastStopId]
				nextTimestamp := stops[nextStopId]
				nowTimestamp := time.Now().Unix()
				progress := float64(nowTimestamp-lastTimestamp) / float64(nextTimestamp-lastTimestamp)
				coordBearing := path.getPointAtProgress(progress)
				trainPositions = append(trainPositions, coordBearing)
			}
			directionMap[direction] = trainPositions
		}
		allPositions[line] = directionMap
	}
	client := getRedisClient()
	err := save(client, TrainPositionKey, allPositions)
	if err != nil {
		log.Printf("Error saving to redis: %s\n", err)
	}
	log.Printf("Finished train position calculation %s\n", time.Now())
}

func findNextStopId(stops map[string]int64, lastStopId string) string {
	keys := make([]string, len(stops))
	i := 0
	for key := range stops {
		keys[i] = key
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return stops[keys[i]] < stops[keys[j]]
	})
	// pretty sure index can't equal -1 since the train's last stop cannot be the final stop, but just in case
	index := indexOf(lastStopId, keys)
	if index == -1 {
		return ""
	}
	return keys[index+1]
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}

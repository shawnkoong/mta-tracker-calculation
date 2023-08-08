package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const fileName string = "station_details.json"

var pathMap map[string]Path
var stationDetailMap map[string]StationDetail

func setupPath() {
	pathMap = make(map[string]Path)
	stationDetailMap = make(map[string]StationDetail)
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error opening file: %s", err)
	}
	err = json.NewDecoder(file).Decode(&stationDetailMap)
	if err != nil {
		log.Fatal("Error reading from file")
	}
	for stopId, stationDetail := range stationDetailMap {
		for nextStopId, coordinates := range stationDetail.North {
			pathName := stopId + "-" + nextStopId
			points := make([]Point, len(coordinates))
			// coordinates doesn't include current or next stops' coordinates
			for i, coordinate := range coordinates {
				points[i] = Point{Longitude: coordinate[0], Latitude: coordinate[1]}
			}
			path := Path{PathName: pathName, Points: points}
			// add current stop's coordinate
			path.appendToFront(Point{Longitude: stationDetail.Longitude, Latitude: stationDetail.Latitude})
			// add next stop's coordinate
			nextStationDetail := stationDetailMap[nextStopId]
			path.appendToEnd(Point{Longitude: nextStationDetail.Longitude, Latitude: nextStationDetail.Latitude})
			pathMap[pathName] = path
		}
	}
	log.Println("Finished building paths")
}

func getPathMap() *map[string]Path {
	return &pathMap
}

func getPath(pathName string) *Path {
	pm := getPathMap()
	if path, ok := (*pm)[pathName]; ok {
		return &path
	}
	nameSplit := strings.Split(pathName, "-")
	path := getPathRecursive(nameSplit[0], nameSplit[1], 0)
	if path != nil {
		(*pm)[pathName] = *path
	}
	return path
}

func getPathRecursive(start string, end string, step int) *Path {
	if step > 15 {
		return nil
	}
	pm := getPathMap()
	if path, ok := (*pm)[start+"-"+end]; ok {
		return &path
	}
	curr := stationDetailMap[start]
	if len(curr.North) == 0 {
		return nil
	}
	for next := range curr.North {
		nextPath := getPathRecursive(next, end, step+1)
		if nextPath == nil {
			continue
		}
		// this will exist as it is built during setup
		pathToNext := (*pm)[start+"-"+next]
		pathToNext.appendAll(nextPath)
		pathToNext.updateName(start + "-" + end)
		return &pathToNext
	}
	// if the for loop terminates, then none of the next stops led to end station
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type MtaResponse struct {
	Timestamp int64            `json:"timestamp"`
	Routes    map[string]Route `json:"routes"`
}

type Route struct {
	Id                     string                `json:"id"`
	Name                   string                `json:"name"`
	Color                  string                `json:"color"`
	Status                 string                `json:"status"`
	Visible                bool                  `json:"visible"`
	Scheduled              bool                  `json:"scheduled"`
	DirectionStatuses      map[string]string     `json:"direction_statuses"`
	DelaySummaries         map[string]string     `json:"delay_summaries"`
	IrregularitySummaries  map[string]string     `json:"service_irregularity_summaries"`
	ServiceChangeSummaries map[string][]string   `json:"service_change_summaries"`
	ActualRoutings         map[string][][]string `json:"actual_routings"`
	ScheduledRoutings      map[string][][]string `json:"scheduled_routings"`
	SlowSections           map[string][]Section  `json:"slow_sections"`
	LongHeadwaySections    map[string][]Section  `json:"long_headway_sections"`
	DelayedSections        map[string][]Section  `json:"delayed_sections"`
	Trips                  map[string][]Trip     `json:"trips"`
}

type Section struct {
	Begin               string  `json:"begin"`
	End                 string  `json:"end"`
	RuntimeDiff         float64 `json:"runtime_diff"`
	MaxActualHeadway    float64 `json:"max_actual_headway"`
	MaxScheduledHeadway float64 `json:"max_scheduled_headway"`
	DelayedTime         float64 `json:"delayed_time"`
}

type Trip struct {
	Id                   string             `json:"id"`
	Stops                map[string]float64 `json:"stops"`
	DelayedTime          float64            `json:"delayed_time"`
	ScheduledDiscrepancy float64            `json:"scheduled_discrepancy"`
	IsDelayed            bool               `json:"is_delayed"`
	IsAssigned           bool               `json:"is_assigned"`
	LastStopMade         string             `json:"last_stop_made"`
}

func main() {

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
	fmt.Println(data)
}

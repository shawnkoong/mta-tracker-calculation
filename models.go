package main

import (
	"errors"
	"math"
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

type Path struct {
	PathName string
	Points   []Point
}

func (p *Path) appendToFront(point Point) {
	p.Points = append([]Point{point}, p.Points...)
}

func (p *Path) appendToEnd(point Point) {
	p.Points = append(p.Points, point)
}

// appends all points into p.Points except the first point from path
func (p *Path) appendAll(path *Path) {
	points := path.Points[1:]
	p.Points = append(p.Points, points...)
}

func (p *Path) updateName(newName string) {
	p.PathName = newName
}

func (p *Path) length() int {
	return len(p.Points)
}

func (p *Path) getPointAtProgress(progress float64) (CoordinateBearing, error) {
	tDist := p.findTotalDistance()
	distTravelled := tDist * progress
	points := p.Points
	for i := 0; i < len(points)-1; i++ {
		distBwtStations := haversineDistance(points[i], points[i+1])
		if distBwtStations < distTravelled {
			distTravelled -= distBwtStations
		} else {
			bearing := calculateBearing(points[i], points[i+1])
			bearing = normalizeBearing(bearing)
			currPoint := pointWithProgress(points[i], distTravelled, bearing)
			return CoordinateBearing{Bearing: bearing * 180 / math.Pi, Latitude: currPoint.Latitude, Longitude: currPoint.Longitude}, nil
		}
	}
	return CoordinateBearing{}, errors.New("not able to find current position")
}

func (p *Path) findTotalDistance() float64 {
	points := p.Points
	if n := len(points); n > 1 {
		distance := 0.0
		for i := 0; i < len(points)-1; i++ {
			distance += haversineDistance(points[i], points[i+1])
		}
		return distance
	}
	return 0.0
}

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type StationDetail struct {
	Name      string                 `json:"name"`
	Longitude float64                `json:"longitude"`
	Latitude  float64                `json:"latitude"`
	North     map[string][][]float64 `json:"north"`
}

type ArrivalTime struct {
	Time       int64 `json:"time"`
	IsDelayed  bool  `json:"isDelayed"`
	IsAssigned bool  `json:"isAssigned"`
}

type CoordinateBearing struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Bearing   float64 `json:"bearing"` // in degrees from north, clock-wise
}

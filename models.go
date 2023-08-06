package main

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
	Id                   string           `json:"id"`
	Stops                map[string]int64 `json:"stops"`
	DelayedTime          float64          `json:"delayed_time"`
	ScheduledDiscrepancy float64          `json:"scheduled_discrepancy"`
	IsDelayed            bool             `json:"is_delayed"`
	IsAssigned           bool             `json:"is_assigned"`
	LastStopMade         string           `json:"last_stop_made"`
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

func (p *Path) getPointAtProgress(progress float64) CoordinateBearing {
	points := p.Points
	// TODO
}

type Point struct {
	Longitude float64
	Latitude  float64
}

type StationDetail struct {
	Name      string                 `json:"name"`
	Longitude float64                `json:"longitude"`
	Latitude  float64                `json:"latitude"`
	North     map[string][][]float64 `json:"north"`
}

type ArrivalTime struct {
	Time       int64
	IsDelayed  bool
	IsAssigned bool
}

type CoordinateBearing struct {
	Longitude float64
	Latitude  float64
	Bearing   float32
}

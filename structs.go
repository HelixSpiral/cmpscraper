package cmpscraper

import "time"

type CMP struct {
	LastUpdate   time.Time
	Total        string
	WithoutPower string
	Counties     map[string]Outage
	NoOutages    bool
}

type Outage struct {
	Total        string
	WithoutPower string
}

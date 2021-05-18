package types

import "time"

type GopsStat struct {
	Property  string
	Value     interface{}
	IntValue  int64
	ValueUnit string
	Label     string
}

type GopsStats struct {
	Stats          []GopsStat
	Started        time.Time
	Duration       time.Duration
	Success        bool
	ConnectionsQty int
}

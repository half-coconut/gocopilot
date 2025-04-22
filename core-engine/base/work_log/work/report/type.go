package report

import "time"

type Report struct {
	Total         int
	Rate          float64
	Throughput    float64
	TotalDuration time.Duration
	Min           time.Duration
	Mean          time.Duration
	Max           time.Duration
	P50           time.Duration
	P90           time.Duration
	P95           time.Duration
	P99           time.Duration
	Ratio         float64
}

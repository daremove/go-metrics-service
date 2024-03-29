package storage

import "errors"

var (
	ErrDataNotFound = errors.New("data is not found")
)

type GaugeMetric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type CounterMetric struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

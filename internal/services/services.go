package services

import "errors"

var (
	ErrMetricNotFound = errors.New("metric data is not found")
)

type MetricEntry struct {
	Name  string
	Value string
}

type MetricSaveParameters struct {
	MetricType  string
	MetricName  string
	MetricValue string
}

type MetricGetParameters struct {
	MetricType string
	MetricName string
}

package metrics

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/storage"
	"strconv"
)

type SaveParameters struct {
	MetricType  string
	MetricName  string
	MetricValue string
}

type Service interface {
	Save(parameters SaveParameters) error
}

type metrics struct {
	store storage.Storage
}

func New(store storage.Storage) Service {
	return &metrics{
		store,
	}
}

func (m *metrics) Save(parameters SaveParameters) error {
	switch parameters.MetricType {
	case "gauge":
		v, err := strconv.ParseFloat(parameters.MetricValue, 64)

		if err != nil {
			return err
		}

		if err := m.store.AddGauge(parameters.MetricName, v); err != nil {
			return err
		}
	case "counter":
		v, err := strconv.ParseInt(parameters.MetricValue, 10, 64)

		if err != nil {
			return err
		}

		if err := m.store.AddCounter(parameters.MetricName, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("metrict type %s isn't defined", parameters.MetricType)
	}

	return nil
}

func IsCounterMetricType(metricName string) bool {
	return metricName == "PollCount"
}

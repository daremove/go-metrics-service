package metrics

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/services"
	"github.com/daremove/go-metrics-service/internal/storage"
	"strconv"
)

type Metrics struct {
	storage Storage
}

type Storage interface {
	GetGaugeMetric(key string) (float64, bool)
	GetGaugeMetrics() []storage.GaugeMetric

	GetCounterMetric(key string) (int64, bool)
	GetCounterMetrics() []storage.CounterMetric

	AddGauge(key string, value float64) error
	AddCounter(key string, value int64) error
}

func New(storage Storage) *Metrics {
	return &Metrics{
		storage,
	}
}

func (m *Metrics) Save(parameters services.MetricSaveParameters) error {
	switch parameters.MetricType {
	case "gauge":
		v, err := strconv.ParseFloat(parameters.MetricValue, 64)

		if err != nil {
			return err
		}

		if err := m.storage.AddGauge(parameters.MetricName, v); err != nil {
			return err
		}
	case "counter":
		v, err := strconv.ParseInt(parameters.MetricValue, 10, 64)

		if err != nil {
			return err
		}

		if err := m.storage.AddCounter(parameters.MetricName, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("metrict type %s isn't defined", parameters.MetricType)
	}

	return nil
}

func (m *Metrics) Get(parameters services.MetricGetParameters) (string, bool) {
	switch parameters.MetricType {
	case "gauge":
		value, ok := m.storage.GetGaugeMetric(parameters.MetricName)

		return fmt.Sprintf("%g", value), ok
	case "counter":
		value, ok := m.storage.GetCounterMetric(parameters.MetricName)

		return fmt.Sprintf("%v", value), ok
	default:
		return "", false
	}
}

func (m *Metrics) GetAll() []services.MetricEntry {
	var result []services.MetricEntry

	for _, item := range m.storage.GetGaugeMetrics() {
		result = append(result, services.MetricEntry{Name: item.Name, Value: fmt.Sprintf("%g", item.Value)})
	}

	for _, item := range m.storage.GetCounterMetrics() {
		result = append(result, services.MetricEntry{Name: item.Name, Value: fmt.Sprintf("%v", item.Value)})
	}

	return result
}

func IsCounterMetricType(metricName string) bool {
	return metricName == "PollCount"
}

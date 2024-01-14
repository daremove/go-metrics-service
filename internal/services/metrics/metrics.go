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

type GetParameters struct {
	MetricType string
	MetricName string
}

type MetricItem struct {
	Name  string
	Value string
}

type Service interface {
	Save(parameters SaveParameters) error

	Get(parameters GetParameters) (string, bool)

	GetAll() []MetricItem
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

func (m *metrics) Get(parameters GetParameters) (string, bool) {
	switch parameters.MetricType {
	case "gauge":
		value, ok := m.store.GetGaugeMetric(parameters.MetricName)

		return fmt.Sprintf("%g", value), ok
	case "counter":
		value, ok := m.store.GetCounterMetric(parameters.MetricName)

		return fmt.Sprintf("%v", value), ok
	default:
		return "", false
	}
}

func (m *metrics) GetAll() []MetricItem {
	var result []MetricItem

	for _, item := range m.store.GetGaugeMetrics() {
		result = append(result, MetricItem{item.Name, fmt.Sprintf("%g", item.Value)})
	}

	for _, item := range m.store.GetCounterMetrics() {
		result = append(result, MetricItem{item.Name, fmt.Sprintf("%v", item.Value)})
	}

	return result
}

func IsCounterMetricType(metricName string) bool {
	return metricName == "PollCount"
}

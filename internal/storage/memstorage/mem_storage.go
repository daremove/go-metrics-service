package memstorage

import (
	"github.com/daremove/go-metrics-service/internal/storage"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (s *MemStorage) GetGaugeMetric(key string) (float64, bool) {
	value, ok := s.gauge[key]

	return value, ok
}

func (s *MemStorage) GetGaugeMetrics() []storage.GaugeMetric {
	var data []storage.GaugeMetric

	for key, value := range s.gauge {
		data = append(data, storage.GaugeMetric{Name: key, Value: value})
	}

	return data
}

func (s *MemStorage) GetCounterMetric(key string) (int64, bool) {
	value, ok := s.counter[key]

	return value, ok
}

func (s *MemStorage) GetCounterMetrics() []storage.CounterMetric {
	var data []storage.CounterMetric

	for key, value := range s.counter {
		data = append(data, storage.CounterMetric{Name: key, Value: value})
	}

	return data
}

func (s *MemStorage) AddGauge(key string, value float64) error {
	s.gauge[key] = value

	return nil
}

func (s *MemStorage) AddCounter(key string, value int64) error {
	s.counter[key] += value

	return nil
}

func New() *MemStorage {
	return &MemStorage{gauge: map[string]float64{}, counter: map[string]int64{}}
}

func NewWithPrefilledData(gauge map[string]float64, counter map[string]int64) *MemStorage {
	return &MemStorage{gauge, counter}
}

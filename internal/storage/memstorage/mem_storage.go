package memstorage

import (
	"context"

	"github.com/daremove/go-metrics-service/internal/storage"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (s *MemStorage) GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error) {
	value, ok := s.gauge[key]

	if !ok {
		return storage.GaugeMetric{}, storage.ErrDataNotFound
	}

	return storage.GaugeMetric{Name: key, Value: value}, nil
}

func (s *MemStorage) GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error) {
	var data []storage.GaugeMetric

	for key, value := range s.gauge {
		data = append(data, storage.GaugeMetric{Name: key, Value: value})
	}

	return data, nil
}

func (s *MemStorage) GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error) {
	value, ok := s.counter[key]

	if !ok {
		return storage.CounterMetric{}, storage.ErrDataNotFound
	}

	return storage.CounterMetric{Name: key, Value: value}, nil
}

func (s *MemStorage) GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error) {
	var data []storage.CounterMetric

	for key, value := range s.counter {
		data = append(data, storage.CounterMetric{Name: key, Value: value})
	}

	return data, nil
}

func (s *MemStorage) AddGaugeMetric(ctx context.Context, key string, value float64) error {
	s.gauge[key] = value

	return nil
}

func (s *MemStorage) AddCounterMetric(ctx context.Context, key string, value int64) error {
	s.counter[key] += value

	return nil
}

func (s *MemStorage) AddMetrics(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) error {
	for _, gaugeMetric := range gaugeMetrics {
		if err := s.AddGaugeMetric(ctx, gaugeMetric.Name, gaugeMetric.Value); err != nil {
			return err
		}
	}

	for _, counterMetric := range counterMetrics {
		if err := s.AddCounterMetric(ctx, counterMetric.Name, counterMetric.Value); err != nil {
			return err
		}
	}

	return nil
}

func New() *MemStorage {
	return &MemStorage{gauge: map[string]float64{}, counter: map[string]int64{}}
}

func NewWithPrefilledData(gauge map[string]float64, counter map[string]int64) *MemStorage {
	return &MemStorage{gauge, counter}
}

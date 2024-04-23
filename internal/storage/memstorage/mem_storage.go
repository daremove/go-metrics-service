// Пакет memstorage предоставляет реализацию хранилища в памяти для метрик.
package memstorage

import (
	"context"

	"github.com/daremove/go-metrics-service/internal/storage"
)

// MemStorage реализует интерфейс Storage, предоставляя операции с метриками, хранящимися в памяти.
type MemStorage struct {
	gauge   map[string]float64 // Хранение метрик типа gauge
	counter map[string]int64   // Хранение метрик типа counter
}

// GetGaugeMetric извлекает метрику типа gauge по ключу.
func (s *MemStorage) GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error) {
	value, ok := s.gauge[key]

	if !ok {
		return storage.GaugeMetric{}, storage.ErrDataNotFound
	}

	return storage.GaugeMetric{Name: key, Value: value}, nil
}

// GetGaugeMetrics возвращает все метрики типа gauge.
func (s *MemStorage) GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error) {
	var data []storage.GaugeMetric

	for key, value := range s.gauge {
		data = append(data, storage.GaugeMetric{Name: key, Value: value})
	}

	return data, nil
}

// GetCounterMetric извлекает метрику типа counter по ключу.
func (s *MemStorage) GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error) {
	value, ok := s.counter[key]

	if !ok {
		return storage.CounterMetric{}, storage.ErrDataNotFound
	}

	return storage.CounterMetric{Name: key, Value: value}, nil
}

// GetCounterMetrics возвращает все метрики типа counter.
func (s *MemStorage) GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error) {
	var data []storage.CounterMetric

	for key, value := range s.counter {
		data = append(data, storage.CounterMetric{Name: key, Value: value})
	}

	return data, nil
}

// AddGaugeMetric добавляет или обновляет метрику типа gauge.
func (s *MemStorage) AddGaugeMetric(ctx context.Context, key string, value float64) error {
	s.gauge[key] = value

	return nil
}

// AddCounterMetric добавляет или инкрементирует метрику типа counter.
func (s *MemStorage) AddCounterMetric(ctx context.Context, key string, value int64) error {
	s.counter[key] += value

	return nil
}

// AddMetrics добавляет набор метрик типа gauge и counter.
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

// New создает новый экземпляр MemStorage с пустыми картами для метрик.
func New() *MemStorage {
	return &MemStorage{gauge: map[string]float64{}, counter: map[string]int64{}}
}

// NewWithPrefilledData создает новый экземпляр MemStorage с предварительно заполненными данными.
func NewWithPrefilledData(gauge map[string]float64, counter map[string]int64) *MemStorage {
	return &MemStorage{gauge, counter}
}

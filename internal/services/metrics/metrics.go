// Пакет metrics предоставляет функции для работы с метриками, включая сохранение, получение и обновление данных.
package metrics

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services"
	"github.com/daremove/go-metrics-service/internal/storage"
)

// Metrics предоставляет методы для управления метриками через определенное хранилище.
type Metrics struct {
	storage Storage
}

// Storage определяет интерфейс для механизмов хранения, используемых системой метрик.
type Storage interface {
	GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error)
	GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error)
	GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error)
	GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error)
	AddGaugeMetric(ctx context.Context, key string, value float64) error
	AddCounterMetric(ctx context.Context, key string, value int64) error
	AddMetrics(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) error
}

// New создает новый экземпляр Metrics.
func New(storage Storage) *Metrics {
	return &Metrics{
		storage,
	}
}

// Save сохраняет одиночную метрику на основе предоставленных параметров.
func (m *Metrics) Save(ctx context.Context, parameters services.MetricSaveParameters) error {
	switch parameters.MetricType {
	case "gauge":
		v, err := strconv.ParseFloat(parameters.MetricValue, 64)

		if err != nil {
			return err
		}

		if err := m.storage.AddGaugeMetric(ctx, parameters.MetricName, v); err != nil {
			return err
		}
	case "counter":
		v, err := strconv.ParseInt(parameters.MetricValue, 10, 64)

		if err != nil {
			return err
		}

		if err := m.storage.AddCounterMetric(ctx, parameters.MetricName, v); err != nil {
			return err
		}
	default:
		return fmt.Errorf("metrict type %s isn't defined", parameters.MetricType)
	}

	return nil
}

// SaveModel сохраняет модель метрики.
func (m *Metrics) SaveModel(ctx context.Context, parameters models.Metrics) error {
	switch parameters.MType {
	case "gauge":
		if err := m.storage.AddGaugeMetric(ctx, parameters.ID, *parameters.Value); err != nil {
			return err
		}
	case "counter":
		if err := m.storage.AddCounterMetric(ctx, parameters.ID, *parameters.Delta); err != nil {
			return err
		}
	default:
		return fmt.Errorf("metrict type %s isn't defined", parameters.MType)
	}

	return nil
}

// SaveModels сохраняет массив метрик.
func (m *Metrics) SaveModels(ctx context.Context, parameters []models.Metrics) error {
	var gaugeMetrics []storage.GaugeMetric
	var counterMetrics []storage.CounterMetric

	for _, parameter := range parameters {
		switch parameter.MType {
		case "gauge":
			gaugeMetrics = append(gaugeMetrics, storage.GaugeMetric{Name: parameter.ID, Value: *parameter.Value})
		case "counter":
			counterMetrics = append(counterMetrics, storage.CounterMetric{Name: parameter.ID, Value: *parameter.Delta})
		default:
			return fmt.Errorf("metrict type %s isn't defined", parameter.MType)
		}
	}

	if err := m.storage.AddMetrics(ctx, gaugeMetrics, counterMetrics); err != nil {
		return err
	}

	return nil
}

// Get возвращает значение метрики по указанным параметрам.
func (m *Metrics) Get(ctx context.Context, parameters services.MetricGetParameters) (string, error) {
	switch parameters.MetricType {
	case "gauge":
		value, err := m.storage.GetGaugeMetric(ctx, parameters.MetricName)

		if err != nil {
			if errors.Is(err, storage.ErrDataNotFound) {
				return "", services.ErrMetricNotFound
			}

			return "", err
		}

		return fmt.Sprintf("%g", value.Value), nil
	case "counter":
		value, err := m.storage.GetCounterMetric(ctx, parameters.MetricName)

		if err != nil {
			if errors.Is(err, storage.ErrDataNotFound) {
				return "", services.ErrMetricNotFound
			}

			return "", err
		}

		return fmt.Sprintf("%v", value.Value), err
	default:
		return "", services.ErrMetricNotFound
	}
}

// GetModel возвращает полную модель метрики.
func (m *Metrics) GetModel(ctx context.Context, parameters models.Metrics) (models.Metrics, error) {
	switch parameters.MType {
	case "gauge":
		value, err := m.storage.GetGaugeMetric(ctx, parameters.ID)

		if err != nil {
			if errors.Is(err, storage.ErrDataNotFound) {
				return models.Metrics{}, services.ErrMetricNotFound
			}

			return models.Metrics{}, err
		}

		return models.Metrics{
			ID:    parameters.ID,
			MType: parameters.MType,
			Value: &value.Value,
		}, nil
	case "counter":
		value, err := m.storage.GetCounterMetric(ctx, parameters.ID)

		if err != nil {
			if errors.Is(err, storage.ErrDataNotFound) {
				return models.Metrics{}, services.ErrMetricNotFound
			}

			return models.Metrics{}, err
		}

		return models.Metrics{
			ID:    parameters.ID,
			MType: parameters.MType,
			Delta: &value.Value,
		}, nil
	default:
		return models.Metrics{}, services.ErrMetricNotFound
	}
}

// GetAll извлекает все метрики из хранилища и формирует список для отображения.
func (m *Metrics) GetAll(ctx context.Context) ([]services.MetricEntry, error) {
	var result []services.MetricEntry

	gaugeMetrics, err := m.storage.GetGaugeMetrics(ctx)

	if err != nil {
		return nil, err
	}

	counterMetrics, err := m.storage.GetCounterMetrics(ctx)

	if err != nil {
		return nil, err
	}

	for _, item := range gaugeMetrics {
		result = append(result, services.MetricEntry{Name: item.Name, Value: fmt.Sprintf("%g", item.Value)})
	}

	for _, item := range counterMetrics {
		result = append(result, services.MetricEntry{Name: item.Name, Value: fmt.Sprintf("%v", item.Value)})
	}

	return result, nil
}

// IsCounterMetricType определяет, является ли метрика счетчиком.
func IsCounterMetricType(metricName string) bool {
	return metricName == "PollCount"
}

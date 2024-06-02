package filestorage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/daremove/go-metrics-service/internal/storage"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	gaugeMetrics   map[string]storage.GaugeMetric
	counterMetrics map[string]storage.CounterMetric
	returnError    bool
}

func (m *MockStorage) GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error) {
	if m.returnError {
		return storage.GaugeMetric{}, errors.New("error")
	}
	metric, ok := m.gaugeMetrics[key]
	if !ok {
		return storage.GaugeMetric{}, errors.New("not found")
	}
	return metric, nil
}

func (m *MockStorage) GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error) {
	if m.returnError {
		return nil, errors.New("error")
	}
	metrics := make([]storage.GaugeMetric, 0, len(m.gaugeMetrics))
	for _, metric := range m.gaugeMetrics {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m *MockStorage) GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error) {
	if m.returnError {
		return storage.CounterMetric{}, errors.New("error")
	}
	metric, ok := m.counterMetrics[key]
	if !ok {
		return storage.CounterMetric{}, errors.New("not found")
	}
	return metric, nil
}

func (m *MockStorage) GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error) {
	if m.returnError {
		return nil, errors.New("error")
	}
	metrics := make([]storage.CounterMetric, 0, len(m.counterMetrics))
	for _, metric := range m.counterMetrics {
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (m *MockStorage) AddGaugeMetric(ctx context.Context, key string, value float64) error {
	if m.returnError {
		return errors.New("error")
	}
	m.gaugeMetrics[key] = storage.GaugeMetric{Name: key, Value: value}
	return nil
}

func (m *MockStorage) AddCounterMetric(ctx context.Context, key string, value int64) error {
	if m.returnError {
		return errors.New("error")
	}
	m.counterMetrics[key] = storage.CounterMetric{Name: key, Value: value}
	return nil
}

func (m *MockStorage) AddMetrics(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) error {
	if m.returnError {
		return errors.New("error")
	}
	for _, metric := range gaugeMetrics {
		m.gaugeMetrics[metric.Name] = metric
	}
	for _, metric := range counterMetrics {
		m.counterMetrics[metric.Name] = metric
	}
	return nil
}

func TestFileStorage(t *testing.T) {
	t.Run("Should successfully backup data", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics: map[string]storage.GaugeMetric{
				"test_gauge": {
					Name:  "test_gauge",
					Value: 123.45,
				},
			},
			counterMetrics: map[string]storage.CounterMetric{
				"test_counter": {
					Name:  "test_counter",
					Value: 678,
				},
			},
		}
		config := Config{FileStoragePath: "test_backup.json", Restore: false}
		defer os.Remove(config.FileStoragePath)
		fs, err := New(context.Background(), mockStorage, config)
		assert.NoError(t, err)
		assert.NotNil(t, fs)

		err = fs.BackupData(context.Background())
		assert.NoError(t, err)

		data, err := os.ReadFile(config.FileStoragePath)
		assert.NoError(t, err)
		var backupData backupFile
		err = json.Unmarshal(data, &backupData)
		assert.NoError(t, err)

		assert.Len(t, backupData.Gauges, 1)
		assert.Equal(t, "test_gauge", backupData.Gauges[0].Name)
		assert.Equal(t, 123.45, backupData.Gauges[0].Value)

		assert.Len(t, backupData.Counters, 1)
		assert.Equal(t, "test_counter", backupData.Counters[0].Name)
		assert.Equal(t, int64(678), backupData.Counters[0].Value)
	})

	t.Run("Should return all counter metrics", func(t *testing.T) {
		mockStorage := &MockStorage{
			counterMetrics: map[string]storage.CounterMetric{
				"test_counter1": {
					Name:  "test_counter1",
					Value: 1,
				},
				"test_counter2": {
					Name:  "test_counter2",
					Value: 2,
				},
			},
		}
		config := Config{FileStoragePath: "test_counter.json", Restore: false}
		defer os.Remove(config.FileStoragePath)
		fs, err := New(context.Background(), mockStorage, config)

		assert.NoError(t, err)
		assert.NotNil(t, fs)

		data, err := fs.GetCounterMetrics(context.Background())

		assert.NoError(t, err)
		assert.ElementsMatch(t, []storage.CounterMetric{{Name: "test_counter1", Value: 1}, {Name: "test_counter2", Value: 2}}, data)
	})

	t.Run("Should return all gauge metrics", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics: map[string]storage.GaugeMetric{
				"test_gauge1": {
					Name:  "test_gauge1",
					Value: 1.2,
				},
				"test_gauge2": {
					Name:  "test_gauge2",
					Value: 3.4,
				},
			},
		}
		config := Config{FileStoragePath: "test_gauge.json", Restore: false}
		defer os.Remove(config.FileStoragePath)
		fs, err := New(context.Background(), mockStorage, config)

		assert.NoError(t, err)
		assert.NotNil(t, fs)

		data, err := fs.GetGaugeMetrics(context.Background())

		assert.NoError(t, err)
		assert.ElementsMatch(t, []storage.GaugeMetric{{Name: "test_gauge1", Value: 1.2}, {Name: "test_gauge2", Value: 3.4}}, data)
	})

	t.Run("Should add a counter metric and backup data", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
		}
		config := Config{FileStoragePath: "add_counter.json", Restore: false}
		defer os.Remove(config.FileStoragePath)
		fs, err := New(context.Background(), mockStorage, config)

		assert.NoError(t, err)
		assert.NotNil(t, fs)

		err = fs.AddCounterMetric(context.Background(), "counter", 1)
		assert.NoError(t, err)

		data, err := os.ReadFile(config.FileStoragePath)
		assert.NoError(t, err)
		var backupData backupFile
		err = json.Unmarshal(data, &backupData)
		assert.NoError(t, err)

		assert.Len(t, backupData.Counters, 1)
		assert.Equal(t, "counter", backupData.Counters[0].Name)
		assert.Equal(t, int64(1), backupData.Counters[0].Value)
	})

	t.Run("Should add metrics and backup data", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
		}
		config := Config{FileStoragePath: "add_metics.json", Restore: false}
		defer os.Remove(config.FileStoragePath)
		fs, err := New(context.Background(), mockStorage, config)

		assert.NoError(t, err)
		assert.NotNil(t, fs)

		err = fs.AddMetrics(context.Background(), []storage.GaugeMetric{{Name: "gauge", Value: 1.23}}, []storage.CounterMetric{{Name: "counter", Value: 456}})
		assert.NoError(t, err)

		data, err := os.ReadFile(config.FileStoragePath)
		assert.NoError(t, err)
		var backupData backupFile
		err = json.Unmarshal(data, &backupData)
		assert.NoError(t, err)

		assert.Len(t, backupData.Counters, 1)
		assert.Equal(t, "counter", backupData.Counters[0].Name)
		assert.Equal(t, int64(456), backupData.Counters[0].Value)

		assert.Len(t, backupData.Gauges, 1)
		assert.Equal(t, "gauge", backupData.Gauges[0].Name)
		assert.Equal(t, 1.23, backupData.Gauges[0].Value)
	})

	t.Run("Should create a new FileStorage instance", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
		}
		config := Config{FileStoragePath: "", Restore: false}
		fs, err := New(context.Background(), mockStorage, config)
		assert.NoError(t, err)
		assert.NotNil(t, fs)
	})

	t.Run("Should add a gauge metric and backup data", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
		}
		config := Config{FileStoragePath: "test_backup.json", Restore: false}
		defer os.Remove(config.FileStoragePath)
		fs, err := New(context.Background(), mockStorage, config)
		assert.NoError(t, err)

		err = fs.AddGaugeMetric(context.Background(), "test_gauge", 123.45)
		assert.NoError(t, err)

		data, err := os.ReadFile(config.FileStoragePath)
		assert.NoError(t, err)
		var backupData backupFile
		err = json.Unmarshal(data, &backupData)
		assert.NoError(t, err)
		assert.Len(t, backupData.Gauges, 1)
		assert.Equal(t, "test_gauge", backupData.Gauges[0].Name)
		assert.Equal(t, 123.45, backupData.Gauges[0].Value)
	})

	t.Run("Should return error if backup fails", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
			returnError:    true,
		}
		config := Config{FileStoragePath: "test_backup.json", Restore: false}
		defer os.Remove(config.FileStoragePath)
		fs, err := New(context.Background(), mockStorage, config)
		assert.NoError(t, err)

		err = fs.AddGaugeMetric(context.Background(), "test_gauge", 123.45)
		assert.Error(t, err)
	})

	t.Run("Should restore data from file", func(t *testing.T) {
		backupData := backupFile{
			Gauges: []storage.GaugeMetric{
				{Name: "test_gauge", Value: 123.45},
			},
			Counters: []storage.CounterMetric{
				{Name: "test_counter", Value: 678},
			},
		}
		data, _ := json.Marshal(backupData)
		filePath := "test_restore.json"
		_ = os.WriteFile(filePath, data, 0666)
		defer os.Remove(filePath)

		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
		}
		config := Config{FileStoragePath: filePath, Restore: true}
		fs, err := New(context.Background(), mockStorage, config)
		assert.NoError(t, err)
		assert.NotNil(t, fs)

		gauge, err := fs.GetGaugeMetric(context.Background(), "test_gauge")
		assert.NoError(t, err)
		assert.Equal(t, "test_gauge", gauge.Name)
		assert.Equal(t, 123.45, gauge.Value)

		counter, err := fs.GetCounterMetric(context.Background(), "test_counter")
		assert.NoError(t, err)
		assert.Equal(t, "test_counter", counter.Name)
		assert.Equal(t, int64(678), counter.Value)
	})

	t.Run("Should return error if file does not exist during restore", func(t *testing.T) {
		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
		}
		config := Config{FileStoragePath: "non_existent.json", Restore: true}
		fs, err := New(context.Background(), mockStorage, config)
		assert.NoError(t, err)
		assert.NotNil(t, fs)
	})

	t.Run("Should return error if deserialization fails during restore", func(t *testing.T) {
		invalidData := []byte(`invalid json`)
		filePath := "test_invalid.json"
		_ = os.WriteFile(filePath, invalidData, 0666)
		defer os.Remove(filePath)

		mockStorage := &MockStorage{
			gaugeMetrics:   make(map[string]storage.GaugeMetric),
			counterMetrics: make(map[string]storage.CounterMetric),
		}
		config := Config{FileStoragePath: filePath, Restore: true}
		fs, err := New(context.Background(), mockStorage, config)
		assert.Error(t, err)
		assert.Nil(t, fs)
	})
}

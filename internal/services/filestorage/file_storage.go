// Пакет filestorage предоставляет реализацию хранилища данных, использующего файловую систему для сохранения и восстановления метрик.
package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/storage"
	"go.uber.org/zap"
)

// backupFile структура для сериализации и десериализации данных метрик в файл.
type backupFile struct {
	Counters []storage.CounterMetric `json:"counters"`
	Gauges   []storage.GaugeMetric   `json:"gauges"`
}

// FileStorage реализует интерфейс Storage, предоставляя методы для работы с метриками, хранящимися в файле.
type FileStorage struct {
	Storage
	storage Storage // Внутреннее хранилище для делегирования операций
	config  Config  // Конфигурация хранилища
}

// Storage интерфейс определяет методы, которые должны быть реализованы для работы с метриками.
type Storage interface {
	GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error)
	GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error)

	GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error)
	GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error)

	AddGaugeMetric(ctx context.Context, key string, value float64) error
	AddCounterMetric(ctx context.Context, key string, value int64) error

	AddMetrics(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) error
}

// Config структура конфигурации FileStorage.
type Config struct {
	StoreInterval   int    // Интервал автоматического сохранения данных в файл
	FileStoragePath string // Путь к файлу хранения данных
	Restore         bool   // Флаг, указывающий на необходимость восстановления данных из файла при инициализации
}

// AddGaugeMetric добавляет значение типа gauge в хранилище и выполняет бэкап данных при необходимости.
func (fs FileStorage) AddGaugeMetric(ctx context.Context, key string, value float64) error {
	if err := fs.storage.AddGaugeMetric(ctx, key, value); err != nil {
		return err
	}

	if fs.config.StoreInterval > 0 {
		return nil
	}

	if err := backupData(ctx, fs.storage, fs.config.FileStoragePath); err != nil {
		return fmt.Errorf("error has occurred during backup data: %s", err)
	}

	return nil
}

// AddCounterMetric добавляет значение типа counter в хранилище и выполняет бэкап данных при необходимости.
func (fs FileStorage) AddCounterMetric(ctx context.Context, key string, value int64) error {
	if err := fs.storage.AddCounterMetric(ctx, key, value); err != nil {
		return err
	}

	if fs.config.StoreInterval > 0 {
		return nil
	}

	if err := backupData(ctx, fs.storage, fs.config.FileStoragePath); err != nil {
		return fmt.Errorf("error has occurred during backup data: %s", err)
	}

	return nil
}

// AddMetrics добавляет несколько метрик в хранилище и выполняет бэкап данных при необходимости.
func (fs FileStorage) AddMetrics(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) error {
	if err := fs.storage.AddMetrics(ctx, gaugeMetrics, counterMetrics); err != nil {
		return err
	}

	if fs.config.StoreInterval > 0 {
		return nil
	}

	if err := backupData(ctx, fs.storage, fs.config.FileStoragePath); err != nil {
		return fmt.Errorf("error has occurred during backup data: %s", err)
	}

	return nil
}

// GetGaugeMetric извлекает метрику типа gauge из хранилища.
func (fs FileStorage) GetGaugeMetric(ctx context.Context, key string) (storage.GaugeMetric, error) {
	return fs.storage.GetGaugeMetric(ctx, key)
}

// GetGaugeMetrics извлекает все метрики типа gauge из хранилища.
func (fs FileStorage) GetGaugeMetrics(ctx context.Context) ([]storage.GaugeMetric, error) {
	return fs.storage.GetGaugeMetrics(ctx)
}

// GetCounterMetric извлекает метрику типа counter из хранилища.
func (fs FileStorage) GetCounterMetric(ctx context.Context, key string) (storage.CounterMetric, error) {
	return fs.storage.GetCounterMetric(ctx, key)
}

// GetCounterMetrics извлекает все метрики типа counter из хранилища.
func (fs FileStorage) GetCounterMetrics(ctx context.Context) ([]storage.CounterMetric, error) {
	return fs.storage.GetCounterMetrics(ctx)
}

func backupData(ctx context.Context, fs Storage, filePath string) error {
	counterMetrics, err := fs.GetCounterMetrics(ctx)

	if err != nil {
		return fmt.Errorf("cannot serialize data: %s", err)
	}

	gaugeMetrics, err := fs.GetGaugeMetrics(ctx)

	if err != nil {
		return fmt.Errorf("cannot serialize data: %s", err)
	}

	serialisedData, err := json.Marshal(backupFile{
		Counters: counterMetrics,
		Gauges:   gaugeMetrics,
	})

	if err != nil {
		return fmt.Errorf("cannot serialize data: %s", err)
	}

	if err := os.WriteFile(filePath, serialisedData, 0666); err != nil {
		return fmt.Errorf("cannot write data to file: %s", err)
	}

	return nil
}

// New создает новый экземпляр FileStorage.
func New(ctx context.Context, storage Storage, config Config) (*FileStorage, error) {
	fileStorage := &FileStorage{
		storage: storage,
		config:  config,
	}

	if config.FileStoragePath == "" {
		return fileStorage, nil
	}

	if config.Restore {
		data, err := os.ReadFile(config.FileStoragePath)

		if err == nil {
			backupData := backupFile{}
			err = json.Unmarshal(data, &backupData)

			if err != nil {
				return nil, fmt.Errorf("cannot deserialize data from file: %s", err)
			}

			for _, counter := range backupData.Counters {
				if err = storage.AddCounterMetric(ctx, counter.Name, counter.Value); err != nil {
					return nil, fmt.Errorf("cannot initialize counter data from file: %s", err)
				}
			}

			for _, gauge := range backupData.Gauges {
				if err = storage.AddGaugeMetric(ctx, gauge.Name, gauge.Value); err != nil {
					return nil, fmt.Errorf("cannot initialize gauge data from file: %s", err)
				}
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot read file: %s", err)
		}
	}

	if config.StoreInterval > 0 {
		go func() {
			for {
				time.Sleep(time.Duration(config.StoreInterval) * time.Second)

				if err := fileStorage.BackupData(ctx); err != nil {
					logger.Log.Error("error has occurred during backup data", zap.Error(err))
					continue
				}
			}
		}()
	}

	return fileStorage, nil
}

// BackupData выполняет бэкап всех данных в файл.
func (fs FileStorage) BackupData(ctx context.Context) error {
	return backupData(ctx, fs.storage, fs.config.FileStoragePath)
}

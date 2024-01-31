package filestorage

import (
	"encoding/json"
	"fmt"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/storage"
	"go.uber.org/zap"
	"os"
	"time"
)

type backupFile struct {
	Counters []storage.CounterMetric `json:"counters"`
	Gauges   []storage.GaugeMetric   `json:"gauges"`
}

type FileStorage struct {
	Storage
	storage Storage
	config  Config
}

type Storage interface {
	GetGaugeMetric(key string) (float64, bool)
	GetGaugeMetrics() []storage.GaugeMetric

	GetCounterMetric(key string) (int64, bool)
	GetCounterMetrics() []storage.CounterMetric

	AddGauge(key string, value float64) error
	AddCounter(key string, value int64) error
}

type Config struct {
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

func (fs FileStorage) AddGauge(key string, value float64) error {
	if err := fs.storage.AddGauge(key, value); err != nil {
		return err
	}

	if fs.config.StoreInterval > 0 {
		return nil
	}

	if err := backupData(fs.storage, fs.config.FileStoragePath); err != nil {
		return fmt.Errorf("error has occurred during backup data: %s", err)
	}

	return nil
}

func (fs FileStorage) AddCounter(key string, value int64) error {
	if err := fs.storage.AddCounter(key, value); err != nil {
		return err
	}

	if fs.config.StoreInterval > 0 {
		return nil
	}

	if err := backupData(fs.storage, fs.config.FileStoragePath); err != nil {
		return fmt.Errorf("error has occurred during backup data: %s", err)
	}

	return nil
}

func (fs FileStorage) GetGaugeMetric(key string) (float64, bool) {
	return fs.storage.GetGaugeMetric(key)
}
func (fs FileStorage) GetGaugeMetrics() []storage.GaugeMetric {
	return fs.storage.GetGaugeMetrics()
}
func (fs FileStorage) GetCounterMetric(key string) (int64, bool) {
	return fs.storage.GetCounterMetric(key)
}
func (fs FileStorage) GetCounterMetrics() []storage.CounterMetric {
	return fs.storage.GetCounterMetrics()
}

func backupData(fs Storage, filePath string) error {
	serialisedData, err := json.Marshal(backupFile{
		Counters: fs.GetCounterMetrics(),
		Gauges:   fs.GetGaugeMetrics(),
	})

	if err != nil {
		return fmt.Errorf("cannot serialize data: %s", err)
	}

	if err := os.WriteFile(filePath, serialisedData, 0666); err != nil {
		return fmt.Errorf("cannot write data to file: %s", err)
	}

	return nil
}

func New(storage Storage, config Config) (*FileStorage, error) {
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
				if err = storage.AddCounter(counter.Name, counter.Value); err != nil {
					return nil, fmt.Errorf("cannot initialize counter data from file: %s", err)
				}
			}

			for _, gauge := range backupData.Gauges {
				if err = storage.AddGauge(gauge.Name, gauge.Value); err != nil {
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

				if err := fileStorage.BackupData(); err != nil {
					logger.Log.Error("error has occurred during backup data", zap.Error(err))
					continue
				}
			}
		}()
	}

	return fileStorage, nil
}

func (fs FileStorage) BackupData() error {
	return backupData(fs.storage, fs.config.FileStoragePath)
}

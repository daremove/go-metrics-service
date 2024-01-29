package backup

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

type Backup struct {
	storage      Storage
	config       Config
	proxyStorage proxyStorage
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

type proxyStorage struct {
	originalStorage Storage
	config          Config
}

func (storage proxyStorage) AddGauge(key string, value float64) error {
	if err := storage.originalStorage.AddGauge(key, value); err != nil {
		return err
	}

	if err := backupData(storage.originalStorage, storage.config.FileStoragePath); err != nil {
		return fmt.Errorf("error has occurred during backup data: %s", err)
	}

	return nil
}

func (storage proxyStorage) AddCounter(key string, value int64) error {
	if err := storage.originalStorage.AddCounter(key, value); err != nil {
		return err
	}

	if err := backupData(storage.originalStorage, storage.config.FileStoragePath); err != nil {
		return fmt.Errorf("error has occurred during backup data: %s", err)
	}

	return nil
}

func (storage proxyStorage) GetGaugeMetric(key string) (float64, bool) {
	return storage.originalStorage.GetGaugeMetric(key)
}
func (storage proxyStorage) GetGaugeMetrics() []storage.GaugeMetric {
	return storage.originalStorage.GetGaugeMetrics()
}
func (storage proxyStorage) GetCounterMetric(key string) (int64, bool) {
	return storage.originalStorage.GetCounterMetric(key)
}
func (storage proxyStorage) GetCounterMetrics() []storage.CounterMetric {
	return storage.originalStorage.GetCounterMetrics()
}

func backupData(storage Storage, filePath string) error {
	serialisedData, err := json.Marshal(backupFile{
		Counters: storage.GetCounterMetrics(),
		Gauges:   storage.GetGaugeMetrics(),
	})

	if err != nil {
		return fmt.Errorf("cannot serialize data: %s", err)
	}

	if err := os.WriteFile(filePath, serialisedData, 0666); err != nil {
		return fmt.Errorf("cannot write data to file: %s", err)
	}

	return nil
}

func New(storage Storage, config Config) (*Backup, error) {
	backupService := &Backup{
		storage:      storage,
		config:       config,
		proxyStorage: proxyStorage{originalStorage: storage, config: config},
	}

	if config.FileStoragePath == "" {
		return backupService, nil
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

				if err := backupData(storage, config.FileStoragePath); err != nil {
					logger.Log.Error("error has occurred during backup data", zap.Error(err))
					continue
				}
			}
		}()
	}

	return backupService, nil
}

func (b *Backup) GetProxyStorage() *Storage {
	if b.config.StoreInterval != 0 {
		return &b.storage
	}

	var st Storage = b.proxyStorage

	return &st
}

func (b *Backup) BackupData() error {
	return backupData(b.storage, b.config.FileStoragePath)
}

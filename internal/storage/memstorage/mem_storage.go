package memstorage

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/storage"
)

type memStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (storage *memStorage) AddGauge(key string, value float64) error {
	storage.gauge[key] = value

	return nil
}

func (storage *memStorage) AddCounter(key string, value int64) error {
	storage.counter[key] += value
	fmt.Println(storage.counter[key])

	return nil
}

func New() storage.Storage {
	return &memStorage{gauge: map[string]float64{}, counter: map[string]int64{}}
}

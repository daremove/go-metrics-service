package storage

type Storage interface {
	AddGauge(key string, value float64) error
	AddCounter(key string, value int64) error
}

package memstorage

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (storage *MemStorage) AddGauge(key string, value float64) {
	storage.gauge[key] = value
}

func (storage *MemStorage) AddCounter(key string, value int64) {
	storage.counter[key] += value
}

func New() *MemStorage {
	return &MemStorage{gauge: map[string]float64{}, counter: map[string]int64{}}
}

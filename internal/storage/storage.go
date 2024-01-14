package storage

type GaugeMetric struct {
	Name  string
	Value float64
}

type CounterMetric struct {
	Name  string
	Value int64
}

type Storage interface {
	GetGaugeMetric(key string) (float64, bool)
	GetGaugeMetrics() []GaugeMetric

	GetCounterMetric(key string) (int64, bool)
	GetCounterMetrics() []CounterMetric

	AddGauge(key string, value float64) error
	AddCounter(key string, value int64) error
}

// Package models предназначен для структур данных, используемых во всем приложении.
package models

const (
	GaugeMetricType   string = "gauge"
	CounterMetricType string = "counter"
)

// Metrics описывает структуру данных метрики, которая может быть типа "gauge" или "counter".
type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Тип метрики, принимает значения "gauge" или "counter"
	Delta *int64   `json:"delta,omitempty"` // Изменение значения для метрик типа "counter"
	Value *float64 `json:"value,omitempty"` // Текущее значение для метрик типа "gauge"
}

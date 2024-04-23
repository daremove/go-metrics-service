// Пакет storage предоставляет определения структур данных и ошибок для системы метрик.
package storage

import "errors"

// ErrDataNotFound ошибка, возникающая когда запрашиваемые данные не найдены.
var ErrDataNotFound = errors.New("data is not found")

// GaugeMetric определяет структуру для метрик типа "gauge", которые представляют собой мгновенное значение.
type GaugeMetric struct {
	Name  string  `json:"name"`  // Имя метрики
	Value float64 `json:"value"` // Значение метрики
}

// CounterMetric определяет структуру для метрик типа "counter", которые накапливают значение со временем.
type CounterMetric struct {
	Name  string `json:"name"`  // Имя метрики
	Value int64  `json:"value"` // Накопленное значение метрики
}

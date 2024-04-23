// Пакет services предоставляет определения типов и переменных ошибок,
// используемых для управления метрическими данными в приложении.
package services

import "errors"

// ErrMetricNotFound ошибка, возвращаемая когда данные метрики не найдены.
var ErrMetricNotFound = errors.New("metric data is not found")

// MetricEntry представляет базовую запись метрики с именем и значением.
type MetricEntry struct {
	Name  string // Имя метрики
	Value string // Значение метрики
}

// MetricSaveParameters содержит параметры для сохранения метрики.
type MetricSaveParameters struct {
	MetricType  string // Тип метрики, например "gauge" или "counter"
	MetricName  string // Имя метрики
	MetricValue string // Значение метрики
}

// MetricGetParameters содержит параметры для получения метрики.
type MetricGetParameters struct {
	MetricType string // Тип метрики
	MetricName string // Имя метрики
}

// Package healthcheck предоставляет инструменты для проверки состояния и доступности хранилища данных.
package healthcheck

import (
	"context"
	"fmt"
)

// HealthCheck структура для проведения проверок состояния хранилища данных.
type HealthCheck struct {
	storage Storage // Интерфейс хранилища, поддерживающий операцию Ping.
}

// Storage интерфейс определяет методы, необходимые для проверки состояния хранилища.
type Storage interface {
	Ping(ctx context.Context) error // Ping проверяет доступность хранилища.
}

// New создает новый экземпляр HealthCheck.
func New(storage Storage) *HealthCheck {
	return &HealthCheck{
		storage,
	}
}

// CheckStorageConnection выполняет проверку соединения с хранилищем данных.
func (hc *HealthCheck) CheckStorageConnection(ctx context.Context) error {
	if hc.storage == nil {
		return fmt.Errorf("storage wasn't initialized")
	}

	return hc.storage.Ping(ctx)
}

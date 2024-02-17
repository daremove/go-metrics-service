package healthcheck

import (
	"context"
	"fmt"
)

type HealthCheck struct {
	storage Storage
}

type Storage interface {
	Ping(ctx context.Context) error
}

func New(storage Storage) *HealthCheck {
	return &HealthCheck{
		storage,
	}
}

func (hc *HealthCheck) CheckStorageConnection(ctx context.Context) error {
	if hc.storage == nil {
		return fmt.Errorf("storage wasn't initialized")
	}

	return hc.storage.Ping(ctx)
}

package healthcheck

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	shouldError bool
}

func (m *MockStorage) Ping(ctx context.Context) error {
	if m.shouldError {
		return errors.New("storage is unavailable")
	}
	return nil
}

func TestHealthCheck(t *testing.T) {
	t.Run("Should create a new HealthCheck instance", func(t *testing.T) {
		storage := &MockStorage{}
		hc := New(storage)
		assert.NotNil(t, hc)
		assert.NotNil(t, hc.storage)
	})

	t.Run("Should return an error if storage is not initialized", func(t *testing.T) {
		hc := &HealthCheck{}
		err := hc.CheckStorageConnection(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "storage wasn't initialized", err.Error())
	})

	t.Run("Should successfully check storage connection", func(t *testing.T) {
		storage := &MockStorage{}
		hc := New(storage)
		err := hc.CheckStorageConnection(context.Background())
		assert.NoError(t, err)
	})

	t.Run("Should return an error if storage ping fails", func(t *testing.T) {
		storage := &MockStorage{shouldError: true}
		hc := New(storage)
		err := hc.CheckStorageConnection(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "storage is unavailable", err.Error())
	})
}

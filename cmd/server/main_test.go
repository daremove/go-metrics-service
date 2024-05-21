package main

import (
	"context"
	"testing"

	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services/healthcheck"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitializeLogger(t *testing.T) {
	t.Run("Should initialize logger with valid log level", func(t *testing.T) {
		err := initializeLogger("info")
		assert.NoError(t, err)
		assert.NotNil(t, logger.Log)
	})

	t.Run("Should return error for invalid log level", func(t *testing.T) {
		err := initializeLogger("invalid")
		assert.Error(t, err)
	})
}

func TestInitializeStorage(t *testing.T) {
	t.Run("Should initialize file storage", func(t *testing.T) {
		ctx := context.Background()
		config := Config{
			StoreInterval:   300,
			FileStoragePath: "/tmp/metrics-db.json",
			Restore:         true,
		}

		storage, healthCheckService, err := initializeStorage(ctx, config)
		require.NoError(t, err)
		assert.NotNil(t, storage)
		assert.NotNil(t, healthCheckService)
	})
}

func TestRunServer(t *testing.T) {
	t.Run("Should run server", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		config := Config{
			Endpoint:   "localhost:8080",
			SigningKey: "test-signing-key",
		}

		storage := memstorage.New()
		healthCheckService := healthcheck.New(nil)

		go func() {
			runServer(ctx, config, storage, healthCheckService, nil)
		}()

		cancel()
	})
}

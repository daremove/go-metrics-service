package main

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/daremove/go-metrics-service/internal/services/metrics"

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
	FileStoragePath := "/tmp/test-init-storage.json"

	if _, err := os.Stat(FileStoragePath); err == nil {
		err := os.Remove(FileStoragePath)

		if err != nil {
			log.Fatalf("Failed to delete file %s: %v", FileStoragePath, err)
		}
	}

	t.Run("Should initialize file storage", func(t *testing.T) {
		ctx := context.Background()
		config := Config{
			StoreInterval:   300,
			FileStoragePath: FileStoragePath,
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
		metricsService := metrics.New(storage)
		healthCheckService := healthcheck.New(nil)

		go func() {
			runServer(ctx, config, metricsService, healthCheckService, nil)
		}()

		cancel()
	})
}

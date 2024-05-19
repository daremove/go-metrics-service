package memstorage

import (
	"context"
	"testing"

	"github.com/daremove/go-metrics-service/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStorage(t *testing.T) {
	ctx := context.Background()
	memStore := New()

	t.Run("Should add and retrieve a gauge metric correctly", func(t *testing.T) {
		err := memStore.AddGaugeMetric(ctx, "temperature", 25.5)
		require.NoError(t, err)

		metric, err := memStore.GetGaugeMetric(ctx, "temperature")
		require.NoError(t, err)
		assert.Equal(t, storage.GaugeMetric{Name: "temperature", Value: 25.5}, metric)
	})

	t.Run("Should return an error when retrieving a non-existent gauge metric", func(t *testing.T) {
		_, err := memStore.GetGaugeMetric(ctx, "humidity")
		assert.Equal(t, storage.ErrDataNotFound, err)
	})

	t.Run("Should add and increment a counter metric correctly", func(t *testing.T) {
		err := memStore.AddCounterMetric(ctx, "requests", 5)
		require.NoError(t, err)

		err = memStore.AddCounterMetric(ctx, "requests", 3)
		require.NoError(t, err)

		metric, err := memStore.GetCounterMetric(ctx, "requests")
		require.NoError(t, err)
		assert.Equal(t, storage.CounterMetric{Name: "requests", Value: 8}, metric)
	})

	t.Run("Should return an error when retrieving a non-existent counter metric", func(t *testing.T) {
		_, err := memStore.GetCounterMetric(ctx, "errors")
		assert.Equal(t, storage.ErrDataNotFound, err)
	})

	t.Run("Should retrieve all stored gauge and counter metrics correctly", func(t *testing.T) {
		gauges, err := memStore.GetGaugeMetrics(ctx)
		require.NoError(t, err)
		assert.Len(t, gauges, 1)

		counters, err := memStore.GetCounterMetrics(ctx)
		require.NoError(t, err)
		assert.Len(t, counters, 1)
	})

	t.Run("Should add multiple gauge and counter metrics successfully", func(t *testing.T) {
		gaugesToAdd := []storage.GaugeMetric{{Name: "speed", Value: 88.0}}
		countersToAdd := []storage.CounterMetric{{Name: "errors", Value: 1}}

		err := memStore.AddMetrics(ctx, gaugesToAdd, countersToAdd)
		require.NoError(t, err)

		gauges, err := memStore.GetGaugeMetrics(ctx)
		require.NoError(t, err)
		assert.Len(t, gauges, 2)

		counters, err := memStore.GetCounterMetrics(ctx)
		require.NoError(t, err)
		assert.Len(t, counters, 2)
	})
}

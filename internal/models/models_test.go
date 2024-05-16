package models

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricsSerialization(t *testing.T) {
	t.Run("Should correctly serialize and deserialize Gauge metric", func(t *testing.T) {
		value := 123.45
		metric := Metrics{
			ID:    "TestGauge",
			MType: GaugeMetricType,
			Value: &value,
		}

		data, err := json.Marshal(metric)
		assert.NoError(t, err)

		var deserializedMetric Metrics
		err = json.Unmarshal(data, &deserializedMetric)
		assert.NoError(t, err)

		assert.Equal(t, metric, deserializedMetric)
	})

	t.Run("Should correctly serialize and deserialize Counter metric", func(t *testing.T) {
		delta := int64(123456)
		metric := Metrics{
			ID:    "TestCounter",
			MType: CounterMetricType,
			Delta: &delta,
		}

		data, err := json.Marshal(metric)
		assert.NoError(t, err)

		var deserializedMetric Metrics
		err = json.Unmarshal(data, &deserializedMetric)
		assert.NoError(t, err)

		assert.Equal(t, metric, deserializedMetric)
	})

	t.Run("Should handle missing optional fields", func(t *testing.T) {
		metric := Metrics{
			ID:    "TestMissingFields",
			MType: GaugeMetricType,
		}

		data, err := json.Marshal(metric)
		assert.NoError(t, err)

		var deserializedMetric Metrics
		err = json.Unmarshal(data, &deserializedMetric)
		assert.NoError(t, err)

		assert.Equal(t, metric.ID, deserializedMetric.ID)
		assert.Equal(t, metric.MType, deserializedMetric.MType)
		assert.Nil(t, deserializedMetric.Delta)
		assert.Nil(t, deserializedMetric.Value)
	})
}

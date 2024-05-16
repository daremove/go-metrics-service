package storage

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGaugeMetricSerialization(t *testing.T) {
	t.Run("Should correctly serialize and deserialize GaugeMetric", func(t *testing.T) {
		metric := GaugeMetric{
			Name:  "TestGauge",
			Value: 123.45,
		}

		data, err := json.Marshal(metric)
		assert.NoError(t, err)

		var deserializedMetric GaugeMetric
		err = json.Unmarshal(data, &deserializedMetric)
		assert.NoError(t, err)

		assert.Equal(t, metric, deserializedMetric)
	})
}

func TestCounterMetricSerialization(t *testing.T) {
	t.Run("Should correctly serialize and deserialize CounterMetric", func(t *testing.T) {
		metric := CounterMetric{
			Name:  "TestCounter",
			Value: 123456,
		}

		data, err := json.Marshal(metric)
		assert.NoError(t, err)

		var deserializedMetric CounterMetric
		err = json.Unmarshal(data, &deserializedMetric)
		assert.NoError(t, err)

		assert.Equal(t, metric, deserializedMetric)
	})
}

func TestErrDataNotFound(t *testing.T) {
	t.Run("Should return correct error message", func(t *testing.T) {
		err := ErrDataNotFound
		assert.EqualError(t, err, "data is not found")
	})
}

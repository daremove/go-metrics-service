package metrics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type storageMock struct {
	gauge   map[string]float64
	counter map[string]int64
}

func (s storageMock) AddGauge(key string, value float64) error {
	s.gauge[key] = value
	return nil
}

func (s storageMock) AddCounter(key string, value int64) error {
	s.counter[key] = value
	return nil
}

func TestMetrics_Save(t *testing.T) {
	storeMock := storageMock{gauge: make(map[string]float64), counter: make(map[string]int64)}
	metricsService := New(storeMock)

	testCases := []struct {
		testName       string
		saveParameters SaveParameters
		testCase       func(t *testing.T, err error)
	}{
		{
			testName: "Should return error if metricType isn't defined",
			saveParameters: SaveParameters{
				MetricType: "incorrect",
			},
			testCase: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			testName: "Should return error if value of gauge metric type wasn't parsed correctly",
			saveParameters: SaveParameters{
				MetricType:  "gauge",
				MetricName:  "metricName",
				MetricValue: "string",
			},
			testCase: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			testName: "Should return error if value of counter metric type wasn't parsed correctly",
			saveParameters: SaveParameters{
				MetricType:  "counter",
				MetricName:  "metricName",
				MetricValue: "1.1",
			},
			testCase: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			testName: "Should save in storage gauge metric type",
			saveParameters: SaveParameters{
				MetricType:  "gauge",
				MetricName:  "metricName",
				MetricValue: "1.1",
			},
			testCase: func(t *testing.T, err error) {
				require.NoError(t, err)
				assert.Equal(t, 1.1, storeMock.gauge["metricName"])
			},
		},
		{
			testName: "Should save in storage counter metric type",
			saveParameters: SaveParameters{
				MetricType:  "counter",
				MetricName:  "metricName",
				MetricValue: "100",
			},
			testCase: func(t *testing.T, err error) {
				require.NoError(t, err)
				assert.Equal(t, int64(100), storeMock.counter["metricName"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			err := metricsService.Save(tc.saveParameters)

			tc.testCase(t, err)
		})
	}
}

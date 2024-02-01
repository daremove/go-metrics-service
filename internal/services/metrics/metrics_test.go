package metrics

import (
	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMetrics_Save(t *testing.T) {
	storeMock := memstorage.NewWithPrefilledData(map[string]float64{}, map[string]int64{})
	metricsService := New(storeMock)

	testCases := []struct {
		testName       string
		saveParameters services.MetricSaveParameters
		testCase       func(t *testing.T, err error)
	}{
		{
			testName: "Should return error if metricType isn't defined",
			saveParameters: services.MetricSaveParameters{
				MetricType: "incorrect",
			},
			testCase: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			testName: "Should return error if value of gauge metric type wasn't parsed correctly",
			saveParameters: services.MetricSaveParameters{
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
			saveParameters: services.MetricSaveParameters{
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
			saveParameters: services.MetricSaveParameters{
				MetricType:  "gauge",
				MetricName:  "metricName",
				MetricValue: "1.1",
			},
			testCase: func(t *testing.T, err error) {
				value, _ := storeMock.GetGaugeMetric("metricName")

				require.NoError(t, err)
				assert.Equal(t, 1.1, value)
			},
		},
		{
			testName: "Should save in storage counter metric type",
			saveParameters: services.MetricSaveParameters{
				MetricType:  "counter",
				MetricName:  "metricName",
				MetricValue: "100",
			},
			testCase: func(t *testing.T, err error) {
				value, _ := storeMock.GetCounterMetric("metricName")

				require.NoError(t, err)
				assert.Equal(t, int64(100), value)
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

func TestMetrics_SaveModel(t *testing.T) {
	var deltaMock int64 = 100
	var valueMock = 1.1

	storeMock := memstorage.NewWithPrefilledData(map[string]float64{}, map[string]int64{})
	metricsService := New(storeMock)

	testCases := []struct {
		testName       string
		saveParameters models.Metrics
		testCase       func(t *testing.T, err error)
	}{
		{
			testName: "Should return error if metricType isn't defined",
			saveParameters: models.Metrics{
				MType: "incorrect",
			},
			testCase: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			testName: "Should save in storage gauge metric type",
			saveParameters: models.Metrics{
				ID:    "metricName",
				MType: "gauge",
				Value: &valueMock,
			},
			testCase: func(t *testing.T, err error) {
				value, _ := storeMock.GetGaugeMetric("metricName")

				require.NoError(t, err)
				assert.Equal(t, 1.1, value)
			},
		},
		{
			testName: "Should save in storage counter metric type",
			saveParameters: models.Metrics{
				ID:    "metricName",
				MType: "counter",
				Delta: &deltaMock,
			},
			testCase: func(t *testing.T, err error) {
				value, _ := storeMock.GetCounterMetric("metricName")

				require.NoError(t, err)
				assert.Equal(t, int64(100), value)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			err := metricsService.SaveModel(tc.saveParameters)

			tc.testCase(t, err)
		})
	}
}

func TestMetrics_GetAll(t *testing.T) {
	metricsService := New(memstorage.NewWithPrefilledData(map[string]float64{"first": 1.11234}, map[string]int64{"second": 1}))

	t.Run("Should return all metrics", func(t *testing.T) {
		result := metricsService.GetAll()

		assert.Equal(t, 2, len(result))
		assert.Contains(t, result, services.MetricEntry{Name: "first", Value: "1.11234"})
		assert.Contains(t, result, services.MetricEntry{Name: "second", Value: "1"})
	})
}

func TestMetrics_Get(t *testing.T) {
	metricsService := New(memstorage.NewWithPrefilledData(map[string]float64{"first": 1.1}, map[string]int64{"second": 1}))

	testCases := []struct {
		testName      string
		getParameters services.MetricGetParameters
		expectedValue string
		expectedOk    bool
	}{
		{
			testName: "Should return gauge metricType value",
			getParameters: services.MetricGetParameters{
				MetricType: "gauge",
				MetricName: "first",
			},
			expectedValue: "1.1",
			expectedOk:    true,
		},
		{
			testName: "Should return counter metricType value",
			getParameters: services.MetricGetParameters{
				MetricType: "counter",
				MetricName: "second",
			},
			expectedValue: "1",
			expectedOk:    true,
		},
		{
			testName: "Should return nothing if store doesn't contain such value",
			getParameters: services.MetricGetParameters{
				MetricType: "test",
				MetricName: "test",
			},
			expectedValue: "",
			expectedOk:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			value, ok := metricsService.Get(services.MetricGetParameters{MetricName: tc.getParameters.MetricName, MetricType: tc.getParameters.MetricType})

			require.Equal(t, tc.expectedOk, ok)
			assert.Equal(t, tc.expectedValue, value)
		})
	}
}

func TestMetrics_GetModel(t *testing.T) {
	var deltaMock int64 = 1
	var valueMock = 1.1

	metricsService := New(memstorage.NewWithPrefilledData(map[string]float64{"first": valueMock}, map[string]int64{"second": deltaMock}))

	testCases := []struct {
		testName      string
		getParameters models.Metrics
		expectedValue models.Metrics
		expectedOk    bool
	}{
		{
			testName: "Should return gauge metricType value",
			getParameters: models.Metrics{
				MType: "gauge",
				ID:    "first",
			},
			expectedValue: models.Metrics{
				MType: "gauge",
				ID:    "first",
				Value: &valueMock,
			},
			expectedOk: true,
		},
		{
			testName: "Should return counter metricType value",
			getParameters: models.Metrics{
				MType: "counter",
				ID:    "second",
			},
			expectedValue: models.Metrics{
				MType: "counter",
				ID:    "second",
				Delta: &deltaMock,
			},
			expectedOk: true,
		},
		{
			testName: "Should return nothing if store doesn't contain such value",
			getParameters: models.Metrics{
				MType: "test",
				ID:    "test",
			},
			expectedValue: models.Metrics{
				MType: "counter",
				ID:    "second",
			},
			expectedOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			value, ok := metricsService.GetModel(tc.getParameters)

			require.Equal(t, tc.expectedOk, ok)

			if ok {
				assert.Equal(t, tc.expectedValue, value)
			}
		})
	}
}

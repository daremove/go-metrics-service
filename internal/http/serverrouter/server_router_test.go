package serverrouter

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/daremove/go-metrics-service/internal/middlewares/gzipm"
	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services"
	"github.com/daremove/go-metrics-service/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendMetricData(t *testing.T) {
	createServer := func(assertFn func(request *http.Request)) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assertFn(r)
		}))
	}
	testCases := []struct {
		testName   string
		testServer *httptest.Server
	}{
		{
			testName: "Should send request with correct path parameters",
			testServer: createServer(func(r *http.Request) {
				assert.Equal(t, "/update/metricType/metricName/metricValue", r.URL.Path)
			}),
		},
		{
			testName: "Should send request with correct http method",
			testServer: createServer(func(r *http.Request) {
				assert.Equal(t, "POST", r.Method)
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			defer tc.testServer.Close()

			err := SendMetricData(SendMetricDataParameters{
				URL:         tc.testServer.URL,
				MetricType:  "metricType",
				MetricName:  "metricName",
				MetricValue: "metricValue",
			})

			assert.NoError(t, err)
		})
	}
}

func TestSendMetricModelData(t *testing.T) {
	var deltaMock int64 = 1
	var valueMock = 2.5

	createServer := func(assertFn func(request *http.Request)) *httptest.Server {
		return httptest.NewServer(gzipm.GzipMiddleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			assertFn(r)
		})))
	}
	testCases := []struct {
		testName   string
		testServer *httptest.Server
	}{
		{
			testName: "Should send request with correct body",
			testServer: createServer(func(r *http.Request) {
				data, err := utils.DecodeJSONRequest[models.Metrics](r)

				require.NoError(t, err)
				assert.Equal(t, data, models.Metrics{
					ID:    "metricName",
					MType: "metricType",
					Delta: &deltaMock,
					Value: &valueMock,
				})
			}),
		},
		{
			testName: "Should send request with correct http method",
			testServer: createServer(func(r *http.Request) {
				assert.Equal(t, "POST", r.Method)
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			defer tc.testServer.Close()

			err := SendMetricModelData(tc.testServer.URL, models.Metrics{
				ID:    "metricName",
				MType: "metricType",
				Delta: &deltaMock,
				Value: &valueMock,
			})

			assert.NoError(t, err)
		})
	}
}

type metricsServiceMock struct {
	data      map[string]string
	modelData map[string]models.Metrics
}

func (m metricsServiceMock) Save(parameters services.MetricSaveParameters) error {
	return nil
}

func (m metricsServiceMock) SaveModel(parameters models.Metrics) error {
	return nil
}

func (m metricsServiceMock) Get(parameters services.MetricGetParameters) (string, bool) {
	value, ok := m.data[parameters.MetricName]

	return value, ok
}

func (m metricsServiceMock) GetModel(parameters models.Metrics) (models.Metrics, bool) {
	value, ok := m.modelData[parameters.ID]

	return value, ok
}

func (m metricsServiceMock) GetAll() []services.MetricEntry {
	var result []services.MetricEntry

	for key, value := range m.data {
		result = append(result, services.MetricEntry{Name: key, Value: value})
	}

	return result
}

type healthCheckServiceMock struct{}

func (hc healthCheckServiceMock) CheckStorageConnection() error {
	return nil
}

func TestServerRouter(t *testing.T) {
	testServer := httptest.NewServer(
		New(metricsServiceMock{
			data: map[string]string{
				"test": "1.1",
			},
		}, healthCheckServiceMock{}, "").Get(),
	)
	defer testServer.Close()

	testCases := []struct {
		testName        string
		methodName      string
		targetURL       string
		expectedCode    int
		expectedMessage string
	}{
		{
			testName:        "Should return 404 if metricName parameter wasn't provided",
			methodName:      http.MethodPost,
			targetURL:       "/update/metricType",
			expectedCode:    http.StatusNotFound,
			expectedMessage: "metricName wasn't provided\n",
		},
		{
			testName:        "Should return 400 if metricValue parameter wasn't provided",
			methodName:      http.MethodPost,
			targetURL:       "/update/metricType/metricName",
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "metricValue wasn't provided\n",
		},
		{
			testName:     "Should return 200",
			methodName:   http.MethodPost,
			targetURL:    "/update/metricType/metricName/100",
			expectedCode: http.StatusOK,
		},
		{
			testName:        "Should return metricValue",
			methodName:      http.MethodGet,
			targetURL:       "/value/gauge/test",
			expectedCode:    http.StatusOK,
			expectedMessage: "1.1",
		},
		{
			testName:        "Should return 404 if metricValue wasn't found",
			methodName:      http.MethodGet,
			targetURL:       "/value/gauge/another",
			expectedCode:    http.StatusNotFound,
			expectedMessage: "Metric value with such parameters wasn't found\n",
		},
		{
			testName:        "Should return html page with metrics",
			methodName:      http.MethodGet,
			targetURL:       "/",
			expectedCode:    http.StatusOK,
			expectedMessage: "<html><head><title>All metrics</title></head><body>test - 1.1</body></html>",
		},
		{
			testName:     "Should return 414 if appropriate content-type wasn't set for json handler",
			methodName:   http.MethodPost,
			targetURL:    "/update",
			expectedCode: http.StatusUnsupportedMediaType,
		},
	}

	t.Run("Should return 405 if http method isn't correct", func(testing *testing.T) {
		for _, methodName := range []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		} {
			res, _ := utils.TestRequest(t, testServer, methodName, "/update/metricType/metricName/123", nil, nil)
			res.Body.Close()

			assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
		}
	})

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			res, mes := utils.TestRequest(t, testServer, tc.methodName, tc.targetURL, nil, nil)
			res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)
			assert.Equal(t, tc.expectedMessage, mes)
		})
	}
}

func TestServerRouterJson(t *testing.T) {
	var deltaMock int64 = 1
	var valueMock = 2.5

	counterDataMock, _ := json.Marshal(models.Metrics{ID: "counter_test", MType: "counter", Delta: &deltaMock})
	gaugeDataMock, _ := json.Marshal(models.Metrics{ID: "gauge_test", MType: "gauge", Value: &valueMock})

	testServer := httptest.NewServer(
		New(metricsServiceMock{
			modelData: map[string]models.Metrics{
				"gauge_test": {ID: "test", MType: "gauge", Value: &valueMock},
			},
		}, healthCheckServiceMock{}, "").Get(),
	)
	defer testServer.Close()

	testCases := []struct {
		testName        string
		methodName      string
		targetURL       string
		expectedCode    int
		expectedMessage string
		headers         map[string]string
		body            io.Reader
	}{
		{
			testName:     "Should return 414 if appropriate content-type wasn't set for json handler",
			methodName:   http.MethodPost,
			targetURL:    "/update",
			expectedCode: http.StatusUnsupportedMediaType,
		},
		{
			testName:        "Should save correctly metric data by using model",
			methodName:      http.MethodPost,
			targetURL:       "/update",
			expectedCode:    http.StatusOK,
			expectedMessage: "{\"id\":\"counter_test\",\"type\":\"counter\",\"delta\":1}",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: bytes.NewBuffer(counterDataMock),
		},
		{
			testName:        "Should return found data by using model",
			methodName:      http.MethodPost,
			targetURL:       "/value",
			expectedCode:    http.StatusOK,
			expectedMessage: "{\"id\":\"test\",\"type\":\"gauge\",\"value\":2.5}",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: bytes.NewBuffer(gaugeDataMock),
		},
		{
			testName:        "Should return 404 if data wasn't found by using model",
			methodName:      http.MethodPost,
			targetURL:       "/value",
			expectedCode:    http.StatusNotFound,
			expectedMessage: "Metric value with such parameters wasn't found\n",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: bytes.NewBuffer(counterDataMock),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			res, mes := utils.TestRequest(t, testServer, tc.methodName, tc.targetURL, tc.headers, tc.body)
			res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)
			assert.Equal(t, tc.expectedMessage, mes)
		})
	}
}

func TestServerRouterGzip(t *testing.T) {
	var valueMock = 1.1
	testServer := httptest.NewServer(
		New(metricsServiceMock{
			data: map[string]string{
				"test": "1.1",
			},
			modelData: map[string]models.Metrics{
				"test": {ID: "test", MType: "gauge", Value: &valueMock},
			},
		}, healthCheckServiceMock{}, "").Get(),
	)
	defer testServer.Close()

	t.Run("Should return gzip data", func(t *testing.T) {
		res, mes := utils.TestRequest(t, testServer, http.MethodGet, "/", map[string]string{
			"Content-Type":    "text/html",
			"Accept-Encoding": "gzip",
		}, nil)
		err := res.Body.Close()

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		reader, err := utils.NewCompressReader(io.NopCloser(strings.NewReader(mes)))

		require.NoError(t, err)

		result, err := io.ReadAll(reader)

		require.NoError(t, err)

		assert.Equal(t, "<html><head><title>All metrics</title></head><body>test - 1.1</body></html>", string(result))
	})

	t.Run("Should accept gzip data", func(t *testing.T) {
		body, err := json.Marshal(models.Metrics{
			ID:    "test",
			MType: "gauge",
			Value: &valueMock,
		})

		require.NoError(t, err)

		var buf bytes.Buffer

		gzipWriter := gzip.NewWriter(&buf)
		_, err = gzipWriter.Write(body)

		require.NoError(t, err)

		err = gzipWriter.Close()

		require.NoError(t, err)

		body = buf.Bytes()
		res, mes := utils.TestRequest(t, testServer, http.MethodPost, "/value", map[string]string{
			"Content-Type":     "application/json",
			"Content-Encoding": "gzip",
		}, bytes.NewBuffer(body))
		err = res.Body.Close()

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		assert.Equal(t, "{\"id\":\"test\",\"type\":\"gauge\",\"value\":1.1}", mes)
	})
}

package serverrouter

import (
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

type metricsServiceMock struct {
	data map[string]string
}

func (m metricsServiceMock) Save(parameters metrics.SaveParameters) error {
	return nil
}

func (m metricsServiceMock) Get(parameters metrics.GetParameters) (string, bool) {
	value, ok := m.data[parameters.MetricName]

	return value, ok
}

func (m metricsServiceMock) GetAll() []metrics.MetricItem {
	var result []metrics.MetricItem

	for key, value := range m.data {
		result = append(result, metrics.MetricItem{Name: key, Value: value})
	}

	return result
}

func TestServerRouter(t *testing.T) {
	testServer := httptest.NewServer(ServerRouter(metricsServiceMock{
		data: map[string]string{
			"test": "1.1",
		},
	}))
	defer testServer.Close()

	testCases := []struct {
		testName        string
		methodName      string
		targetURL       string
		expectedCode    int
		expectedMessage string
	}{
		{
			testName:        "Should return 400 if metricType parameter wasn't provided",
			methodName:      http.MethodPost,
			targetURL:       "/update",
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "metricType wasn't provided\n",
		},
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
			res, _ := utils.TestRequest(t, testServer, methodName, "/update/metricType/metricName/123")
			res.Body.Close()

			assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
		}
	})

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			res, mes := utils.TestRequest(t, testServer, tc.methodName, tc.targetURL)
			res.Body.Close()

			assert.Equal(t, tc.expectedCode, res.StatusCode)
			assert.Equal(t, tc.expectedMessage, mes)
		})
	}
}

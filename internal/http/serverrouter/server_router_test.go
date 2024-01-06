package serverrouter

import (
	"github.com/daremove/go-metrics-service/internal/services/metrics"
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

type metricsServiceMock struct{}

func (m metricsServiceMock) Save(parameters metrics.SaveParameters) error {
	return nil
}

func TestUpdateMetricHandler(t *testing.T) {
	testCases := []struct {
		testName     string
		methodName   string
		targetURL    string
		expectedCode int
	}{
		{
			testName:     "Should return 400 if metricType parameter wasn't provided",
			methodName:   http.MethodPost,
			targetURL:    "/update",
			expectedCode: http.StatusBadRequest,
		},
		{
			testName:     "Should return 404 if metricName parameter wasn't provided",
			methodName:   http.MethodPost,
			targetURL:    "/update/metricType",
			expectedCode: http.StatusNotFound,
		},
		{
			testName:     "Should return 400 if metricValue parameter wasn't provided",
			methodName:   http.MethodPost,
			targetURL:    "/update/metricType/metricName",
			expectedCode: http.StatusBadRequest,
		},
		{
			testName:     "Should return 200",
			methodName:   http.MethodPost,
			targetURL:    "/update/metricType/metricName/100",
			expectedCode: http.StatusOK,
		},
	}

	t.Run("Should return 404 if http method isn't correct", func(testing *testing.T) {
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
			r := httptest.NewRequest(methodName, "/update", nil)
			w := httptest.NewRecorder()

			updateMetricHandler(metricsServiceMock{})(w, r)

			assert.Equal(t, http.StatusNotFound, w.Code)
		}
	})

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			r := httptest.NewRequest(tc.methodName, tc.targetURL, nil)
			w := httptest.NewRecorder()

			updateMetricHandler(metricsServiceMock{})(w, r)

			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

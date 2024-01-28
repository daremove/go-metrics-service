package gzipm

import (
	"github.com/daremove/go-metrics-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func TestGzipMiddleware(t *testing.T) {
	router := chi.NewRouter()
	router.Use(GzipMiddleware)
	router.HandleFunc("/", dummyHandler)

	testServer := httptest.NewServer(
		router,
	)
	defer testServer.Close()

	testCases := []struct {
		testName                string
		expectedContentEncoding string
		headers                 map[string]string
	}{
		{
			testName: "Should accept application/json content type",
			headers: map[string]string{
				"Content-type":    "application/json",
				"Accept-Encoding": "gzip",
			},
			expectedContentEncoding: "gzip",
		},
		{
			testName: "Should accept text/html content type",
			headers: map[string]string{
				"Content-type":    "text/html",
				"Accept-Encoding": "gzip",
			},
			expectedContentEncoding: "gzip",
		},
		{
			testName: "Should not compress data without accept encoding",
			headers: map[string]string{
				"Content-type": "text/html",
			},
			expectedContentEncoding: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			res, _ := utils.TestRequest(t, testServer, http.MethodGet, "/", tc.headers, nil)
			err := res.Body.Close()

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tc.expectedContentEncoding, res.Header.Get("Content-Encoding"))
		})
	}
}

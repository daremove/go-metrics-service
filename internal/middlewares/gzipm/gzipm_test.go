package gzipm

import (
	"bytes"
	"compress/gzip"
	"github.com/daremove/go-metrics-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := io.ReadAll(r.Body)

	w.Write(resp)
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
		testName        string
		rawBody         []byte
		headers         map[string]string
		expectedMessage string
	}{
		{
			testName: "Should decompress request if Content-Encoding is set",
			headers: map[string]string{
				"Content-Encoding": "gzip",
			},
			rawBody:         []byte("test"),
			expectedMessage: "test",
		},
		{
			testName:        "Shouldn't decompress request if Content-Encoding isn't set",
			rawBody:         []byte("test"),
			expectedMessage: "test",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			var body []byte

			if tc.headers["Content-Encoding"] == "gzip" {
				var buf bytes.Buffer

				gzipWriter := gzip.NewWriter(&buf)
				_, err := gzipWriter.Write(tc.rawBody)

				require.NoError(t, err)

				err = gzipWriter.Close()

				require.NoError(t, err)

				body = buf.Bytes()
			} else {
				body = tc.rawBody
			}

			res, mes := utils.TestRequest(t, testServer, http.MethodPost, "/", tc.headers, bytes.NewBuffer(body))
			err := res.Body.Close()

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, tc.expectedMessage, mes)
		})
	}
}

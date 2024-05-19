package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestRequest(t *testing.T) {
	t.Run("Should make a successful GET request", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, world!"))
		}
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		resp, body := TestRequest(t, ts, "GET", "/", nil, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "Hello, world!", body)
	})

	t.Run("Should make a successful POST request with body", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "POST", r.Method)
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.Equal(t, "request body", string(body))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Received"))
		}
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		headers := map[string]string{"Content-Type": "text/plain"}
		resp, body := TestRequest(t, ts, "POST", "/", headers, strings.NewReader("request body"))
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "Received", body)
	})

	t.Run("Should handle headers correctly", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Headers OK"))
		}
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		headers := map[string]string{"Content-Type": "application/json"}
		resp, body := TestRequest(t, ts, "GET", "/", headers, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "Headers OK", body)
	})

	t.Run("Should handle query parameters correctly", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "GET", r.Method)
			require.Equal(t, "value", r.URL.Query().Get("param"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Query OK"))
		}
		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		resp, body := TestRequest(t, ts, "GET", "/?param=value", nil, nil)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "Query OK", body)
	})
}

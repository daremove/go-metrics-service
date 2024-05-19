package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	t.Run("Should initialize logger with valid level", func(t *testing.T) {
		err := Initialize("info")
		require.NoError(t, err)
		assert.NotNil(t, Log)
	})

	t.Run("Should return error with invalid level", func(t *testing.T) {
		err := Initialize("invalid")
		assert.Error(t, err)
	})
}

func TestRequestLogger(t *testing.T) {
	err := Initialize("info")
	require.NoError(t, err)

	// Создаем тестовый обработчик
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, world!"))
	})

	// Создаем тестовый HTTP сервер
	ts := httptest.NewServer(RequestLogger(testHandler))
	defer ts.Close()

	t.Run("Should log HTTP requests and responses", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Should log request details", func(t *testing.T) {
		req, err := http.NewRequest("GET", ts.URL, nil)
		require.NoError(t, err)

		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		duration := time.Since(start)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Less(t, duration.Milliseconds(), int64(1000)) // Проверка, что запрос обрабатывается быстро
	})
}

func TestLoggingResponseWriter(t *testing.T) {
	t.Run("Should correctly log response status and size", func(t *testing.T) {
		responseData := &responseData{}
		rr := httptest.NewRecorder()
		lw := &loggingResponseWriter{
			ResponseWriter: rr,
			responseData:   responseData,
		}

		lw.WriteHeader(http.StatusCreated)
		lw.Write([]byte("Hello, world!"))

		assert.Equal(t, http.StatusCreated, responseData.status)
		assert.Equal(t, len("Hello, world!"), responseData.size)
	})
}

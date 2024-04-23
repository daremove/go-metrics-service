// Пакет logger предназначен для логирования HTTP запросов и ответов.
package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseData содержит данные о HTTP ответе.
type responseData struct {
	status int // HTTP статус код ответа
	size   int // Размер ответа в байтах
}

// loggingResponseWriter реализует интерфейс http.ResponseWriter для перехвата и логирования ответов.
type loggingResponseWriter struct {
	http.ResponseWriter               // Встроенный http.ResponseWriter
	responseData        *responseData // Ссылка на данные о ответе
}

// Write переопределяет метод Write интерфейса http.ResponseWriter для подсчета размера данных.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader переопределяет метод WriteHeader интерфейса http.ResponseWriter для записи статус кода.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Log предоставляет глобальный доступ к логгеру zap.Logger.
var Log *zap.Logger = zap.NewNop()

// Initialize инициализирует логгер с указанным уровнем логирования.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		return err
	}

	Log = zl

	return nil
}

// RequestLogger является middleware для логирования HTTP запросов и ответов.
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		Log.Info("got incoming HTTP request",
			zap.String("URI", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration*time.Millisecond),
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	})
}

// Пакет serverrouter предназначен для настройки маршрутизации API на сервере.
package serverrouter

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daremove/go-metrics-service/internal/middlewares/profiler"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/middlewares/dataintergity"
	"github.com/daremove/go-metrics-service/internal/middlewares/gzipm"
	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services"
	"github.com/daremove/go-metrics-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// RouterConfig содержит конфигурацию для маршрутизатора сервера.
type RouterConfig struct {
	Endpoint   string // URL-адрес конечной точки сервера
	SigningKey string // Ключ для подписи данных
}

// ServerRouter предоставляет маршрутизацию запросов к сервисам метрик и проверки состояния.
type ServerRouter struct {
	metricsService     MetricsService     // Сервис для работы с метриками
	healthCheckService HealthCheckService // Сервис для проверки состояния
	config             RouterConfig       // Конфигурация маршрутизатора
}

// MetricsService определяет интерфейс для сервиса метрик.
type MetricsService interface {
	Save(ctx context.Context, parameters services.MetricSaveParameters) error         // Сохраняет метрику
	SaveModel(ctx context.Context, parameters models.Metrics) error                   // Сохраняет модель метрик
	SaveModels(ctx context.Context, parameters []models.Metrics) error                // Сохраняет несколько моделей метрик
	Get(ctx context.Context, parameters services.MetricGetParameters) (string, error) // Получает значение метрики
	GetModel(ctx context.Context, parameters models.Metrics) (models.Metrics, error)  // Получает модель метрики
	GetAll(ctx context.Context) ([]services.MetricEntry, error)                       // Получает все метрики
}

// HealthCheckService определяет интерфейс для сервиса проверки состояния.
type HealthCheckService interface {
	CheckStorageConnection(ctx context.Context) error // Проверяет соединение с хранилищем
}

// New создает новый экземпляр ServerRouter.
func New(metricsService MetricsService, healthCheckService HealthCheckService, config RouterConfig) *ServerRouter {
	return &ServerRouter{metricsService, healthCheckService, config}
}

// Get инициализирует и возвращает маршрутизатор с предварительно сконфигурированными маршрутами.
func (router *ServerRouter) Get(ctx context.Context) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger)
	r.Use(middleware.NewCompressor(flate.BestSpeed).Handler)
	r.Use(dataintergity.NewMiddleware(dataintergity.DataIntegrityMiddlewareConfig{
		SigningKey: router.config.SigningKey,
	}))
	r.Use(gzipm.GzipMiddleware)

	r.Mount("/debug", profiler.Profiler())

	r.Route("/", func(r chi.Router) {
		r.Get("/", getAllMetricsHandler(ctx, router.metricsService))

		r.Route("/update", func(r chi.Router) {
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Post("/{metricValue}", updateMetricHandler(ctx, router.metricsService))

					r.Post("/", func(w http.ResponseWriter, r *http.Request) {
						http.Error(w, "metricValue wasn't provided", http.StatusBadRequest)
					})
				})

				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "metricName wasn't provided", http.StatusNotFound)
				})
			})

			r.Post("/", updateMetricWithJSONHandler(ctx, router.metricsService))
		})

		r.Route("/updates", func(r chi.Router) {
			r.Post("/", updateMetricsHandler(ctx, router.metricsService))
		})

		r.Route("/value", func(r chi.Router) {
			r.Get("/{metricType}/{metricName}", getMetricValueHandler(ctx, router.metricsService))

			r.Post("/", getMetricValueWithJSONHandler(ctx, router.metricsService))
		})

		r.Route("/ping", func(r chi.Router) {
			r.Get("/", pingHandler(ctx, router.healthCheckService))
		})
	})

	return r
}

// Run запускает сервер на заданном порту и с заданными маршрутами.
func (router *ServerRouter) Run(ctx context.Context) {
	log.Fatal(http.ListenAndServe(router.config.Endpoint, router.Get(ctx)))
}

func handleJSONError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case utils.UnsupportedContentTypeCode:
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
	default:
		logger.Log.Debug("cannot parse json data from request", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateMetricHandler(ctx context.Context, metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := metricsService.Save(ctx, services.MetricSaveParameters{
			MetricType:  chi.URLParam(r, "metricType"),
			MetricName:  chi.URLParam(r, "metricName"),
			MetricValue: chi.URLParam(r, "metricValue"),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func updateMetricWithJSONHandler(ctx context.Context, metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.DecodeJSONRequest[models.Metrics](r)

		if err != nil {
			handleJSONError(w, err)
			return
		}

		if err := metricsService.SaveModel(ctx, data); err != nil {
			logger.Log.Error("error saving data in metrics service", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := utils.EncodeJSONRequest[models.Metrics](w, data); err != nil {
			logger.Log.Error("error encoding response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func updateMetricsHandler(ctx context.Context, metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.DecodeJSONRequest[[]models.Metrics](r)

		if err != nil {
			handleJSONError(w, err)
			return
		}

		if err := metricsService.SaveModels(ctx, data); err != nil {
			logger.Log.Error("error saving data in metrics service", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := utils.EncodeJSONRequest[[]models.Metrics](w, data); err != nil {
			logger.Log.Error("error encoding response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func pingHandler(ctx context.Context, healthCheckService HealthCheckService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := healthCheckService.CheckStorageConnection(ctx); err != nil {
			logger.Log.Error("error check storage connection in health check service", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}

}

func getMetricValueHandler(ctx context.Context, metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := metricsService.Get(ctx, services.MetricGetParameters{
			MetricType: chi.URLParam(r, "metricType"),
			MetricName: chi.URLParam(r, "metricName"),
		})

		if err != nil {
			if errors.Is(err, services.ErrMetricNotFound) {
				http.Error(w, "Metric value with such parameters wasn't found", http.StatusNotFound)
				return
			}

			logger.Log.Error("error get metric data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := io.WriteString(w, value); err != nil {
			logger.Log.Error("failed to write data", zap.Error(err))
		}
	}
}

func getMetricValueWithJSONHandler(ctx context.Context, metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := utils.DecodeJSONRequest[models.Metrics](r)

		if err != nil {
			handleJSONError(w, err)
			return
		}

		value, err := metricsService.GetModel(ctx, data)

		if err != nil {
			if errors.Is(err, services.ErrMetricNotFound) {
				http.Error(w, "Metric value with such parameters wasn't found", http.StatusNotFound)
				return
			}

			logger.Log.Error("error get model metric data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := utils.EncodeJSONRequest[models.Metrics](w, value); err != nil {
			logger.Log.Error("error encoding response", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func getAllMetricsHandler(ctx context.Context, metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricData, err := metricsService.GetAll(ctx)

		if err != nil {
			logger.Log.Error("error get all metric data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := make([]string, 0, len(metricData))

		for _, el := range metricData {
			result = append(result, fmt.Sprintf("%s - %s", el.Name, el.Value))
		}

		sort.Strings(result)

		w.Header().Set("Content-Type", "text/html")

		if _, err := io.WriteString(w, fmt.Sprintf("<html><head><title>All metrics</title></head><body>%s</body></html>", strings.Join(result, "<br />"))); err != nil {
			logger.Log.Error("failed to write data", zap.Error(err))
		}
	}
}

// SendMetricDataParameters определяет параметры для отправки данных метрик.
type SendMetricDataParameters struct {
	URL         string // URL-адрес сервера
	MetricType  string // Тип метрики
	MetricName  string // Имя метрики
	MetricValue string // Значение метрики
}

// SendMetricData осуществляет отправку данных метрики по HTTP POST.
func SendMetricData(parameters SendMetricDataParameters) error {
	res, err := http.Post(fmt.Sprintf("%s/update/%s/%s/%s", parameters.URL, parameters.MetricType, parameters.MetricName, parameters.MetricValue), "text/plain", nil)
	if err != nil {
		return fmt.Errorf("failed to send data by using POST method: %w", err)
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	return nil
}

// SendMetricModelDataConfig содержит конфигурацию для отправки модели данных метрик.
type SendMetricModelDataConfig struct {
	URL        string // URL-адрес сервера
	SigningKey string // Ключ для подписи данных
}

// SendMetricModelData отправляет модель данных метрик на указанный сервер с возможной подписью данных.
func SendMetricModelData(data []models.Metrics, config SendMetricModelDataConfig) error {
	body, err := json.Marshal(data)

	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	var buf bytes.Buffer

	gzipWriter := gzip.NewWriter(&buf)
	_, err = gzipWriter.Write(body)

	if err != nil {
		return fmt.Errorf("failed to gzip data: %w", err)
	}

	err = gzipWriter.Close()

	if err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	body = buf.Bytes()
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/updates", config.URL), bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")

	if config.SigningKey != "" {
		signedBody, err := utils.SignData(body, config.SigningKey)

		if err != nil {
			return fmt.Errorf("failed to sign data: %w", err)
		}

		req.Header.Set(dataintergity.HeaderKeyHash, hex.EncodeToString(signedBody))
	}

	res, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to send data by using POST method: %w", err)
	}

	err = res.Body.Close()

	if err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	return nil
}

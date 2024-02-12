package serverrouter

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/middlewares/gzipm"
	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services"
	"github.com/daremove/go-metrics-service/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

type ServerRouter struct {
	metricsService     MetricsService
	healthCheckService HealthCheckService
	endpoint           string
}

type MetricsService interface {
	Save(ctx context.Context, parameters services.MetricSaveParameters) error

	SaveModel(ctx context.Context, parameters models.Metrics) error

	SaveModels(ctx context.Context, parameters []models.Metrics) error

	Get(ctx context.Context, parameters services.MetricGetParameters) (string, error)

	GetModel(ctx context.Context, parameters models.Metrics) (models.Metrics, error)

	GetAll(ctx context.Context) ([]services.MetricEntry, error)
}

type HealthCheckService interface {
	CheckStorageConnection(ctx context.Context) error
}

func New(metricsService MetricsService, healthCheckService HealthCheckService, endpoint string) *ServerRouter {
	return &ServerRouter{metricsService, healthCheckService, endpoint}
}

func (router *ServerRouter) Get(ctx context.Context) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.RequestLogger)
	r.Use(middleware.NewCompressor(flate.DefaultCompression).Handler)
	r.Use(gzipm.GzipMiddleware)

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

func (router *ServerRouter) Run(ctx context.Context) {
	log.Fatal(http.ListenAndServe(router.endpoint, router.Get(ctx)))
}

func handleJSONError(w http.ResponseWriter, err error) {
	switch err.Error() {
	case utils.UnsupportedContentTypeCode:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	default:
		logger.Log.Debug("cannot parse json data from request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
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
		var result []string

		metricData, err := metricsService.GetAll(ctx)

		if err != nil {
			logger.Log.Error("error get all metric data", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

type SendMetricDataParameters struct {
	URL         string
	MetricType  string
	MetricName  string
	MetricValue string
}

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

func SendMetricModelData(url string, data []models.Metrics) error {
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
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/updates", url), bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Encoding", "gzip")

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

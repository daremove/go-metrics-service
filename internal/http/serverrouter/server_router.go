package serverrouter

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

type ServerRouter struct {
	metricsService MetricsService
	endpoint       string
}

type MetricsService interface {
	Save(parameters services.MetricSaveParameters) error

	Get(parameters services.MetricGetParameters) (string, bool)

	GetAll() []services.MetricEntry
}

func New(metricsService MetricsService, endpoint string) *ServerRouter {
	return &ServerRouter{metricsService, endpoint}
}

func (router *ServerRouter) Get() chi.Router {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)

	r.Route("/", func(r chi.Router) {
		r.Get("/", getAllMetricsHandler(router.metricsService))

		r.Route("/update", func(r chi.Router) {
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Post("/{metricValue}", updateMetricHandler(router.metricsService))

					r.Post("/", func(w http.ResponseWriter, r *http.Request) {
						http.Error(w, "metricValue wasn't provided", http.StatusBadRequest)
					})
				})

				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "metricName wasn't provided", http.StatusNotFound)
				})
			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "metricType wasn't provided", http.StatusBadRequest)
			})
		})

		r.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
			r.Get("/", getMetricValueHandler(router.metricsService))
		})
	})

	return r
}

func (router *ServerRouter) Run() {
	log.Fatal(http.ListenAndServe(router.endpoint, router.Get()))
}

func updateMetricHandler(metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := metricsService.Save(services.MetricSaveParameters{
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

func getMetricValueHandler(metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, ok := metricsService.Get(services.MetricGetParameters{
			MetricType: chi.URLParam(r, "metricType"),
			MetricName: chi.URLParam(r, "metricName"),
		})

		if !ok {
			http.Error(w, "Metric value with such parameters wasn't found", http.StatusNotFound)
		}

		if _, err := io.WriteString(w, value); err != nil {
			logger.Log.Error("failed to write data", zap.Error(err))
		}
	}
}

func getAllMetricsHandler(metricsService MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var result []string

		for _, el := range metricsService.GetAll() {
			result = append(result, fmt.Sprintf("%s - %s", el.Name, el.Value))
		}

		sort.Strings(result)

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

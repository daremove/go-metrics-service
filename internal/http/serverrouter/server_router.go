package serverrouter

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
)

type serverRouter struct {
	metricsService metrics.Service
	endpoint       string
}

func New(metricsService metrics.Service, endpoint string) *serverRouter {
	return &serverRouter{metricsService, endpoint}
}

func ServerRouter(metricsService metrics.Service) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/", getAllMetricsHandler(metricsService))

		r.Route("/update", func(r chi.Router) {
			r.Route("/{metricType}", func(r chi.Router) {
				r.Route("/{metricName}", func(r chi.Router) {
					r.Post("/{metricValue}", updateMetricHandler(metricsService))

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
			r.Get("/", getMetricValueHandler(metricsService))
		})
	})

	return r
}

func (router *serverRouter) Run() {
	log.Fatal(http.ListenAndServe(router.endpoint, ServerRouter(router.metricsService)))
}

func updateMetricHandler(metricsService metrics.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := metricsService.Save(metrics.SaveParameters{
			MetricType:  chi.URLParam(r, "metricType"),
			MetricName:  chi.URLParam(r, "metricName"),
			MetricValue: chi.URLParam(r, "metricValue"),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func getMetricValueHandler(metricsService metrics.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, ok := metricsService.Get(metrics.GetParameters{
			MetricType: chi.URLParam(r, "metricType"),
			MetricName: chi.URLParam(r, "metricName"),
		})

		if !ok {
			http.Error(w, "Metric value with such parameters wasn't found", http.StatusNotFound)
		}

		if _, err := io.WriteString(w, value); err != nil {
			panic(err)
		}
	}
}

func getAllMetricsHandler(metricsService metrics.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var result []string

		for _, el := range metricsService.GetAll() {
			result = append(result, fmt.Sprintf("%s - %s", el.Name, el.Value))
		}

		sort.Strings(result)

		_, err := io.WriteString(w, fmt.Sprintf("<html><head><title>All metrics</title></head><body>%s</body></html>", strings.Join(result, "<br />")))

		if err != nil {
			panic(err)
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
		return err
	}

	err = res.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

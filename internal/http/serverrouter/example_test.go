package serverrouter

import (
	"context"
	"log"
	"net/http"

	"github.com/daremove/go-metrics-service/internal/services/healthcheck"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/database"
	"github.com/go-chi/chi/v5"
)

func Example() {
	ctx := context.Background()
	router := chi.NewRouter()

	serverAddress := "localhost:8080"
	dsn := "postgresql://localhost/dbname"
	storage, _ := database.New(ctx, dsn)

	metricsService := metrics.New(storage)
	healthCheckService := healthcheck.New(storage)

	router.Get("/", getAllMetricsHandler(ctx, metricsService))

	router.Post("/update/{metricType}/{metricName}/{metricValue}", updateMetricHandler(ctx, metricsService))
	router.Post("/update", updateMetricWithJSONHandler(ctx, metricsService))
	router.Post("/updates", updateMetricsHandler(ctx, metricsService))

	router.Get("/value/{metricType}/{metricName}", getMetricValueHandler(ctx, metricsService))
	router.Post("/value", getMetricValueWithJSONHandler(ctx, metricsService))

	router.Get("/ping", pingHandler(ctx, healthCheckService))

	log.Fatal(http.ListenAndServe(serverAddress, router))
}

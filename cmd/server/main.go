package main

import (
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
)

func main() {
	store := memstorage.New()
	metricsService := metrics.New(store)
	router := serverrouter.New(metricsService, 8080)

	router.Run()
}

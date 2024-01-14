package main

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
)

func main() {
	parseFlags()

	store := memstorage.New()
	metricsService := metrics.New(store)
	router := serverrouter.New(metricsService, endpoint)

	fmt.Printf("Running server on %s\n", endpoint)

	router.Run()
}

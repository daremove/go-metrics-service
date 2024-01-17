package main

import (
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"log"
)

func main() {
	config := NewConfig()

	store := memstorage.New()
	metricsService := metrics.New(store)
	router := serverrouter.New(metricsService, config.endpoint)

	log.Printf("Running server on %s\n", config.endpoint)

	router.Run()
}

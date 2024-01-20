package main

import (
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"log"
)

func main() {
	config := NewConfig()

	if err := logger.Initialize("info"); err != nil {
		log.Fatalf("Logger wasn't initialized due to %s", err)
	}

	store := memstorage.New()
	metricsService := metrics.New(store)
	router := serverrouter.New(metricsService, config.endpoint)

	log.Printf("Running server on %s\n", config.endpoint)

	router.Run()
}

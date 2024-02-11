package main

import (
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services/filestorage"
	"github.com/daremove/go-metrics-service/internal/services/healthcheck"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/database"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"github.com/daremove/go-metrics-service/internal/utils"
	"log"
)

func main() {
	config := NewConfig()

	if err := logger.Initialize(config.logLevel); err != nil {
		log.Fatalf("Logger wasn't initialized due to %s", err)
	}

	store := memstorage.New()
	db := database.New(config.dsn)
	fileStorage, err := filestorage.New(store, filestorage.Config{
		StoreInterval:   config.storeInterval,
		FileStoragePath: config.fileStoragePath,
		Restore:         config.restore,
	})

	if err != nil {
		log.Fatalf("Backup service wasn't initialized due to %s", err)
	}

	metricsService := metrics.New(fileStorage)
	healthCheckService := healthcheck.New(db)
	router := serverrouter.New(metricsService, healthCheckService, config.endpoint)

	utils.HandleTerminationProcess(func() {
		if err := fileStorage.BackupData(); err != nil {
			log.Fatalf("Cannot backup data data after termination process %s", err)
		}
	})

	log.Printf("Running server on %s\n", config.endpoint)

	router.Run()
}

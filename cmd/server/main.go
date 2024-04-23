package main

import (
	"context"
	"log"

	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services/filestorage"
	"github.com/daremove/go-metrics-service/internal/services/healthcheck"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/database"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"github.com/daremove/go-metrics-service/internal/utils"
)

func main() {
	ctx := context.Background()
	config := NewConfig()

	if err := logger.Initialize(config.logLevel); err != nil {
		log.Fatalf("Logger wasn't initialized due to %s", err)
	}

	var (
		storage            metrics.Storage
		healthCheckService *healthcheck.HealthCheck
	)

	if config.dsn == "" {
		fileStorage, err := filestorage.New(ctx, memstorage.New(), filestorage.Config{
			StoreInterval:   config.storeInterval,
			FileStoragePath: config.fileStoragePath,
			Restore:         config.restore,
		})

		if err != nil {
			log.Fatalf("Backup service wasn't initialized due to %s", err)
		}

		storage = fileStorage
		healthCheckService = healthcheck.New(nil)

		utils.HandleTerminationProcess(func() {
			if err := fileStorage.BackupData(ctx); err != nil {
				log.Fatalf("Cannot backup data data after termination process %s", err)
			}
		})
	} else {
		db, err := database.New(ctx, config.dsn)

		if err != nil {
			log.Fatalf("Database wasn't initialized due to %s", err)
		}

		storage = db
		healthCheckService = healthcheck.New(db)
	}

	metricsService := metrics.New(storage)
	router := serverrouter.New(metricsService, healthCheckService, serverrouter.RouterConfig{
		Endpoint:   config.endpoint,
		SigningKey: config.signingKey,
	})

	log.Printf("Running server on %s\n", config.endpoint)

	router.Run(ctx)
}

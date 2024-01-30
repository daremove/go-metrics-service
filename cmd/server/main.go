package main

import (
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services/backup"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"github.com/daremove/go-metrics-service/internal/utils"
	"log"
)

func main() {
	config := NewConfig()

	if err := logger.Initialize("error"); err != nil {
		log.Fatalf("Logger wasn't initialized due to %s", err)
	}

	store := memstorage.New()
	backupService, err := backup.New(store, backup.Config{
		StoreInterval:   config.storeInterval,
		FileStoragePath: config.fileStoragePath,
		Restore:         config.restore,
	})

	if err != nil {
		log.Fatalf("Backup service wasn't initialized due to %s", err)
	}

	metricsService := metrics.New(backupService.FileStorage)
	router := serverrouter.New(metricsService, config.endpoint)

	utils.HandleTerminationProcess(func() {
		if err := backupService.BackupData(); err != nil {
			log.Fatalf("Cannot backup data data after termination process %s", err)
		}
	})

	log.Printf("Running server on %s\n", config.endpoint)

	router.Run()
}

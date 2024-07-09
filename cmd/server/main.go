package main

import (
	"context"
	"crypto/rsa"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/daremove/go-metrics-service/internal/proto"
	"google.golang.org/grpc"

	_ "github.com/daremove/go-metrics-service/cmd/buildversion"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services/filestorage"
	"github.com/daremove/go-metrics-service/internal/services/healthcheck"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/storage/database"
	"github.com/daremove/go-metrics-service/internal/storage/memstorage"
	"github.com/daremove/go-metrics-service/internal/utils"

	pb "github.com/daremove/go-metrics-service/internal/proto/metrics"
)

func initializeLogger(logLevel string) error {
	return logger.Initialize(logLevel)
}

func initializeStorage(ctx context.Context, config Config) (metrics.Storage, *healthcheck.HealthCheck, error) {
	var storage metrics.Storage
	var healthCheckService *healthcheck.HealthCheck

	if config.Dsn == "" {
		fileStorage, err := filestorage.New(ctx, memstorage.New(), filestorage.Config{
			StoreInterval:   config.StoreInterval,
			FileStoragePath: config.FileStoragePath,
			Restore:         config.Restore,
		})

		if err != nil {
			return nil, nil, err
		}

		storage = fileStorage
		healthCheckService = healthcheck.New(nil)

		utils.HandleTerminationProcess(func() {
			if err := fileStorage.BackupData(ctx); err != nil {
				log.Fatalf("Cannot backup data after termination process %s", err)
			}
		})
	} else {
		db, err := database.New(ctx, config.Dsn)

		if err != nil {
			return nil, nil, err
		}

		storage = db
		healthCheckService = healthcheck.New(db)
	}

	return storage, healthCheckService, nil
}

func runServer(ctx context.Context, config Config, metricsService *metrics.Metrics, healthCheckService *healthcheck.HealthCheck, privateKey *rsa.PrivateKey) *http.Server {

	router := serverrouter.New(metricsService, healthCheckService, serverrouter.RouterConfig{
		Endpoint:      config.Endpoint,
		SigningKey:    config.SigningKey,
		PrivateKey:    privateKey,
		TrustedSubnet: config.TrustedSubnet,
	})

	server := &http.Server{
		Addr:    config.Endpoint,
		Handler: router.Get(ctx),
	}

	go func() {
		log.Printf("Running server on %s\n", config.Endpoint)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", config.Endpoint, err)
		}
	}()

	return server
}

func runGRPCServer(metricsService *metrics.Metrics) *grpc.Server {
	address := ":3200"
	server := grpc.NewServer()
	metricsServer := proto.NewMetricsServer(metricsService)

	go func() {
		listen, err := net.Listen("tcp", address)

		if err != nil {
			log.Fatal(err)
		}

		pb.RegisterMetricsServiceServer(server, metricsServer)

		log.Printf("Running gRPC server on %s\n", address)

		if err := server.Serve(listen); err != nil {
			log.Fatalf("Could not listen gRPC on %s: %v\n", address, err)
		}
	}()

	return server
}

func main() {
	ctx := context.Background()
	config := NewConfig()

	privateKey, privateKeyErr := utils.LoadPrivateKey(config.CryptoKey)

	if privateKeyErr != nil {
		log.Fatalf("Private key wasn't loaded due to %s", privateKeyErr)
	}

	if err := initializeLogger(config.LogLevel); err != nil {
		log.Fatalf("Logger wasn't initialized due to %s", err)
	}

	storage, healthCheckService, err := initializeStorage(ctx, config)

	if err != nil {
		log.Fatalf("Storage wasn't initialized due to %s", err)
	}

	metricsService := metrics.New(storage)
	server := runServer(ctx, config, metricsService, healthCheckService, privateKey)
	grpcServer := runGRPCServer(metricsService)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-stop

	log.Println("Shutting down the server...")

	grpcServer.GracefulStop()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped gracefully.")
}

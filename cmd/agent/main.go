package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/daremove/go-metrics-service/cmd/buildversion"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/services/stats"
	"github.com/daremove/go-metrics-service/internal/utils"
)

type Job struct {
	metricName  string
	metricValue float64
}

func jobWorker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, config Config, publicKey *rsa.PublicKey) {
	defer wg.Done()

	var (
		ticker  = time.NewTicker(time.Duration(config.ReportInterval) * time.Second)
		payload []models.Metrics
	)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			for d := range jobs {
				payloadItem := models.Metrics{
					ID:    d.metricName,
					MType: models.GaugeMetricType,
				}

				if metrics.IsCounterMetricType(d.metricName) {
					value := int64(d.metricValue)

					payloadItem.MType = models.CounterMetricType
					payloadItem.Delta = &value
				} else {
					value := d.metricValue
					payloadItem.Value = &value
				}

				payload = append(payload, payloadItem)
			}

			if len(payload) > 0 {
				if err := serverrouter.SendMetricModelData(payload, serverrouter.SendMetricModelDataConfig{
					URL:        fmt.Sprintf("http://%s", config.Endpoint),
					SigningKey: config.SigningKey,
					PublicKey:  publicKey,
				}); err != nil {
					log.Printf("failed to send metric data: %s", err)
				}
			}
			return
		case d := <-jobs:
			payloadItem := models.Metrics{
				ID:    d.metricName,
				MType: models.GaugeMetricType,
			}

			if metrics.IsCounterMetricType(d.metricName) {
				value := int64(d.metricValue)

				payloadItem.MType = models.CounterMetricType
				payloadItem.Delta = &value
			} else {
				value := d.metricValue
				payloadItem.Value = &value
			}

			payload = append(payload, payloadItem)
		case <-ticker.C:
			if err := serverrouter.SendMetricModelData(payload, serverrouter.SendMetricModelDataConfig{
				URL:        fmt.Sprintf("http://%s", config.Endpoint),
				SigningKey: config.SigningKey,
				PublicKey:  publicKey,
			}); err != nil {
				log.Printf("failed to send metric data: %s", err)
			} else {
				payload = nil
			}
		}
	}
}

var (
	cpuProvider  = &stats.RealCPUUsageProvider{}
	diskProvider = &stats.RealDiskUsageProvider{}
)

func startReadMetrics(ctx context.Context, wg *sync.WaitGroup, config Config) chan Job {
	jobsCh := make(chan Job, 100)

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(jobsCh)

		var (
			ticker       = time.NewTicker(time.Duration(config.PollInterval) * time.Second)
			statsService = stats.New(cpuProvider, diskProvider)
		)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				utils.Parallelize(
					func() {
						for metricName, metricValue := range statsService.Read() {
							jobsCh <- Job{metricName, metricValue}
						}
					},
					func() {
						data, err := statsService.ReadGopsUtil()

						if err != nil {
							log.Fatal(err)
						}

						for metricName, metricValue := range data {
							jobsCh <- Job{metricName, metricValue}
						}
					},
				)
			}
		}
	}()

	return jobsCh
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	var (
		config = NewConfig()
		wg     sync.WaitGroup
		jobsCh = startReadMetrics(ctx, &wg, config)
	)

	pubicKey, err := utils.LoadPublicKey(config.CryptoKey)

	if err != nil {
		log.Fatalf("Crypto key wasn't loaded due to %s", err)
	}

	log.Printf(
		"Starting read stats data every %v and send it every %v to %s",
		time.Duration(config.PollInterval)*time.Second,
		time.Duration(config.ReportInterval)*time.Second,
		config.Endpoint,
	)

	for i := 0; i < int(config.RateLimit); i++ {
		wg.Add(1)
		go jobWorker(ctx, &wg, jobsCh, config, pubicKey)
	}

	<-stop
	log.Println("Shutting down the agent...")

	cancel()
	wg.Wait()
	log.Println("Agent stopped gracefully.")
}

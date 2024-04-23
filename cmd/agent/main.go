package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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

func jobWorker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, config Config) {
	defer wg.Done()

	var (
		ticker  = time.NewTicker(time.Duration(config.reportInterval) * time.Second)
		payload []models.Metrics
	)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case d := <-jobs:
			payloadItem := models.Metrics{
				ID:    d.metricName,
				MType: "gauge",
			}

			if metrics.IsCounterMetricType(d.metricName) {
				value := int64(d.metricValue)

				payloadItem.MType = "counter"
				payloadItem.Delta = &value
			} else {
				value := d.metricValue
				payloadItem.Value = &value
			}

			payload = append(payload, payloadItem)
		case <-ticker.C:
			if err := serverrouter.SendMetricModelData(payload, serverrouter.SendMetricModelDataConfig{
				URL:        fmt.Sprintf("http://%s", config.endpoint),
				SigningKey: config.signingKey,
			}); err != nil {
				log.Printf("failed to send metric data: %s", err)
			} else {
				payload = nil
			}
		}
	}
}

func startReadMetrics(ctx context.Context, wg *sync.WaitGroup, config Config) chan Job {
	jobsCh := make(chan Job, 100)

	wg.Add(1)

	go func() {
		defer wg.Done()
		defer close(jobsCh)

		var (
			ticker       = time.NewTicker(time.Duration(config.pollInterval) * time.Second)
			statsService = stats.New()
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

	var (
		config = NewConfig()
		wg     sync.WaitGroup
		jobsCh = startReadMetrics(ctx, &wg, config)
	)

	log.Printf(
		"Starting read stats data every %v and send it every %v to %s",
		time.Duration(config.pollInterval)*time.Second,
		time.Duration(config.reportInterval)*time.Second,
		config.endpoint,
	)

	for i := 0; i < int(config.rateLimit); i++ {
		wg.Add(1)
		go jobWorker(ctx, &wg, jobsCh, config)
	}

	wg.Wait()
}

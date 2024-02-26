package main

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/services/stats"
	"github.com/daremove/go-metrics-service/internal/utils"
	"log"
	"time"
)

type Job struct {
	metricName  string
	metricValue float64
}

func requestWorker(jobs <-chan Job, config Config) {
	var (
		ticker  = time.NewTicker(time.Duration(config.reportInterval) * time.Second)
		payload []models.Metrics
	)

	for {
		select {
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

func main() {
	config := NewConfig()

	statsService := stats.New()
	jobsCh := make(chan Job, 100)
	defer close(jobsCh)

	log.Printf(
		"Starting read stats data every %v and send it every %v to %s",
		time.Duration(config.pollInterval)*time.Second,
		time.Duration(config.reportInterval)*time.Second,
		config.endpoint,
	)

	for i := 0; i < int(config.rateLimit); i++ {
		go requestWorker(jobsCh, config)
	}

	utils.Parallelize(
		func() {
			for {
				time.Sleep(time.Duration(config.pollInterval) * time.Second)

				for metricName, metricValue := range statsService.Read() {
					jobsCh <- Job{metricName, metricValue}
				}
			}
		},
		func() {
			for {
				time.Sleep(time.Duration(config.pollInterval) * time.Second)

				data, err := statsService.ReadGopsUtil()

				if err != nil {
					log.Fatal(err)
				}

				for metricName, metricValue := range data {
					jobsCh <- Job{metricName, metricValue}
				}
			}
		},
	)
}

package main

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/models"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/services/stats"
	"github.com/daremove/go-metrics-service/internal/utils"
	"log"
	"sync"
	"time"
)

func main() {
	config := NewConfig()

	var data map[string]float64
	var mutex sync.Mutex
	statsService := stats.New()
	backoff := utils.NewBackoff()

	log.Printf(
		"Starting read stats data every %v and send it every %v to %s",
		time.Duration(config.pollInterval)*time.Second,
		time.Duration(config.reportInterval)*time.Second,
		config.endpoint,
	)

	utils.Parallelize(
		func() {
			for {
				time.Sleep(time.Duration(config.pollInterval) * time.Second)

				mutex.Lock()
				data = statsService.Read()
				mutex.Unlock()
			}
		},
		func() {
			for {
				time.Sleep(time.Duration(config.reportInterval) * time.Second)

				backoff.Reset()
				mutex.Lock()
				var payload []models.Metrics

				for metricName, metricValue := range data {
					payloadItem := models.Metrics{
						ID:    metricName,
						MType: "gauge",
					}

					if metrics.IsCounterMetricType(metricName) {
						value := int64(metricValue)

						payloadItem.MType = "counter"
						payloadItem.Delta = &value
					} else {
						value := metricValue
						payloadItem.Value = &value
					}

					payload = append(payload, payloadItem)
				}

				for {
					if err := serverrouter.SendMetricModelData(fmt.Sprintf("http://%s", config.endpoint), payload); err != nil {
						d, ok := backoff.Duration()

						if !ok {
							log.Printf("failed to send metric data: %s", err)
							break
						} else {
							time.Sleep(d)
							continue
						}
					}
					break
				}

				mutex.Unlock()
			}
		},
	)
}

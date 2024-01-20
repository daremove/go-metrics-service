package main

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/logger"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/services/stats"
	"github.com/daremove/go-metrics-service/internal/utils"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

func main() {
	config := NewConfig()

	var data map[string]float64
	var mutex sync.Mutex
	statsService := stats.New()

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

				mutex.Lock()
				for metricName, metricValue := range data {
					metricType := "gauge"

					if metrics.IsCounterMetricType(metricName) {
						metricType = "counter"
					}

					err := serverrouter.SendMetricData(serverrouter.SendMetricDataParameters{
						URL:         fmt.Sprintf("http://%s", config.endpoint),
						MetricType:  metricType,
						MetricName:  metricName,
						MetricValue: fmt.Sprintf("%v", metricValue),
					})

					if err != nil {
						logger.Log.Error("failed to send metric data", zap.Error(err))
					}
				}
				mutex.Unlock()
			}
		},
	)
}

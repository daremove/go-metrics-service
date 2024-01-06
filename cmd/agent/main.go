package main

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/services/stats"
	"reflect"
	"time"
)

func main() {
	statsService := stats.New()
	iteration := 0

	for {
		time.Sleep(2 * time.Second)
		iteration += 1

		data := statsService.Read()

		if iteration%5 == 0 {
			v := reflect.ValueOf(data)

			for i := 0; i < v.NumField(); i++ {
				metricType := "gauge"
				metricName := v.Type().Field(i).Name
				metricValue := v.Field(i)

				if metrics.IsCounterMetricType(metricName) {
					metricType = "counter"
				}

				err := serverrouter.SendMetricData(serverrouter.SendMetricDataParameters{
					URL:         "http://localhost:8080",
					MetricType:  metricType,
					MetricName:  metricName,
					MetricValue: fmt.Sprintf("%v", metricValue),
				})

				if err != nil {
					panic(err)
				}
			}
		}
	}
}

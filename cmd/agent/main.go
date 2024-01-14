package main

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/http/serverrouter"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/services/stats"
	"github.com/daremove/go-metrics-service/internal/utils"
	"reflect"
	"time"
)

func main() {
	parseFlags()

	var data stats.ReadResult
	statsService := stats.New()

	fmt.Printf("Starting read stats data every %v and send it every %v to %s", time.Duration(pollInterval)*time.Second, time.Duration(reportInterval)*time.Second, endpoint)

	utils.Parallelize(
		func() {
			for {
				time.Sleep(time.Duration(pollInterval) * time.Second)

				data = statsService.Read()
			}
		},
		func() {
			for {
				time.Sleep(time.Duration(reportInterval) * time.Second)

				v := reflect.ValueOf(data)

				for i := 0; i < v.NumField(); i++ {
					metricType := "gauge"
					metricName := v.Type().Field(i).Name
					metricValue := v.Field(i)

					if metrics.IsCounterMetricType(metricName) {
						metricType = "counter"
					}

					err := serverrouter.SendMetricData(serverrouter.SendMetricDataParameters{
						URL:         fmt.Sprintf("http://%s", endpoint),
						MetricType:  metricType,
						MetricName:  metricName,
						MetricValue: fmt.Sprintf("%v", metricValue),
					})

					if err != nil {
						panic(err)
					}
				}
			}
		},
	)
}

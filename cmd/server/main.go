package main

import (
	"github.com/daremove/go-metrics-service/internal/memstorage"
	"github.com/daremove/go-metrics-service/internal/uriparser"
	"net/http"
	"strconv"
)

const (
	updateMetricRoute = "/update/"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	parser := uriparser.New(r.RequestURI, "/method/metricType/metricName/metricValue")

	metricType, ok := parser.GetPathValue("metricType")

	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	metricName, ok := parser.GetPathValue("metricName")

	if !ok {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	metricValue, ok := parser.GetPathValue("metricValue")

	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	switch metricType {
	case Gauge:
		v, err := strconv.ParseFloat(metricValue, 64)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		storage.AddGauge(metricName, v)
	case Counter:
		v, err := strconv.ParseInt(metricValue, 10, 64)

		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		storage.AddCounter(metricName, v)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
}

var storage = memstorage.New()

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(updateMetricRoute, updateMetricHandler)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}
}

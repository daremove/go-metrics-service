package serverrouter

import (
	"fmt"
	"github.com/daremove/go-metrics-service/internal/services/metrics"
	"github.com/daremove/go-metrics-service/internal/utils/uriparser"
	"net/http"
)

type serverRouter struct {
	metricsService metrics.Service
	port           int
}

func New(metricsService metrics.Service, port int) *serverRouter {
	return &serverRouter{metricsService, port}
}

func (router *serverRouter) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updateMetricHandler(router.metricsService))

	err := http.ListenAndServe(fmt.Sprintf(":%v", router.port), mux)

	if err != nil {
		panic(err)
	}
}

// http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
func updateMetricHandler(metricsService metrics.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusNotFound)
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

		if err := metricsService.Save(metrics.SaveParameters{
			MetricType:  metricType,
			MetricName:  metricName,
			MetricValue: metricValue,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusOK)
	}
}

type SendMetricDataParameters struct {
	URL         string
	MetricType  string
	MetricName  string
	MetricValue string
}

func SendMetricData(parameters SendMetricDataParameters) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/update/%s/%s/%s", parameters.URL, parameters.MetricType, parameters.MetricName, parameters.MetricValue), nil)

	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	err = res.Body.Close()

	if err != nil {
		return err
	}

	return nil
}

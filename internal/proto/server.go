// Package proto предназначен для хранения абстракций, связанных с gRPC.
package proto

import (
	"context"

	"github.com/daremove/go-metrics-service/internal/models"
	pb "github.com/daremove/go-metrics-service/internal/proto/metrics"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServiceServer
	metricsService MetricsService
}

type MetricsService interface {
	SaveModels(ctx context.Context, parameters []models.Metrics) error
}

func NewMetricsServer(metricsService MetricsService) *MetricsServer {
	return &MetricsServer{
		metricsService: metricsService,
	}
}

func (metricsServer *MetricsServer) UpdateMetrics(ctx context.Context, in *pb.UpdateMetricsRequest) (*pb.UpdateMetricsResponse, error) {
	var payload = make([]models.Metrics, len(in.Metrics))

	for i, value := range in.Metrics {
		payload[i] = models.Metrics{
			MType: value.Type,
			ID:    value.Id,
			Value: &value.Value,
			Delta: &value.Delta,
		}
	}

	err := metricsServer.metricsService.SaveModels(ctx, payload)

	return nil, err
}

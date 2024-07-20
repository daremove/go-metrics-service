// Package proto предназначен для хранения абстракций, связанных с gRPC.
package proto

import (
	"context"
	"log"

	"github.com/daremove/go-metrics-service/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/daremove/go-metrics-service/internal/proto/metrics"
)

func SendMetricModelData(ctx context.Context, data []models.Metrics) error {
	address := ":3200"
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMetricsServiceClient(conn)
	payload := make([]*pb.Metrics, len(data))

	for i, metric := range data {
		if metric.Value == nil {
			payload[i] = &pb.Metrics{
				Type:  metric.MType,
				Id:    metric.ID,
				Delta: *metric.Delta,
			}
		} else {
			payload[i] = &pb.Metrics{
				Type:  metric.MType,
				Id:    metric.ID,
				Value: *metric.Value,
			}
		}
	}

	_, err = client.UpdateMetrics(ctx, &pb.UpdateMetricsRequest{
		Metrics: payload,
	})

	return err
}

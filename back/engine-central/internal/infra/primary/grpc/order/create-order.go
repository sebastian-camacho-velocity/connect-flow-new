package order

import (
	"context"
	"engine-central/internal/infra/primary/grpc/order/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateOrder implementa el método gRPC para crear órdenes
func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*emptypb.Empty, error) {
	h.logger.Info().
		Interface("CreateOrder request", req).
		Msg("Recibiendo rques Grcp")
	return &emptypb.Empty{}, nil
}

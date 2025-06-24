package order

import (
	"engine-central/internal/app/usecaseorders"
	"engine-central/internal/infra/primary/grpc/order/proto"
	"engine-central/internal/infra/shared/log"

	"google.golang.org/grpc"
)

// OrderGRPCHandler maneja las peticiones gRPC de órdenes
type OrderGRPCHandler struct {
	proto.UnimplementedOrderServiceServer
	usecaseOrder usecaseorders.IUseCaseOrders
	logger       log.ILogger
}

// RegisterGRPCServer crea el handler y lo registra en el servidor gRPC
func NewHandler(usecaseOrder usecaseorders.IUseCaseOrders, logger log.ILogger) *OrderGRPCHandler {
	return &OrderGRPCHandler{
		usecaseOrder: usecaseOrder,
		logger:       logger,
	}
}
func RegisterGRPCServer(grpcServer *grpc.Server, handler *OrderGRPCHandler) {
	proto.RegisterOrderServiceServer(grpcServer, handler)
}

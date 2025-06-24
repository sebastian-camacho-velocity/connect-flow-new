package grpc

import (
	"engine-central/internal/app/usecaseorders"
	"engine-central/internal/infra/primary/grpc/order"
	"engine-central/internal/infra/shared/log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server       *grpc.Server
	listener     net.Listener
	usecaseOrder usecaseorders.IUseCaseOrders
	logger       log.ILogger
}

func New(
	address string,
	usecaseOrder usecaseorders.IUseCaseOrders,
	logger log.ILogger,
) (
	*GRPCServer, error,
) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()

	grpcServer := &GRPCServer{
		server:   server,
		listener: listener,
		logger:   logger,
	}

	// Registrar todos los servicios gRPC
	grpcServer.registerAllServices()

	return grpcServer, nil
}

func (s *GRPCServer) Start() error {
	return s.server.Serve(s.listener)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

func (s *GRPCServer) Register(registerFunc func(*grpc.Server)) {
	registerFunc(s.server)
}

// registerAllServices registra todos los servicios gRPC disponibles
func (s *GRPCServer) registerAllServices() {
	// Registrar servicio de órdenes
	orderHandler := order.NewHandler(s.usecaseOrder, s.logger)
	order.RegisterGRPCServer(s.server, orderHandler)

	// Aquí puedes agregar más servicios conforme los vayas creando:
	// product.RegisterGRPCServer(s.server)
	// customer.RegisterGRPCServer(s.server)
	// shipping.RegisterGRPCServer(s.server)
}

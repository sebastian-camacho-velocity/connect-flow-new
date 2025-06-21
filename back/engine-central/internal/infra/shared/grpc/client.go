package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
)

type GRPCClientConfig struct {
	Address           string
	Timeout           time.Duration
	MaxBackoff        time.Duration
	MinConnectTimeout time.Duration
	WithInsecure      bool // true si no hay TLS
}

// NewGRPCClient crea una conexión gRPC configurable
func NewGRPCClient(cfg GRPCClientConfig) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   cfg.MaxBackoff,
			},
			MinConnectTimeout: cfg.MinConnectTimeout,
		}),
		grpc.WithBlock(), // espera hasta que conecte
	}

	if cfg.WithInsecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		// Puedes agregar aquí configuración TLS si la necesitas
	}

	return grpc.DialContext(ctx, cfg.Address, opts...)
}

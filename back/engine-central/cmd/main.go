package main

import (
	"context"

	"os"
	"os/signal"
	"syscall"

	"engine-central/internal/infra/primary/server"
)

func main() {
	ctx := context.Background()
	services, err := server.InitServer(ctx)
	if err != nil {
		if services != nil && services.Logger != nil {
			services.Logger.Error(ctx).Err(err).Msg("No se pudo inicializar el servidor")
		}
		os.Exit(1)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	services.Logger.Info(ctx).Msg("Apagando servidor...")
}

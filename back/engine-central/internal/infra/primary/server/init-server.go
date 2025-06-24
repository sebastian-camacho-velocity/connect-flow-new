package server

import (
	"context"
	"engine-central/internal/app/usecaseorders"
	"engine-central/internal/domain"
	"engine-central/internal/infra/primary/grpc"
	"engine-central/internal/infra/secundary/httpclient/orderbroker"
	"engine-central/internal/infra/shared/db"
	"engine-central/internal/infra/shared/env"
	"engine-central/internal/infra/shared/log"
	"engine-central/internal/infra/shared/nats"
	"engine-central/internal/infra/shared/s3"
	"fmt"
	"net/http"
)

type AppServices struct {
	Config      env.IConfig
	Logger      log.ILogger
	DB          db.IDatabase
	Nats        nats.INatsClient
	S3          s3.IS3
	OrderBroker domain.IOrderBroker
	GRPCServer  *grpc.GRPCServer
}

func InitServer(ctx context.Context) (*AppServices, error) {

	logger := log.New()

	env, err := env.New(logger)
	if err != nil {
		return nil, err
	}

	database := db.New(logger, env)
	if err := database.Connect(ctx); err != nil {
		return nil, err
	}

	natsClient := nats.New(env, logger)
	if natsClient == nil {
		return nil, err
	}

	s3Client, err := s3.New(env, logger)
	if err != nil {
		return nil, err
	}

	orderBrokerClient := orderbroker.NewClient(env, logger)
	usecaseOrder := usecaseorders.NewOrderUseCase(orderBrokerClient, logger)

	// Inicializar servidor gRPC (con todos los servicios registrados automáticamente)
	grpcPort := env.Get("GRPC_PORT")
	grpcAddr := fmt.Sprintf(":%s", grpcPort)
	grpcServer, err := grpc.New(grpcAddr, usecaseOrder, logger)
	if err != nil {
		return nil, fmt.Errorf("error al crear servidor gRPC: %w", err)
	}

	services := &AppServices{
		Config:      env,
		Logger:      logger,
		DB:          database,
		Nats:        natsClient,
		S3:          s3Client,
		OrderBroker: orderBrokerClient,
		GRPCServer:  grpcServer,
	}

	// Iniciar servidor HTTP
	port := services.Config.Get("API_PORT")
	addr := fmt.Sprintf(":%s", port)
	serverURL := fmt.Sprintf("http://localhost:%s", port)
	coloredURL := fmt.Sprintf("\033[34;4m%s\033[0m", serverURL)

	services.Logger.Info(ctx).Msg("")
	services.Logger.Info(ctx).Msg("")
	services.Logger.Info(ctx).Msgf("    🚀 Servidor HTTP iniciado correctamente")
	services.Logger.Info(ctx).Msgf("    📍 Disponible en: %s", coloredURL)
	services.Logger.Info(ctx).Msg("")

	// Iniciar servidor gRPC
	grpcURL := fmt.Sprintf("grpc://localhost:%s", grpcPort)
	coloredGRPCURL := fmt.Sprintf("\033[32;4m%s\033[0m", grpcURL)
	services.Logger.Info(ctx).Msgf("    🔌 Servidor gRPC iniciado correctamente")
	services.Logger.Info(ctx).Msgf("    📍 Disponible en: %s", coloredGRPCURL)
	services.Logger.Info(ctx).Msg("")

	// Log informativo para NATS
	natsHost := env.Get("NATS_HOST")
	natsPort := env.Get("NATS_PORT")
	natsURL := fmt.Sprintf("nats://%s:%s", natsHost, natsPort)
	coloredNatsURL := fmt.Sprintf("\033[33;4m%s\033[0m", natsURL)

	services.Logger.Info(ctx).Msgf("    📡 Conexión NATS establecida")
	services.Logger.Info(ctx).Msgf("    📍 Conectado a: %s", coloredNatsURL)
	services.Logger.Info(ctx).Msg("")

	// Log informativo para S3
	s3Region := env.Get("S3_REGION")
	s3Bucket := env.Get("S3_BUCKET")
	s3URL := fmt.Sprintf("s3://%s.s3.%s.amazonaws.com", s3Bucket, s3Region)
	coloredS3URL := fmt.Sprintf("\033[35;4m%s\033[0m", s3URL)

	services.Logger.Info(ctx).Msgf("    ☁️  Conexión S3 establecida")
	services.Logger.Info(ctx).Msgf("    📍 Bucket: %s", coloredS3URL)
	services.Logger.Info(ctx).Msg("")

	// Log informativo para MySQL
	dbHost := env.Get("DB_HOST")
	dbPort := env.Get("DB_PORT")
	dbName := env.Get("DB_NAME")
	dbURL := fmt.Sprintf("mysql://%s:%s/%s", dbHost, dbPort, dbName)
	coloredDBURL := fmt.Sprintf("\033[36;4m%s\033[0m", dbURL)

	services.Logger.Info(ctx).Msgf("    🗄️  Conexión MySQL establecida")
	services.Logger.Info(ctx).Msgf("    📍 Base de datos: %s", coloredDBURL)
	services.Logger.Info(ctx).Msg("")

	// Iniciar servidor HTTP en goroutine
	go func() {
		if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
			services.Logger.Error(ctx).Err(err).Msg("Error al iniciar el servidor HTTP")
		}
	}()

	// Iniciar servidor gRPC en goroutine
	go func() {
		if err := grpcServer.Start(); err != nil {
			services.Logger.Error(ctx).Err(err).Msg("Error al iniciar el servidor gRPC")
		}
	}()

	return services, nil
}

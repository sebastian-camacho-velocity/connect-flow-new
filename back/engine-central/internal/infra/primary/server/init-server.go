package server

import (
	"context"
	"engine-central/internal/infra/shared/db"
	"engine-central/internal/infra/shared/env"
	"engine-central/internal/infra/shared/log"
	"engine-central/internal/infra/shared/nats"
	"engine-central/internal/infra/shared/s3"
	"fmt"
	"net/http"
)

type AppServices struct {
	Config env.IConfig
	Logger log.ILogger
	DB     db.IDatabase
	Nats   nats.INatsClient
	S3     s3.IS3
}

func InitServer(ctx context.Context) (*AppServices, error) {

	logger := log.New()
	logger.Info(ctx).Msg("Initializing server...")
	env, err := env.New(logger)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("error loading environment variables")
		return nil, err
	}

	database := db.New(logger, env)
	if err := database.Connect(ctx); err != nil {
		logger.Error(ctx).Err(err).Msg("error connecting to database")
		return nil, err
	}

	natsClient := nats.New(env, logger)
	if natsClient == nil {
		logger.Error(ctx).Msg("error initializing NATS client")
		return nil, err
	}

	s3Client, err := s3.New(env, logger)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("error initializing S3 client")
		return nil, err
	}

	services := &AppServices{
		Config: env,
		Logger: logger,
		DB:     database,
		Nats:   natsClient,
		S3:     s3Client,
	}

	port := services.Config.Get("API_PORT")
	addr := fmt.Sprintf(":%s", port)
	serverURL := fmt.Sprintf("http://localhost:%s", port)

	// Formato con color azul y subrayado para la URL
	coloredURL := fmt.Sprintf("\033[34;4m%s\033[0m", serverURL)

	services.Logger.Info(ctx).Msg("")
	services.Logger.Info(ctx).Msg("")
	services.Logger.Info(ctx).Msgf("    🚀 Servidor HTTP iniciado correctamente")
	services.Logger.Info(ctx).Msgf("    📍 Disponible en: %s", coloredURL)
	services.Logger.Info(ctx).Msg("")
	services.Logger.Info(ctx).Msg("")

	go func() {
		if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
			services.Logger.Error(ctx).Err(err).Msg("Error al iniciar el servidor HTTP")
		}
	}()

	return services, nil
}

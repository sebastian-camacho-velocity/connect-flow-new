package db

import (
	"context"
	"engine-central/internal/infra/shared/env"
	"engine-central/internal/infra/shared/log"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// IDatabase define la interfaz para la conexión a la base de datos
type IDatabase interface {
	Connect(ctx context.Context) error
	Close() error
	Conn(ctx context.Context) *gorm.DB
	WithContext(ctx context.Context) *gorm.DB
}

// database implementa la interfaz IDatabase
type database struct {
	conn     *gorm.DB
	log      log.ILogger
	dbLogger DBLogger
	config   env.IConfig
}

// New crea una nueva instancia de IDatabase con logger y config inyectados
func New(logger log.ILogger, config env.IConfig) IDatabase {
	dbLog := logger.With().Str("component", "database").Logger()
	return &database{
		log:      logger,
		dbLogger: NewDBLogger(dbLog).LogMode(getLogLevel(config)),
		config:   config,
	}
}

// Connect establece la conexión con la base de datos
func (d *database) Connect(ctx context.Context) error {
	tzEscaped := url.QueryEscape(d.config.Get("DB_TIMEZONE"))
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s&interpolateParams=True",
		d.config.Get("DB_USER"),
		d.config.Get("DB_PASS"),
		d.config.Get("DB_HOST"),
		d.config.Get("DB_PORT"),
		d.config.Get("DB_NAME"),
		tzEscaped,
	)

	var err error
	d.conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:              false,
		DisableNestedTransaction: true,
		Logger:                   d.dbLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}
	d.log.Info(ctx).Msg("Database connection established successfully")
	d.conn = d.conn.Omit(clause.Associations).Session(&gorm.Session{
		FullSaveAssociations: false,
	})

	sqlDB, err := d.conn.DB()
	if err != nil {
		return err
	}

	// Valores por defecto si no existen en config
	maxIdle := 25
	maxOpen := 25
	maxLifetime := 5 * time.Minute

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	return nil
}

// Close cierra la conexión con la base de datos
func (d *database) Close() error {
	if d.conn != nil {
		sqlDB, err := d.conn.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetConnection retorna la conexión actual
func (d *database) Conn(ctx context.Context) *gorm.DB {
	return d.conn.WithContext(ctx)
}

// WithContext retorna una nueva instancia de la conexión con el contexto especificado
func (d *database) WithContext(ctx context.Context) *gorm.DB {
	return d.conn.WithContext(ctx)
}

func getLogLevel(config env.IConfig) logger.LogLevel {
	level := config.Get("DB_LOG_LEVEL")
	switch level {
	case "debug":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Error
	}
}

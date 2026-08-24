package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	pgxzero "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/multitracer"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/kushalktamang/blueprint-fullstack/internal/config"
	loggerConfig "github.com/kushalktamang/blueprint-fullstack/internal/logger"
	"github.com/newrelic/go-agent/v3/integrations/nrpgx5"
	"github.com/rs/zerolog"
)

type Database struct {
	Pool *pgxpool.Pool
	log  *zerolog.Logger
}

const DatabasePingTimeout = 10

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerConfig.LoggerService) (*Database, error) {

	dsn := cfg.Database.URL
	if dsn == "" {
		hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

		// URL-encode the password
		encodedPassword := url.QueryEscape(cfg.Database.Password)
		dsn = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
			cfg.Database.User,
			encodedPassword,
			hostPort,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
	}

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	// New Relic implements
	//  QueryTracer, BatchTracer, PrepareTracer, and  ConnectTracer
	// keeps a typed handle so we can register it against
	//  each interface below rather than just QueryTracer.
	var nrTracer *nrpgx5.Tracer
	if loggerService != nil && loggerService.GetApplication() != nil {
		nrTracer = nrpgx5.NewTracer()
	}

	var localTracer *tracelog.TraceLog
	if cfg.Primary.Env == "local" {
		globalLevel := logger.GetLevel()
		pgxLogger := loggerConfig.NewPgxLogger(cfg.Observability, globalLevel)
		localTracer = &tracelog.TraceLog{
			Logger:   pgxzero.NewLogger(pgxLogger),
			LogLevel: tracelog.LogLevel(loggerConfig.GetPgxTraceLogLevel(globalLevel)),
		}
	}

	// Combine whichever tracers are active.
	// multitracer dispatches each event
	// (query/batch/prepare/connect) only to tracers that
	// implement the matching interface,
	// so New Relic's batch/prepare/connect instrumentation
	// isn't silently dropped like it would be with a
	// hand-rolled QueryTracer-only wrapper.
	switch {
	case nrTracer != nil && localTracer != nil:
		pgxPoolConfig.ConnConfig.Tracer = multitracer.New(nrTracer, localTracer)
	case nrTracer != nil:
		pgxPoolConfig.ConnConfig.Tracer = nrTracer
	case localTracer != nil:
		pgxPoolConfig.ConnConfig.Tracer = localTracer
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	database := &Database{
		Pool: pool,
		log:  logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().Msg("connected to the database")
	return database, nil
}

func (db *Database) Close() error {
	db.log.Info().Msg("closing database connection pool")
	db.Pool.Close()
	return nil
}

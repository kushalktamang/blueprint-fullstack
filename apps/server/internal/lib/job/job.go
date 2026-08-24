package job

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/kushalktamang/blueprint-fullstack/internal/config"
	"github.com/kushalktamang/blueprint-fullstack/internal/lib/email"
	"github.com/rs/zerolog"
)

type JobService struct {
	Client      *asynq.Client
	server      *asynq.Server
	logger      *zerolog.Logger
	emailClient *email.Client
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config) (*JobService, error) {
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6, // Higher priority queue for important emails
				"default":  3, // Default priority for most emails
				"low":      1, // Lower priority for non-urgent emails
			},
		},
	)

	emailClient, err := email.NewClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to init email client: %w", err)
	}

	return &JobService{
		Client:      client,
		server:      server,
		logger:      logger,
		emailClient: emailClient,
	}, nil
}

func (j *JobService) Start() error {
	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)

	j.logger.Info().Msg("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")
	j.server.Shutdown()
	j.Client.Close()
}

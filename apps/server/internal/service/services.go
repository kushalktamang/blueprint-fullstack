package service

import (
	"github.com/kushalktamang/blueprint-fullstack/internal/lib/job"
	"github.com/kushalktamang/blueprint-fullstack/internal/repository"
	"github.com/kushalktamang/blueprint-fullstack/internal/server"
)

type Services struct {
	Auth *AuthService
	Job  *job.JobService
}

func NewServices(s *server.Server, repos *repository.Repositories) (*Services, error) {
	authService := NewAuthService(s)

	return &Services{
		Job:  s.Job,
		Auth: authService,
	}, nil
}

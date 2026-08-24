package handler

import (
	"github.com/kushalktamang/blueprint-fullstack/internal/server"
	"github.com/kushalktamang/blueprint-fullstack/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
	}
}

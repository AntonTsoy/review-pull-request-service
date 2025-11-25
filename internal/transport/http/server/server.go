package server

import (
	"context"
	"time"

	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/api"
	"github.com/AntonTsoy/review-pull-request-service/internal/transport/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
}

func New(handlers api.ServerInterface) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	})

	app.Use(middleware.Timeout(3 * time.Second))

	api.RegisterHandlers(app, handlers)

	return &Server{
		app: app,
	}
}

func (s *Server) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

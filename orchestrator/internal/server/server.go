package server

import (
	"CalcYA/orchestrator/config"
	"CalcYA/orchestrator/pkg/db"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	fiber *fiber.App
	cfg   *config.Config
	store *db.ExpressionStore
}

func NewServer(cfg *config.Config, store *db.ExpressionStore) (*Server, error) {
	app := fiber.New(fiber.Config{
		ColorScheme:       fiber.DefaultColors,
		Prefork:           false,
		EnablePrintRoutes: true,
		ServerHeader:      "CalcYA",
		AppName:           cfg.App.Name,
	})
	server := &Server{
		fiber: app,
		cfg:   cfg,
		store: store,
	}
	server.SetupRoutes()
	app.Use(cors.New())
	return server, nil
}

func (s *Server) Run() error {
	address := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	if err := s.fiber.Listen(address); err != nil {
		return err
	}
	return nil
}

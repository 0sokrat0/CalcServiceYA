package server

import (
	"CalcYA/orchestrator/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

type Server struct {
	fiber *fiber.App
	cfg   *config.Config
	db    *gorm.DB
}

func NewServer(cfg *config.Config, db *gorm.DB) (*Server, error) {
	app := fiber.New()
	server := &Server{
		fiber: app,
		cfg:   cfg,
		db:    db,
	}
	server.SetupRoutes()
	app.Use(cors.New())
	return server, nil
}

func (s *Server) SetupRoutes() {
	s.fiber.Get("/api/v1/calculate", s.Calculate)
}

func (s *Server) Run() error {
	address := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	if err := s.fiber.Listen(address); err != nil {
		return err
	}
	return nil
}

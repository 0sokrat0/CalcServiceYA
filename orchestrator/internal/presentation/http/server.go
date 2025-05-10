package http

import (
	"fmt"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/app/expr"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/http/handlers"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcClient "github.com/0sokrat0/GoApiYA/orchestrator/pkg/gen/api/auth"
)

type Server struct {
	fiber    *fiber.App
	cfg      *config.Config
	uc       expr.CalcOrchUsecase
	log      *zap.Logger
	handlers *handlers.Handlers
}

func NewServer(
	cfg *config.Config,
	calcUseCase expr.CalcOrchUsecase,
	log *zap.Logger,
) (*Server, error) {
	app := fiber.New(fiber.Config{
		ColorScheme:       fiber.DefaultColors,
		Prefork:           false,
		EnablePrintRoutes: true,
		ServerHeader:      "CalcYA",
		AppName:           cfg.App.Name,
	})

	srv := &Server{
		fiber: app,
		cfg:   cfg,
		uc:    calcUseCase,
		log:   log,
	}

	authAddr := fmt.Sprintf("%s:%s", cfg.Auth.Host, cfg.Auth.Port)
	conn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial auth service", zap.Error(err))
	}

	authClient := grpcClient.NewAuthClient(conn)

	srv.handlers = handlers.NewHandlers(app, cfg, calcUseCase, log, authClient)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
	}))

	jwtMiddleware := middleware.JWTProtected(cfg.JWT.JWTSecret)
	srv.SetupRoutes(jwtMiddleware)

	return srv, nil
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
	return s.fiber.Listen(addr)
}

func (s *Server) App() *fiber.App {
	return s.fiber
}

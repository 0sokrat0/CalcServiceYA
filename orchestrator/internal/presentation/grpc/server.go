package grpc

import (
	"context"
	"net"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/app/expr"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/grpc/handlers"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/0sokrat0/GoApiYA/orchestrator/pkg/gen/api/task"
)

type Server struct {
	server *grpc.Server
	log    *zap.Logger
	cfg    *config.Config
}

func NewServer(log *zap.Logger, cfg *config.Config, uc expr.CalcOrchUsecase) *Server {
	interceptorChain := grpc_middleware.ChainUnaryServer(
		grpc_recovery.UnaryServerInterceptor(),
		loggingInterceptor(log),
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptorChain),
	)

	taskHandler := handlers.NewTaskHandler(uc, log)

	pb.RegisterTaskServiceServer(grpcServer, taskHandler)

	return &Server{
		server: grpcServer,
		log:    log,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.cfg.Grpc.Port)
	if err != nil {
		s.log.Error("failed to listen", zap.Error(err))
		return err
	}

	s.log.Info("gRPC server starting", zap.String("port", s.cfg.Grpc.Port))
	return s.server.Serve(lis)
}

func (s *Server) GracefulShutdown() {
	s.log.Info("shutting down gRPC server")
	s.server.GracefulStop()
}

func loggingInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Info("gRPC request",
			zap.String("method", info.FullMethod),
			zap.Any("request", req),
		)

		resp, err := handler(ctx, req)
		if err != nil {
			log.Error("gRPC error",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
		}
		return resp, err
	}
}

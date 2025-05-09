package main

import (
	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/app/expr"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/infrastructure/persistence"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/grpc"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/presentation/http"
	sqlite "github.com/0sokrat0/GoApiYA/orchestrator/pkg/db/SQLite"
	"github.com/0sokrat0/GoApiYA/orchestrator/pkg/logger"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	cfg := config.GetConfig()
	log := logger.InitLogger(cfg)
	conn := sqlite.SQLiteConnect()

	expStore := persistence.NewExpressionRepoGORM(conn)
	taskStore := persistence.NewTaskRepoGORM(conn)

	service := expr.NewCalcOrch(taskStore, expStore, cfg, log)

	server, err := http.NewServer(cfg, service, log)
	if err != nil {
		log.Fatal("Fatal to create server ❌", zap.Error(err))
	}
	grpcServer := grpc.NewServer(log, cfg, service)
	go grpcServer.Start()

	server.Run()
}

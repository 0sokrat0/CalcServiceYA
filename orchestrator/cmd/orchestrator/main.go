package main

import (
	"CalcYA/orchestrator/config"
	"CalcYA/orchestrator/internal/server"
	"CalcYA/orchestrator/pkg/db"
	"log"

	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {

	cfg, err := config.LoadConfig("./orchestrator/config")
	if err != nil {
		log.Fatal("Fatal to load config", err)
	}

	zap.L().Info("Loaded config",
		zap.String("AppName", cfg.App.Name),
		zap.String("ServerHost", cfg.Server.Host),
		zap.String("ServerPort", cfg.Server.Port),
		zap.String("SQLitePath", cfg.SQLite.Path),
	)

	db.InitDB(cfg)

	server, err := server.NewServer(cfg, nil)
	if err != nil {
		zap.L().Fatal("Fatal to create server", zap.Error(err))
	}

	server.Run()

}

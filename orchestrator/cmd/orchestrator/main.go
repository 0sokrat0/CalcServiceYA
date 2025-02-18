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
		log.Fatal("Fatal to load config ❌", err)
	}

	store := db.NewExpressionStore()

	server, err := server.NewServer(cfg, store)
	if err != nil {
		zap.L().Fatal("Fatal to create server ❌", zap.Error(err))
	}

	server.Run()

}

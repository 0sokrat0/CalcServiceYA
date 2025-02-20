package main

import (
	"agent/config"
	"agent/internal/server"

	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	cfg, err := config.LoadConfig(".agent/config")
	if err != nil {
		zap.L().Panic("Failed to load config")
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		zap.L().Fatal("Fatal to create server ❌", zap.Error(err))
	}

	server.Run()

}

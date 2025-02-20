package main

import (
	"log"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {

	cfg, err := config.LoadConfig("../config")
	if err != nil {
		log.Fatal("Fatal to load config ❌", err)
	}

	stores := db.NewStores()
	zap.L().Info("In-memory хранилеще запущено")

	server, err := server.NewServer(cfg, stores)
	if err != nil {
		zap.L().Fatal("Fatal to create server ❌", zap.Error(err))
	}

	server.Run()
}

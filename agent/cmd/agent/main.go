package main

import (
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/0sokrat0/GoApiYA/agent/config"
	"github.com/0sokrat0/GoApiYA/agent/internal/deamon"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		slog.Error("Fatal to load config", "error", err)
		panic(err)
	}

	slog.Info("Starting deamon")
	tasksChan := make(chan deamon.Task, 50)

	go deamon.CalcPool(tasksChan, cfg)

	go deamon.GetTask(*cfg, tasksChan)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigs
	slog.Info("Received signal, shutting down", "signal", s)
}

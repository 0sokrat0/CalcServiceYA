package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/0sokrat0/GoApiYA/agent/config"
	"github.com/0sokrat0/GoApiYA/agent/internal/deamon"
	"github.com/0sokrat0/GoApiYA/agent/internal/grpc"
	"github.com/0sokrat0/GoApiYA/agent/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		fmt.Printf("Fatal to load config: %v\n", err)
		panic(err)
	}

	log := logger.InitLogger(cfg)
	log.Info("Starting deamon")
	tasksChan := make(chan deamon.Task, 50)

	conn, err := grpc.NewCalcClient(cfg)
	if err != nil {
		log.Fatal("err connect gRPC", zap.Error(err))
	}

	dm := deamon.NewDemon(cfg, conn, log)

	go dm.CalcPool(tasksChan)

	go dm.GetTask(tasksChan)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigs
	log.Info("Received signal, shutting down", zap.Any("signal", s))
}

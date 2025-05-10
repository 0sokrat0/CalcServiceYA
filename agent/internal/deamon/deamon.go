package deamon

import (
	"context"
	"fmt"
	"time"

	"github.com/0sokrat0/GoApiYA/agent/config"
	gen "github.com/0sokrat0/GoApiYA/agent/pkg/gen/api/task"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
)

type Demon struct {
	cfg    *config.Config
	client gen.TaskServiceClient
	log    *zap.Logger
}

func NewDemon(cfg *config.Config,
	client gen.TaskServiceClient,
	log *zap.Logger) *Demon {
	return &Demon{
		cfg:    cfg,
		client: client,
		log:    log,
	}
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
}

func (d *Demon) CalcPool(tasks <-chan Task) {
	workers := d.cfg.App.COMPUTING_POWER
	d.log.Info("Starting worker pool", zap.Int("workers", workers))
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			for task := range tasks {
				d.log.Info("Worker starts task",
					zap.Int("workerID", workerID),
					zap.String("taskID", task.ID),
				)
				result, err := compute(task, d.log)
				if err != nil {
					d.log.Error("Compute failed", zap.String("taskID", task.ID), zap.Error(err))
					continue
				}
				d.ResultTask(task.ID, result)
				d.log.Info("Worker finished task",
					zap.Int("workerID", workerID),
					zap.String("taskID", task.ID),
					zap.Float64("result", result),
				)
			}
		}(i)
	}
}

func compute(task Task, log *zap.Logger) (float64, error) {
	log.Info("Compute op",
		zap.String("taskID", task.ID),
		zap.String("operation", task.Operation),
		zap.Float64("arg1", task.Arg1),
		zap.Float64("arg2", task.Arg2),
	)
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = task.Arg1 / task.Arg2
	default:
		return 0, fmt.Errorf("unsupported operation %q", task.Operation)
	}
	return result, nil
}

func (d *Demon) ResultTask(id string, result float64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req := &gen.UpdateTaskRequest{
		Id:     id,
		Result: result,
	}
	resp, err := d.client.UpdateTask(ctx, req)
	if err != nil {
		d.log.Error("gRPC UpdateTask failed", zap.String("taskID", id), zap.Error(err))
		return
	}

	d.log.Info("Task updated on server",
		zap.String("taskID", resp.ID),
		zap.Float64("result", result),
	)
}

func (d *Demon) GetTask(tasksChan chan<- Task) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		resp, err := d.client.GetNextTask(ctx, &empty.Empty{})
		if err != nil {
			d.log.Error("gRPC GetNextTask failed", zap.Error(err))
			time.Sleep(10 * time.Second)
			continue
		}
		cancel()

		if !resp.Found || resp.Task.ID == "" {
			d.log.Debug("No task available, retrying")
			time.Sleep(time.Second)
			continue
		}

		t := Task{
			ID:            resp.Task.ID,
			Arg1:          float64(resp.Task.Arg1),
			Arg2:          float64(resp.Task.Arg2),
			Operation:     resp.Task.Operation,
			OperationTime: resp.Task.OperationTime,
		}
		d.log.Info("Received task", zap.String("taskID", t.ID), zap.String("operation", t.Operation))
		tasksChan <- t

		time.Sleep(time.Second)
	}
}

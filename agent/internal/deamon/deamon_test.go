package deamon

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/0sokrat0/GoApiYA/agent/config"
	gen "github.com/0sokrat0/GoApiYA/agent/pkg/gen/api/task"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type mockClient struct {
	updateFunc  func(context.Context, *gen.UpdateTaskRequest, ...grpc.CallOption) (*gen.Task, error)
	getNextFunc func(context.Context, *emptypb.Empty, ...grpc.CallOption) (*gen.GetNextTaskResponse, error)
}

func (m *mockClient) UpdateTask(ctx context.Context, req *gen.UpdateTaskRequest, opts ...grpc.CallOption) (*gen.Task, error) {
	return m.updateFunc(ctx, req, opts...)
}
func (m *mockClient) GetNextTask(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*gen.GetNextTaskResponse, error) {
	return m.getNextFunc(ctx, in, opts...)
}

func TestComputeSubtract(t *testing.T) {
	log, _ := zap.NewDevelopment()
	res, err := compute(Task{"sub", 10, 4, "-", 0}, log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != 6 {
		t.Errorf("expected 6, got %v", res)
	}
}

func TestComputeMultiply(t *testing.T) {
	log, _ := zap.NewDevelopment()
	res, err := compute(Task{"mul", 3, 5, "*", 0}, log)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != 15 {
		t.Errorf("expected 15, got %v", res)
	}
}

func TestComputeDivideZero(t *testing.T) {
	log, _ := zap.NewDevelopment()
	_, err := compute(Task{"div0", 1, 0, "/", 0}, log)
	exp := "division by zero"
	if err == nil || !strings.Contains(err.Error(), exp) {
		t.Errorf("expected error containing %q, got %v", exp, err)
	}
}

func TestResultTask_NoClientCallOnError(t *testing.T) {
	cfg := &config.Config{App: config.AppConfig{COMPUTING_POWER: 1}}
	called := false
	client := &mockClient{
		updateFunc: func(ctx context.Context, req *gen.UpdateTaskRequest, _ ...grpc.CallOption) (*gen.Task, error) {
			called = true
			return nil, errors.New("rpc fail")
		},
	}
	d := NewDemon(cfg, client, zap.NewNop())
	d.ResultTask("id", 123)
	if !called {
		t.Error("expected UpdateTask to be called even on error")
	}
}

func TestGetTask_ErrorRetry(t *testing.T) {
	cfg := &config.Config{App: config.AppConfig{COMPUTING_POWER: 1}}
	count := 0

	client := &mockClient{
		getNextFunc: func(ctx context.Context, in *emptypb.Empty, _ ...grpc.CallOption) (*gen.GetNextTaskResponse, error) {
			count++
			if count == 1 {
				return nil, errors.New("network")
			}
			return &gen.GetNextTaskResponse{Found: false}, nil
		},
	}

	tasks := make(chan Task, 1)
	d := NewDemon(cfg, client, zap.NewNop())
	go d.GetTask(tasks)

	select {
	case <-time.After(11 * time.Second):
		if count < 2 {
			t.Errorf("expected at least 2 attempts, got %d", count)
		}
	}
}

func TestCalcPool_DrainChannel(t *testing.T) {
	cfg := &config.Config{App: config.AppConfig{COMPUTING_POWER: 2}}

	client := &mockClient{
		updateFunc: func(ctx context.Context, req *gen.UpdateTaskRequest, _ ...grpc.CallOption) (*gen.Task, error) {
			return &gen.Task{ID: req.Id}, nil
		},
	}

	d := NewDemon(cfg, client, zap.NewNop())
	tasks := make(chan Task)

	d.CalcPool(tasks)

	tasks <- Task{"p1", 1, 1, "+", 0}
	close(tasks)

	time.Sleep(100 * time.Millisecond)
}

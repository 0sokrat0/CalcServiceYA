package handlers

import (
	"context"

	"github.com/0sokrat0/GoApiYA/orchestrator/internal/app/expr"
	pb "github.com/0sokrat0/GoApiYA/orchestrator/pkg/gen/api/task"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskHandler struct {
	pb.UnimplementedTaskServiceServer
	uc  expr.CalcOrchUsecase
	log *zap.Logger
}

func NewTaskHandler(uc expr.CalcOrchUsecase, log *zap.Logger) *TaskHandler {
	return &TaskHandler{uc: uc, log: log}
}

func (h *TaskHandler) GetNextTask(ctx context.Context, _ *emptypb.Empty) (*pb.GetNextTaskResponse, error) {
	task, found, err := h.uc.GetNextTask()
	if err != nil {
		h.log.Error("GetNextTask failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	if !found {
		return &pb.GetNextTaskResponse{Found: false}, nil
	}
	return &pb.GetNextTaskResponse{
		Found: true,
		Task: &pb.Task{
			ID:            task.ID,
			Arg1:          float32(task.Arg1),
			Arg2:          float32(task.Arg2),
			Operation:     task.Operation,
			OperationTime: task.OperationTime,
		},
	}, nil
}

func (h *TaskHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.Task, error) {
	updatedTask, err := h.uc.UpdateTaskResult(req.Id, req.Result)
	if err != nil {
		h.log.Error("UpdateTask failed", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "internal error")
	}
	if err := h.uc.ProcessTaskResult(req.Id, req.Result); err != nil {
		h.log.Error("ProcessTaskResult failed", zap.Error(err))
	}
	return &pb.Task{
		ID:            updatedTask.ID,
		Arg1:          float32(updatedTask.Arg1),
		Arg2:          float32(updatedTask.Arg2),
		Operation:     updatedTask.Operation,
		OperationTime: updatedTask.OperationTime,
	}, nil
}

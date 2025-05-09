package expr

import (
	"errors"
	"fmt"

	"github.com/0sokrat0/GoApiYA/orchestrator/config"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CalcOrch struct {
	taskRepo repository.TaskRepository
	expRepo  repository.ExpressionRepository
	cfg      *config.Config
	log      *zap.Logger
}

type CalcOrchUsecase interface {
	CreateExp(expression string, ownerID string) (*entity.Expression, error)
	GetExpression(id string) (*entity.Expression, error)
	ListExpressions(ownerID string) ([]entity.Expression, error)
	ListTasks() ([]entity.Task, error)
	ProcessTaskResult(taskID string, result float64) error
	GetNextTask() (entity.Task, bool, error)
	UpdateTaskResult(taskID string, result float64) (*entity.Task, error)
}

func NewCalcOrch(tR repository.TaskRepository, eR repository.ExpressionRepository, cfg *config.Config, log *zap.Logger) CalcOrchUsecase {
	return &CalcOrch{
		taskRepo: tR,
		expRepo:  eR,
		cfg:      cfg,
		log:      log,
	}
}

func (cu *CalcOrch) CreateExp(expression string, ownerID string) (*entity.Expression, error) {
	exprRecord := entity.Expression{
		ID:         uuid.NewString(),
		OwnerID:    ownerID,
		RootTaskID: "",
		Status:     entity.StatusPending,
		Result:     0,
	}
	cu.expRepo.Create(exprRecord)

	tokens := tokenize(expression)
	parser := NewParser(tokens)
	ast := parser.parseExpression(0)

	rootID, tasks := GenerateTasks(ast, cu.cfg)
	cu.log.Info("Корневой таск ID: %s\n", zap.String("root ID:", rootID))

	exprRecord.RootTaskID = rootID
	cu.expRepo.Update(exprRecord)

	for _, task := range tasks {
		if err := cu.taskRepo.Create(task); err != nil {
			return &exprRecord, err
		}
	}
	return &exprRecord, nil
}

func (cu *CalcOrch) GetExpression(id string) (*entity.Expression, error) {
	expr, ok := cu.expRepo.GetByID(id)
	if !ok {
		return nil, errors.New("expression not found")
	}
	if expr.RootTaskID != "" {
		task, found, err := cu.taskRepo.GetByID(expr.RootTaskID)
		if err != nil {
			return nil, err
		}
		if found && task.Status == entity.StatusSuccess {
			expr.Result = task.Result
			expr.Status = entity.StatusSuccess
			if err := cu.expRepo.Update(*expr); err != nil {
				return nil, err
			}
		}
	}

	return expr, nil
}

func (uc *CalcOrch) ListExpressions(OwnerID string) ([]entity.Expression, error) {
	list, err := uc.expRepo.List(OwnerID)
	if err != nil {
		return nil, err
	}

	for i := range list {
		expr := &list[i]
		if expr.RootTaskID == "" {
			continue
		}

		task, found, err := uc.taskRepo.GetByID(expr.RootTaskID)
		if err != nil {
			return nil, err
		}

		if found && task.Status == entity.StatusSuccess {
			expr.Result = task.Result
			expr.Status = entity.StatusSuccess

			if err := uc.expRepo.Update(*expr); err != nil {
				return nil, fmt.Errorf("failed to update expression %s: %w", expr.ID, err)
			}
		}
	}

	return list, nil
}

func (uc *CalcOrch) ListTasks() ([]entity.Task, error) {
	tasks, err := uc.taskRepo.List()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (cu *CalcOrch) ProcessTaskResult(taskID string, result float64) error {
	parents, err := cu.taskRepo.List()
	if err != nil {
		return err
	}
	for _, p := range parents {
		updated := false
		if p.LeftTaskID == taskID {
			p.Arg1 = result
			p.LeftTaskID = ""
			updated = true
		}
		if p.RightTaskID == taskID {
			p.Arg2 = result
			p.RightTaskID = ""
			updated = true
		}
		if !updated {
			continue
		}
		if p.LeftTaskID == "" && p.RightTaskID == "" {
			p.Status = entity.StatusReady
		}
		if err := cu.taskRepo.Replace(p); err != nil {
			cu.log.Warn("taskRepo.Replace failed", zap.Error(err), zap.String("parentID", p.ID))
		}
	}
	return nil
}

func (cu *CalcOrch) GetNextTask() (entity.Task, bool, error) {
	return cu.taskRepo.GetNextTask()
}

func (cu *CalcOrch) UpdateTaskResult(taskID string, result float64) (*entity.Task, error) {
	updated, err := cu.taskRepo.Update(taskID, result)
	if err != nil {
		cu.log.Error("taskRepo.Update failed", zap.Error(err), zap.String("taskID", taskID))
		return &entity.Task{}, err
	}
	return updated, nil
}

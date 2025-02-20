package db

import (
	"fmt"
	"sync"
)

type Stores struct {
	Expression *ExpressionStore
	Task       *TaskStore
}

func NewStores() *Stores {
	return &Stores{
		Expression: NewExpressionStore(),
		Task:       NewTaskStore(),
	}
}

type ExpressionRepository interface {
	GetByID(id string) (*Expression, bool)
	List() []Expression
	Delete(id string)
	Create(expr Expression)
	Update(expr Expression)
}

type Status string

const (
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
	StatusError   Status = "error"
	StatusPending Status = "pending"
	StatusWaiting Status = "waiting"
	StatusReady   Status = "ready"
	StatusRunning Status = "running"
)

type Expression struct {
	ID         string  `json:"id"`
	RootTaskID string  `json:"rootTaskID"`
	Status     Status  `json:"status"`
	Result     float64 `json:"result"`
}

type ExpressionStore struct {
	mu          sync.RWMutex
	expressions map[string]Expression
}

func NewExpressionStore() *ExpressionStore {
	return &ExpressionStore{
		expressions: make(map[string]Expression),
	}
}

func (s *ExpressionStore) GetByID(id string) (*Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expr, ok := s.expressions[id]
	return &expr, ok
}

func (s *ExpressionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.expressions, id)
}

func (s *ExpressionStore) List() ([]Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expressions := make([]Expression, 0, len(s.expressions))
	for _, expr := range s.expressions {
		expressions = append(expressions, expr)
	}
	return expressions, nil
}

func (s *ExpressionStore) Create(expr Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expressions[expr.ID] = expr
}

func (s *ExpressionStore) Update(expr Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.expressions[expr.ID] = expr
}

func (s *ExpressionStore) UpdateResult(exprID string, result float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	expr, ok := s.expressions[exprID]
	if !ok {
		return fmt.Errorf("expression with id %s not found", exprID)
	}

	expr.Result = result
	expr.Status = StatusSuccess
	s.expressions[exprID] = expr
	return nil
}

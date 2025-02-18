package db

import "sync"

type ExpressionRepository interface {
	GetByID(id string) (*Expression, bool)
	List() []Expression
	Delete(id string)
	Create(expr Expression)
	Update(expr Expression)
}

type Status string

const (
	StatusSuccess  Status = "success"
	StatusFailed   Status = "failed"
	StatusError    Status = "error"
	StatusAccepted Status = "accepted"
)

type Expression struct {
	ID     string  `json:"id"`
	Status Status  `json:"status"`
	Result float64 `json:"result"`
}

type ExpressionStore struct {
	mu         sync.RWMutex
	expression map[string]Expression
}

func NewExpressionStore() *ExpressionStore {
	return &ExpressionStore{
		expression: make(map[string]Expression),
	}
}

func (s *ExpressionStore) GetByID(id string) (*Expression, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expr, ok := s.expression[id]
	return &expr, ok
}

func (s *ExpressionStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.expression, id)
}

func (s *ExpressionStore) List() ([]Expression, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expressions := make([]Expression, 0, len(s.expression))
	for _, expr := range s.expression {
		expressions = append(expressions, expr)
	}
	return expressions, nil
}

func (s *ExpressionStore) Create(expr Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.expression[expr.ID] = expr
}

func (s *ExpressionStore) Update(expr Expression) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.expression[expr.ID] = expr
}

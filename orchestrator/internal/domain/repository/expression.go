package repository

import (
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
)

type ExpressionRepository interface {
	GetByID(id string) (*entity.Expression, bool)
	List(OwnerID string) ([]entity.Expression, error)
	Create(expr entity.Expression) error
	Update(expr entity.Expression) error
	Delete(id string) error
	UpdateResult(exprID string, result float64) error
}

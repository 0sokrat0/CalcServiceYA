package repository

import (
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
)

type TaskRepository interface {
	List() ([]entity.Task, error)
	Delete(id string) error
	Create(task entity.Task) error
	ListResponses() ([]entity.Task, error)
	Update(id string, result float64) (*entity.Task, error)
	Replace(task entity.Task) error
	GetNextTask() (entity.Task, bool, error)
	GetByID(id string) (entity.Task, bool, error)
}

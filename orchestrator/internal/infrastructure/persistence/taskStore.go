package persistence

import (
	"errors"
	"fmt"

	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/repository"
	"github.com/0sokrat0/GoApiYA/orchestrator/migrations/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// type TasksResponse struct {
// 	ID            string  `json:"id"`
// 	Arg1          float64 `json:"arg1"`
// 	Arg2          float64 `json:"arg2"`
// 	Operation     string  `json:"operation"`
// 	OperationTime int64   `json:"operation_time"`
// }

var ErrNotFound = errors.New("not found")

type taskRepoGORM struct {
	db *gorm.DB
}

func NewTaskRepoGORM(db *gorm.DB) repository.TaskRepository {
	db.AutoMigrate(&models.Task{})
	return &taskRepoGORM{db: db}
}

func (r *taskRepoGORM) List() ([]entity.Task, error) {
	var tasks []entity.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepoGORM) Delete(id string) error {
	return r.db.Delete(&entity.Task{}, "id = ?", id).Error
}

func (r *taskRepoGORM) Create(task entity.Task) error {
	return r.db.Create(&task).Error
}

func (r *taskRepoGORM) ListResponses() ([]entity.Task, error) {
	tasks, err := r.List()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepoGORM) Update(id string, result float64) (*entity.Task, error) {
	var task entity.Task
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	task.Result = result
	task.Status = entity.StatusSuccess
	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepoGORM) Replace(task entity.Task) error {
	var existing entity.Task
	if err := r.db.First(&existing, "id = ?", task.ID).Error; err != nil {
		return fmt.Errorf("task not found")
	}
	return r.db.Save(&task).Error
}

func (r *taskRepoGORM) GetNextTask() (entity.Task, bool, error) {
	tx := r.db.Begin()
	defer tx.Rollback()

	var task entity.Task
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("status = ?", entity.StatusReady).
		First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Task{}, false, nil
		}
		return entity.Task{}, false, err
	}

	if err := tx.Model(&task).Update("status", entity.StatusRunning).Error; err != nil {
		return entity.Task{}, false, err
	}
	tx.Commit()

	return task, true, nil
}

func (r *taskRepoGORM) GetByID(id string) (entity.Task, bool, error) {
	var task entity.Task
	if err := r.db.First(&task, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Task{}, false, nil
		}
		return entity.Task{}, false, err
	}
	return task, true, nil
}

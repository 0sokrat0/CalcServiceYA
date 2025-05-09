package persistence

import (
	"fmt"

	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/entity"
	"github.com/0sokrat0/GoApiYA/orchestrator/internal/domain/repository"
	"github.com/0sokrat0/GoApiYA/orchestrator/migrations/models"
	"gorm.io/gorm"
)

type expressionRepoGORM struct {
	db *gorm.DB
}

func NewExpressionRepoGORM(db *gorm.DB) repository.ExpressionRepository {
	return &expressionRepoGORM{db: db}
}

func (r *expressionRepoGORM) GetByID(id string) (*entity.Expression, bool) {
	var expr entity.Expression
	if err := r.db.First(&expr, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(err)
	}
	return &expr, true
}

func (r *expressionRepoGORM) List(ownerID string) ([]entity.Expression, error) {
	var list []entity.Expression
	if err := r.db.Where("owner_id = ?", ownerID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *expressionRepoGORM) Create(expr entity.Expression) error {
	return r.db.Create(&expr).Error
}

func (r *expressionRepoGORM) Update(expr entity.Expression) error {
	return r.db.Save(&expr).Error
}

func (r *expressionRepoGORM) Delete(id string) error {
	return r.db.Delete(&entity.Expression{}, "id = ?", id).Error
}

func (r *expressionRepoGORM) UpdateResult(exprID string, result float64) error {
	res := r.db.Model(&entity.Expression{}).
		Where("id = ?", exprID).
		Updates(map[string]interface{}{
			"result": result,
			"status": models.StatusSuccess,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("expression with id %s not found", exprID)
	}
	return nil
}

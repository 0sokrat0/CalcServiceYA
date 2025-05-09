package sqlite

import (
	"github.com/0sokrat0/GoApiYA/orchestrator/migrations/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SQLiteConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("calc.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Expression{}, &models.Task{})
	return db
}

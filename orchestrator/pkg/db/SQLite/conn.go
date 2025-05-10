package sqlite

import (
	"github.com/0sokrat0/GoApiYA/orchestrator/migrations/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SQLiteConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("calc.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Expression{}, &models.Task{})
	return db
}

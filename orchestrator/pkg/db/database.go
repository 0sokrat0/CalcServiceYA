package db

import (
	"CalcYA/orchestrator/config"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.SQLite.Path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate()

	zap.L().Info("Connected to database")

	return db, nil

}

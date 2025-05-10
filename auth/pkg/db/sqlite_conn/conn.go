package sqliteconn

import (
	"auth/config"
	"auth/pkg/db/sqlite_conn/models"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLiteDB(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database.SQLitePath), &gorm.Config{})
	if err != nil {
		log.Error("failed to connect sqlite", zap.Error(err))
		return nil, err
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Error("failed to migrate sqlite", zap.Error(err))
		return nil, err
	}
	log.Info("sqlite connection established")
	return db, nil
}

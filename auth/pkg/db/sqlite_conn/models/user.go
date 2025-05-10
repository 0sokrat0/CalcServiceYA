package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"size:255;not null;uniqueIndex"`
	PasswordHash string    `gorm:"type:text;not null"`
	Role         string    `gorm:"size:20;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

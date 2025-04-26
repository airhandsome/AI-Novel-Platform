package models

import (
	"time"
)

type Favorite struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null"`
	NovelID   uint      `gorm:"not null"`
	CreatedAt time.Time
} 
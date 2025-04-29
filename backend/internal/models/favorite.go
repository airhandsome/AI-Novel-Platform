package models

import (
	"time"
)

// Favorite 表示用户收藏的小说
type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	NovelID   uint      `gorm:"index" json:"novelId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
} 
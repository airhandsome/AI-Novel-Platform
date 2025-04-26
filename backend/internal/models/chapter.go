package models

import (
	"time"
)

type Chapter struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	NovelID   uint      `json:"novelId" gorm:"not null"`
	Title     string    `json:"title" gorm:"size:100;not null"`
	Content   string    `json:"content" gorm:"type:text"`
	WordCount int       `json:"wordCount"`
	Order     int       `json:"order" gorm:"not null"` // 章节顺序
	Status    int       `json:"status" gorm:"default:0"`
	Novel     Novel     `json:"-" gorm:"foreignKey:NovelID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateTime"`
}

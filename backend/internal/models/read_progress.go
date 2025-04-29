package models

import (
	"time"
)

// ReadProgress 用户阅读进度记录
type ReadProgress struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"not null"`
	NovelID   uint      `json:"novelId" gorm:"not null"`
	ChapterID uint      `json:"chapterId" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// 外键关联
	User    User    `json:"-" gorm:"foreignKey:UserID"`
	Novel   Novel   `json:"-" gorm:"foreignKey:NovelID"`
	Chapter Chapter `json:"-" gorm:"foreignKey:ChapterID"`
}

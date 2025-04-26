package models

import (
	"time"
	
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primarykey"`
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	Avatar       string `gorm:"default:''"`
	Bio          string `gorm:"type:text"`
	ResetToken   string `gorm:"default:''"`
	ResetExpires time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// BeforeCreate 在创建记录前设置默认值
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Avatar == "" {
		u.Avatar = ""
	}
	if u.Bio == "" {
		u.Bio = ""
	}
	return nil
}

// BeforeUpdate 在更新记录前设置默认值
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return nil
}
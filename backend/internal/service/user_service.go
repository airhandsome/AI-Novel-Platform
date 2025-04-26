package service

import (
	"ai-novel-platform/internal/models"
	"ai-novel-platform/internal/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(username, password, email string) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建新用户
	user := models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		Bio:          "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ResetExpires: time.Now().Add(1 * time.Hour),
	}

	return s.db.Create(&user).Error
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func (s *UserService) RequestPasswordReset(email string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	// 生成重置token
	resetToken := utils.GenerateRandomString(32)
	resetExpires := time.Now().Add(1 * time.Hour)

	// 更新用户的重置token
	user.ResetToken = resetToken
	user.ResetExpires = resetExpires
	if err := s.db.Save(&user).Error; err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *UserService) ResetPassword(token, newPassword string) error {
	var user models.User
	if err := s.db.Where("reset_token = ? AND reset_expires > ?", token, time.Now()).First(&user).Error; err != nil {
		return errors.New("invalid or expired reset token")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码并清除重置token
	user.PasswordHash = string(hashedPassword)
	user.ResetToken = ""
	user.ResetExpires = time.Time{}
	return s.db.Save(&user).Error
}

func (s *UserService) UpdateUser(userID uint, updates map[string]interface{}) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (s *UserService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.db.Save(&user).Error
}

func (s *UserService) UpdateAvatar(userID uint, avatarURL string) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL).Error
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserStats 获取用户统计信息 - 简化版
func (s *UserService) GetUserStats(userID uint) (map[string]int64, error) {
	stats := map[string]int64{
		"novelCount":    0,
		"wordCount":     0,
		"favoriteCount": 0,
	}

	// 获取小说数量 - 忽略错误
	var novelCount int64
	s.db.Model(&models.Novel{}).Where("author_id = ?", userID).Count(&novelCount)
	stats["novelCount"] = novelCount

	// 获取总字数 - 忽略错误
	var wordCount int64
	s.db.Model(&models.Novel{}).
		Where("author_id = ?", userID).
		Select("COALESCE(SUM(word_count), 0)").
		Scan(&wordCount)
	stats["wordCount"] = wordCount

	// 获取收藏数 - 忽略错误
	var favoriteCount int64
	s.db.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&favoriteCount)
	stats["favoriteCount"] = favoriteCount

	return stats, nil
}

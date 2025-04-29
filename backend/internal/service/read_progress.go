package service

import (
	"ai-novel-platform/internal/models"

	"gorm.io/gorm"
)

type ReadProgressService struct {
	db *gorm.DB
}

func NewReadProgressService(db *gorm.DB) *ReadProgressService {
	return &ReadProgressService{db: db}
}

// UpdateReadProgress 更新用户阅读进度
func (s *ReadProgressService) UpdateReadProgress(userID, novelID, chapterID uint) error {
	progress := models.ReadProgress{}

	// 查找是否存在记录
	result := s.db.Where("user_id = ? AND novel_id = ?", userID, novelID).First(&progress)

	if result.Error == gorm.ErrRecordNotFound {
		// 不存在则创建新记录
		progress = models.ReadProgress{
			UserID:    userID,
			NovelID:   novelID,
			ChapterID: chapterID,
		}
		return s.db.Create(&progress).Error
	} else if result.Error != nil {
		return result.Error
	}

	// 存在则更新
	progress.ChapterID = chapterID
	return s.db.Save(&progress).Error
}

// GetReadProgress 获取用户阅读进度
func (s *ReadProgressService) GetReadProgress(userID, novelID uint) (*models.ReadProgress, error) {
	var progress models.ReadProgress
	err := s.db.Where("user_id = ? AND novel_id = ?", userID, novelID).First(&progress).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &progress, nil
}

package service

import (
	"ai-novel-platform/internal/models"
	"errors"
	"gorm.io/gorm"
)

type ChapterService struct {
	db *gorm.DB
}

func NewChapterService(db *gorm.DB) *ChapterService {
	return &ChapterService{db: db}
}

// CreateChapter 创建新章节
func (s *ChapterService) CreateChapter(chapter *models.Chapter) error {
	// 获取当前最大的order
	var maxOrder struct {
		MaxOrder int
	}
	s.db.Model(&models.Chapter{}).
		Where("novel_id = ?", chapter.NovelID).
		Select("COALESCE(MAX(`order`), 0) as max_order").
		Scan(&maxOrder)

	chapter.Order = maxOrder.MaxOrder + 1
	return s.db.Create(chapter).Error
}

// GetChapter 获取章节详情
func (s *ChapterService) GetChapter(id uint) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := s.db.First(&chapter, id).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

// UpdateChapter 更新章节
func (s *ChapterService) UpdateChapter(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Chapter{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteChapter 删除章节
func (s *ChapterService) DeleteChapter(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var chapter models.Chapter
		if err := tx.First(&chapter, id).Error; err != nil {
			return err
		}

		// 删除章节
		if err := tx.Delete(&chapter).Error; err != nil {
			return err
		}

		// 更新后续章节的顺序
		return tx.Model(&models.Chapter{}).
			Where("novel_id = ? AND `order` > ?", chapter.NovelID, chapter.Order).
			UpdateColumn("`order`", gorm.Expr("`order` - 1")).
			Error
	})
}

// ListNovelChapters 获取小说的章节列表
func (s *ChapterService) ListNovelChapters(novelID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := s.db.Where("novel_id = ?", novelID).
		Order("`order` asc").
		Find(&chapters).Error
	return chapters, err
}

// MoveChapter 移动章节位置
func (s *ChapterService) MoveChapter(id uint, direction string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var chapter models.Chapter
		if err := tx.First(&chapter, id).Error; err != nil {
			return err
		}

		var targetChapter models.Chapter
		var err error

		switch direction {
		case "up":
			if chapter.Order <= 1 {
				return errors.New("already at the top")
			}
			err = tx.Where("novel_id = ? AND `order` = ?", chapter.NovelID, chapter.Order-1).
				First(&targetChapter).Error
		case "down":
			err = tx.Where("novel_id = ? AND `order` = ?", chapter.NovelID, chapter.Order+1).
				First(&targetChapter).Error
		default:
			return errors.New("invalid direction")
		}

		if err != nil {
			return err
		}

		// 交换顺序
		chapter.Order, targetChapter.Order = targetChapter.Order, chapter.Order
		if err := tx.Save(&chapter).Error; err != nil {
			return err
		}
		return tx.Save(&targetChapter).Error
	})
}

// UpdateChapterContent 更新章节内容
func (s *ChapterService) UpdateChapterContent(id uint, title string, content string, wordCount int) error {
	updates := map[string]interface{}{
		"title":      title,
		"content":    content,
		"word_count": wordCount,
	}
	return s.db.Model(&models.Chapter{}).Where("id = ?", id).Updates(updates).Error
}

// GetChaptersByNovelID 获取小说的所有章节
func (s *ChapterService) GetChaptersByNovelID(novelID uint) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := s.db.Where("novel_id = ?", novelID).
		Order("`order` asc").
		Find(&chapters).Error
	return chapters, err
}

// UpdateChapterStatus 更新章节状态
func (s *ChapterService) UpdateChapterStatus(id uint, status int) error {
	return s.db.Model(&models.Chapter{}).Where("id = ?", id).Update("status", status).Error
} 
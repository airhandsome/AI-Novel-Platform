package service

import (
	"ai-novel-platform/internal/models"
	"errors"

	"gorm.io/gorm"
)

type NovelService struct {
	db *gorm.DB
}

func NewNovelService(db *gorm.DB) *NovelService {
	return &NovelService{db: db}
}

func (s *NovelService) CreateNovel(novel *models.Novel) error {
	return s.db.Create(novel).Error
}

func (s *NovelService) GetNovel(id uint) (*models.Novel, error) {
	var novel models.Novel
	if err := s.db.Preload("Author").First(&novel, id).Error; err != nil {
		return nil, err
	}
	return &novel, nil
}

func (s *NovelService) UpdateNovel(novel *models.Novel) error {
	return s.db.Save(novel).Error
}

func (s *NovelService) DeleteNovel(id uint, authorID uint) error {
	result := s.db.Where("id = ? AND author_id = ?", id, authorID).Delete(&models.Novel{})
	if result.RowsAffected == 0 {
		return errors.New("novel not found or not authorized")
	}
	return result.Error
}

func (s *NovelService) ListNovels(page, pageSize int, category string, status int, keyword string) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	query := s.db.Model(&models.Novel{}).Preload("Author")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != -1 {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&novels).Error; err != nil {
		return nil, 0, err
	}

	return novels, total, nil
}

func (s *NovelService) ListAuthorNovels(authorID uint, page, pageSize int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	query := s.db.Model(&models.Novel{}).Where("author_id = ?", authorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&novels).Error; err != nil {
		return nil, 0, err
	}

	return novels, total, nil
}

// GetAuthorStats 获取作者统计数据
func (s *NovelService) GetAuthorStats(authorID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 获取小说总数
	var novelCount int64
	if err := s.db.Model(&models.Novel{}).Where("author_id = ?", authorID).Count(&novelCount).Error; err != nil {
		return nil, err
	}
	stats["totalNovels"] = novelCount

	// 获取总字数
	var totalWords int64
	if err := s.db.Model(&models.Novel{}).Where("author_id = ?", authorID).Select("COALESCE(SUM(word_count), 0)").Scan(&totalWords).Error; err != nil {
		return nil, err
	}
	stats["totalWords"] = totalWords

	// 获取总阅读量
	var totalReads int64
	if err := s.db.Model(&models.Novel{}).Where("author_id = ?", authorID).Select("COALESCE(SUM(read_count), 0)").Scan(&totalReads).Error; err != nil {
		return nil, err
	}
	stats["totalReads"] = totalReads

	// 获取总收藏数
	var totalFavorites int64
	if err := s.db.Model(&models.Novel{}).Where("author_id = ?", authorID).Select("COALESCE(SUM(favorite_count), 0)").Scan(&totalFavorites).Error; err != nil {
		return nil, err
	}
	stats["totalFavorites"] = totalFavorites

	// 获取连载中的小说数量
	var ongoingNovels int64
	if err := s.db.Model(&models.Novel{}).Where("author_id = ? AND status = ?", authorID, 1).Count(&ongoingNovels).Error; err != nil {
		return nil, err
	}
	stats["ongoingNovels"] = ongoingNovels

	// 获取已完结的小说数量
	var completedNovels int64
	if err := s.db.Model(&models.Novel{}).Where("author_id = ? AND status = ?", authorID, 2).Count(&completedNovels).Error; err != nil {
		return nil, err
	}
	stats["completedNovels"] = completedNovels

	return stats, nil
}

// GetNovelOutline 获取小说大纲
func (s *NovelService) GetNovelOutline(novelID uint, authorID uint) (*models.NovelOutline, error) {
	var novel models.Novel
	if err := s.db.Where("id = ? AND author_id = ?", novelID, authorID).First(&novel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.NovelOutline{
				Outline: []models.OutlineItem{},
				WorldBuilding: models.WorldBuilding{
					Background:  "",
					Characters:  []models.Character{},
					Locations:  []models.Location{},
				},
			}, nil
		}
		return nil, err
	}
	
	if novel.NovelOutline.Outline == nil {
		novel.NovelOutline.Outline = []models.OutlineItem{}
	}
	if novel.NovelOutline.WorldBuilding.Characters == nil {
		novel.NovelOutline.WorldBuilding.Characters = []models.Character{}
	}
	if novel.NovelOutline.WorldBuilding.Locations == nil {
		novel.NovelOutline.WorldBuilding.Locations = []models.Location{}
	}
	
	return &novel.NovelOutline, nil
}

// UpdateNovelOutline 更新小说大纲
func (s *NovelService) UpdateNovelOutline(novelID uint, authorID uint, outline *models.NovelOutline) error {
	// 确保切片不为 nil
	if outline.Outline == nil {
		outline.Outline = []models.OutlineItem{}
	}
	if outline.WorldBuilding.Characters == nil {
		outline.WorldBuilding.Characters = []models.Character{}
	}
	if outline.WorldBuilding.Locations == nil {
		outline.WorldBuilding.Locations = []models.Location{}
	}

	result := s.db.Model(&models.Novel{}).
		Where("id = ? AND author_id = ?", novelID, authorID).
		Update("novel_outline", outline)
	
	if result.RowsAffected == 0 {
		return errors.New("novel not found or not authorized")
	}
	return result.Error
}

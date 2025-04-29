package service

import (
	"ai-novel-platform/internal/models"

	"gorm.io/gorm"
)

// AuthorStats 作者统计数据结构
type AuthorStats struct {
	TotalNovels   int64   `json:"total_novels"`   // 总小说数
	TotalChapters int64   `json:"total_chapters"` // 总章节数
	TotalWords    int64   `json:"total_words"`    // 总字数
	TotalViews    int64   `json:"total_views"`    // 总浏览量
	TotalLikes    int64   `json:"total_likes"`    // 总点赞数
	TotalComments int64   `json:"total_comments"` // 总评论数
	AverageRating float64 `json:"average_rating"` // 平均评分
}

// GetAuthorStats 获取作者统计数据
func (s *NovelService) GetNovelStats(novelID uint) (*AuthorStats, error) {
	stats := &AuthorStats{}

	// 获取总小说数
	if err := s.db.Model(&models.Novel{}).Where("id = ?", novelID).Count(&stats.TotalNovels).Error; err != nil {
		return nil, err
	}

	// 获取总章节数和总字数
	var result struct {
		TotalChapters int64
		TotalWords    int64
	}
	if err := s.db.Model(&models.Chapter{}).
		Select("COUNT(*) as total_chapters, SUM(word_count) as total_words").
		Joins("JOIN novels ON chapters.novel_id = novels.id").
		Where("novels.id = ?", novelID).
		Scan(&result).Error; err != nil {
		return nil, err
	}
	stats.TotalChapters = result.TotalChapters
	stats.TotalWords = result.TotalWords

	// 获取总浏览量、点赞数和评论数
	if err := s.db.Model(&models.Novel{}).
		Select("SUM(views) as total_views, SUM(likes) as total_likes, SUM(comments) as total_comments").
		Where("id = ?", novelID).
		Scan(&struct {
			TotalViews    int64
			TotalLikes    int64
			TotalComments int64
		}{
			TotalViews:    stats.TotalViews,
			TotalLikes:    stats.TotalLikes,
			TotalComments: stats.TotalComments,
		}).Error; err != nil {
		return nil, err
	}

	// 获取平均评分
	if err := s.db.Model(&models.Novel{}).
		Select("COALESCE(AVG(rating), 0) as average_rating").
		Where("id = ?", novelID).
		Scan(&stats.AverageRating).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// FavoriteNovel 收藏小说
func (s *NovelService) FavoriteNovel(userID, novelID uint) error {
	// 检查小说是否存在
	var novel models.Novel
	if err := s.db.First(&novel, novelID).Error; err != nil {
		return err
	}

	// 检查是否已经收藏
	var favorite models.Favorite
	result := s.db.Where("user_id = ? AND novel_id = ?", userID, novelID).First(&favorite)
	if result.Error == nil {
		return nil // 已经收藏，直接返回
	}

	// 创建收藏记录
	favorite = models.Favorite{
		UserID:  userID,
		NovelID: novelID,
	}
	if err := s.db.Create(&favorite).Error; err != nil {
		return err
	}

	// 更新小说的收藏数
	if err := s.db.Model(&novel).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
		return err
	}

	return nil
}

// UnfavoriteNovel 取消收藏
func (s *NovelService) UnfavoriteNovel(userID, novelID uint) error {
	// 检查小说是否存在
	var novel models.Novel
	if err := s.db.First(&novel, novelID).Error; err != nil {
		return err
	}

	// 删除收藏记录
	result := s.db.Where("user_id = ? AND novel_id = ?", userID, novelID).Delete(&models.Favorite{})
	if result.Error != nil {
		return result.Error
	}

	// 如果确实删除了记录，更新小说的收藏数
	if result.RowsAffected > 0 {
		if err := s.db.Model(&novel).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetFavorites 获取用户的收藏列表
func (s *NovelService) GetFavorites(userID uint, page, limit int, novelID uint) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// 构建查询条件
	query := s.db.Model(&models.Favorite{}).Where("user_id = ?", userID)
	if novelID > 0 {
		query = query.Where("novel_id = ?", novelID)
	}

	// 获取收藏的小说ID列表
	var favoriteIDs []uint
	if err := query.Pluck("novel_id", &favoriteIDs).Error; err != nil {
		return nil, 0, err
	}

	// 获取小说详情
	novelQuery := s.db.Model(&models.Novel{}).
		Where("id IN ?", favoriteIDs).
		Preload("Author")

	// 获取总数
	if err := novelQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := novelQuery.Offset(offset).Limit(limit).Find(&novels).Error; err != nil {
		return nil, 0, err
	}

	return novels, total, nil
}

// IsFavorited 检查用户是否收藏了小说
func (s *NovelService) IsFavorited(userID, novelID uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.Favorite{}).
		Where("user_id = ? AND novel_id = ?", userID, novelID).
		Count(&count).Error
	return count > 0, err
}

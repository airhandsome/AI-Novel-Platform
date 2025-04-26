package service

import "ai-novel-platform/internal/models"

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

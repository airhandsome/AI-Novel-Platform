package handlers

import (
	"ai-novel-platform/internal/models"
	"ai-novel-platform/internal/service"
	"ai-novel-platform/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ChapterHandler struct {
	chapterService *service.ChapterService
	novelService   *service.NovelService
}

func NewChapterHandler(chapterService *service.ChapterService, novelService *service.NovelService) *ChapterHandler {
	return &ChapterHandler{
		chapterService: chapterService,
		novelService:   novelService,
	}
}

// CreateChapter 创建新章节
func (h *ChapterHandler) CreateChapter(c *gin.Context) {
	var chapter models.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证小说所有权
	userID := utils.GetUserIDFromContext(c)
	novel, err := h.novelService.GetNovel(chapter.NovelID)
	if err != nil || novel.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此小说"})
		return
	}

	if err := h.chapterService.CreateChapter(&chapter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chapter)
}

// GetChapter 获取章节详情
func (h *ChapterHandler) GetChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	chapter, err := h.chapterService.GetChapter(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	c.JSON(http.StatusOK, chapter)
}

// UpdateChapter 更新章节
func (h *ChapterHandler) UpdateChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var updates struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		WordCount int    `json:"wordCount"`
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.chapterService.UpdateChapterContent(uint(id), updates.Title, updates.Content, updates.WordCount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter updated successfully"})
}

// DeleteChapter 删除章节
func (h *ChapterHandler) DeleteChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	if err := h.chapterService.DeleteChapter(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter deleted successfully"})
}

// ListNovelChapters 获取小说的章节列表
func (h *ChapterHandler) ListNovelChapters(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("novelId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid novel ID"})
		return
	}

	chapters, err := h.chapterService.ListNovelChapters(uint(novelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chapters": chapters})
}

// MoveChapter 移动章节位置
func (h *ChapterHandler) MoveChapter(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	direction := c.Query("direction")
	if direction != "up" && direction != "down" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid direction"})
		return
	}

	if err := h.chapterService.MoveChapter(uint(id), direction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter moved successfully"})
}

// UpdateChapterStatus 更新章节状态
func (h *ChapterHandler) UpdateChapterStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chapter ID"})
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证小说所有权
	chapter, err := h.chapterService.GetChapter(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}

	userID := utils.GetUserIDFromContext(c)
	novel, err := h.novelService.GetNovel(chapter.NovelID)
	if err != nil || novel.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此章节"})
		return
	}

	if err := h.chapterService.UpdateChapterStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chapter status updated successfully"})
} 
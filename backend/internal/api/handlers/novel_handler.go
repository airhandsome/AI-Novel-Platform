package handlers

import (
	"ai-novel-platform/internal/models"
	"ai-novel-platform/internal/service"
	"ai-novel-platform/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type NovelHandler struct {
	novelService *service.NovelService
}

func NewNovelHandler(novelService *service.NovelService) *NovelHandler {
	return &NovelHandler{novelService: novelService}
}

func (h *NovelHandler) CreateNovel(c *gin.Context) {
	var novel models.Novel
	if err := c.ShouldBindJSON(&novel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	userID := utils.GetUserIDFromContext(c)
	novel.AuthorID = userID

	if err := h.novelService.CreateNovel(&novel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, novel)
}

func (h *NovelHandler) GetNovel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid novel ID"})
		return
	}

	novel, err := h.novelService.GetNovel(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found"})
		return
	}

	c.JSON(http.StatusOK, novel)
}

func (h *NovelHandler) UpdateNovel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid novel ID"})
		return
	}

	var novel models.Novel
	if err := c.ShouldBindJSON(&novel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exist, err := h.novelService.GetNovel(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found"})
		return
	}
	exist.UpdatedAt = time.Now()
	exist.Title = novel.Title
	exist.CoverURL = novel.CoverURL
	exist.Description = novel.Description
	exist.Category = novel.Category
	exist.Status = novel.Status
	exist.Tags = novel.Tags

	if err := h.novelService.UpdateNovel(exist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, novel)
}

func (h *NovelHandler) DeleteNovel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid novel ID"})
		return
	}

	userID := utils.GetUserIDFromContext(c)
	if err := h.novelService.DeleteNovel(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Novel deleted successfully"})
}

func (h *NovelHandler) ListNovels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))
	keyword := c.Query("keyword")

	novels, total, err := h.novelService.ListNovels(page, pageSize, category, status, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"novels": novels,
		"total":  total,
	})
}

func (h *NovelHandler) ListAuthorNovels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	userID := utils.GetUserIDFromContext(c)

	novels, total, err := h.novelService.ListAuthorNovels(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"novels": novels,
		"total":  total,
	})
}

// GetAuthorStats 获取作者统计数据
func (h *NovelHandler) GetAuthorStats(c *gin.Context) {
	// 从上下文中获取作者ID
	authorID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 调用service层获取统计数据
	stats, err := h.novelService.GetAuthorStats(authorID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": stats,
	})
}

// GetNovelOutline 获取小说大纲
func (h *NovelHandler) GetNovelOutline(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid novel ID",
			"data": models.NovelOutline{
				Outline: []models.OutlineItem{},
				WorldBuilding: models.WorldBuilding{},
			},
		})
		return
	}

	userID := utils.GetUserIDFromContext(c)
	outline, err := h.novelService.GetNovelOutline(uint(novelID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  err.Error(),
			"data": models.NovelOutline{
				Outline: []models.OutlineItem{},
				WorldBuilding: models.WorldBuilding{},
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": outline,
	})
}

// UpdateNovelOutline 更新小说大纲
func (h *NovelHandler) UpdateNovelOutline(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid novel ID",
		})
		return
	}

	var outline models.NovelOutline
	if err := c.ShouldBindJSON(&outline); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Invalid outline data",
			"error": err.Error(),
		})
		return
	}

	userID := utils.GetUserIDFromContext(c)
	if err := h.novelService.UpdateNovelOutline(uint(novelID), userID, &outline); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Failed to update outline",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

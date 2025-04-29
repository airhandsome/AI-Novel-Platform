package handlers

import (
	"ai-novel-platform/internal/service"
	"ai-novel-platform/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReadProgressHandler struct {
	readProgressService *service.ReadProgressService
}

func NewReadProgressHandler(readProgressService *service.ReadProgressService) *ReadProgressHandler {
	return &ReadProgressHandler{readProgressService: readProgressService}
}

// UpdateReadProgress 更新阅读进度
func (h *ReadProgressHandler) UpdateReadProgress(c *gin.Context) {
	var req struct {
		NovelID   uint `json:"novelId" binding:"required"`
		ChapterID uint `json:"chapterId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	// 从上下文获取用户ID
	userID := utils.GetUserIDFromContext(c)

	err := h.readProgressService.UpdateReadProgress(userID, req.NovelID, req.ChapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reading progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reading progress updated successfully"})
}

// GetReadProgress 获取阅读进度
func (h *ReadProgressHandler) GetReadProgress(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("novelId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "can't get novelId"})
		return
	}
	userID := utils.GetUserIDFromContext(c)

	progress, err := h.readProgressService.GetReadProgress(userID, uint(novelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reading progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": progress})
}

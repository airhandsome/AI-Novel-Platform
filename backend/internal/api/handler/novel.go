package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/your-project/model"
	"github.com/your-project/service"
)

// GetAuthorStats 获取作者统计数据
func (h *NovelHandler) GetAuthorStats(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"message": "未授权访问",
		})
		return
	}

	// 获取作者统计数据
	stats, err := h.novelService.GetAuthorStats(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"message": "获取统计数据失败",
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"message": "success",
		"data": stats,
	})
} 
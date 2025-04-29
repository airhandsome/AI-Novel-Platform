package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteNovel 收藏小说
func (h *NovelHandler) FavoriteNovel(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的小说ID"})
		return
	}

	userID := c.GetUint("userID")
	if err := h.novelService.FavoriteNovel(userID, uint(novelID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "收藏成功"})
}

// UnfavoriteNovel 取消收藏
func (h *NovelHandler) UnfavoriteNovel(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的小说ID"})
		return
	}

	userID := c.GetUint("userID")
	if err := h.novelService.UnfavoriteNovel(userID, uint(novelID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消收藏失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消收藏成功"})
}

// GetFavorites 获取收藏列表
func (h *NovelHandler) GetFavorites(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	novelID, _ := strconv.ParseUint(c.DefaultQuery("novelId", "0"), 10, 32)
	userID := c.GetUint("userID")

	novels, total, err := h.novelService.GetFavorites(userID, page, limit, uint(novelID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取收藏列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"novels": novels,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}
func (h *NovelHandler) CheckFavorite(c *gin.Context) {
	novelID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的小说ID"})
		return
	}

	userID := c.GetUint("userID")
	favorited := false
	if favorited, err = h.novelService.IsFavorited(userID, uint(novelID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "收藏失败"})
		return
	}

	if favorited {
		c.JSON(http.StatusOK, gin.H{"message": "true"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "false"})
	}

}

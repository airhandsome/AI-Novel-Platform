package router

import (
	"github.com/gin-gonic/gin"
	"github.com/your-project/middleware"
	"github.com/your-project/handler"
)

func SetupRouter(r *gin.Engine) {
	// 小说相关路由
	novelGroup := r.Group("/api/novel")
	{
		novelGroup.GET("/author/stats", middleware.AuthMiddleware(), novelHandler.GetAuthorStats) // 获取作者统计数据
	}
} 
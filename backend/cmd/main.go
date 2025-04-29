package main

import (
	"ai-novel-platform/internal/api/handlers"
	"ai-novel-platform/internal/middleware"
	"ai-novel-platform/internal/models"
	"ai-novel-platform/internal/service"
	"ai-novel-platform/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	db, err := utils.InitDB("novel_user", "novel_password", "localhost", "3306", "novel_platform")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 删除现有表（仅在开发环境使用）
	//db.Migrator().DropTable(&models.Novel{}, &models.User{}, &models.Favorite{}, &models.Chapter{})

	// 自动迁移数据库表
	if err := db.AutoMigrate(&models.User{}, &models.Novel{}, &models.Favorite{}, &models.Chapter{}, &models.ReadProgress{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 确保 novel_outline 字段类型正确
	if err := db.Exec("ALTER TABLE novels MODIFY novel_outline JSON").Error; err != nil {
		log.Printf("Warning: Failed to modify novel_outline column type: %v", err)
	}

	// 初始化Redis连接
	rdb, err := utils.InitRedis("localhost", "6379", "", 0)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	// 初始化Gin
	r := gin.Default()

	// 使用 CORS 中间件
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
		}
		c.Next()
	})

	// 初始化服务和处理器
	userService := service.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService, "your_jwt_secret")
	novelService := service.NewNovelService(db)
	novelHandler := handlers.NewNovelHandler(novelService)
	chapterService := service.NewChapterService(db)
	chapterHandler := handlers.NewChapterHandler(chapterService, novelService)
	readProgressService := service.NewReadProgressService(db)
	readProgressHandler := handlers.NewReadProgressHandler(readProgressService)

	// 用户相关路由
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/reset-password-request", userHandler.RequestPasswordReset)
		auth.POST("/reset-password", userHandler.ResetPassword)
		auth.POST("/refresh-token", middleware.JWTAuth(), userHandler.RefreshToken)
	}

	// 需要认证的路由
	user := r.Group("/api/v1/user")
	user.Use(middleware.JWTAuth())
	{
		user.GET("/profile", userHandler.GetProfile)
		user.PUT("/profile", userHandler.UpdateUser)
		user.PUT("/password", userHandler.UpdatePassword)
		user.POST("/avatar", userHandler.UploadAvatar)
	}

	// 小说相关路由
	novels := r.Group("/api/v1/novels")
	{
		// 公开接口
		novels.GET("", novelHandler.ListNovels)
		novels.GET("/:id", novelHandler.GetNovel)

		// 需要认证的接口
		authorized := novels.Group("")
		authorized.Use(middleware.JWTAuth())
		{
			authorized.POST("", novelHandler.CreateNovel)
			authorized.PUT("/:id", novelHandler.UpdateNovel)
			authorized.PUT("/:id/status", novelHandler.UpdateNovelStatus)
			authorized.DELETE("/:id", novelHandler.DeleteNovel)
			authorized.GET("/author/me", novelHandler.ListAuthorNovels)
			authorized.GET("/author/stats", novelHandler.GetAuthorStats)
			authorized.GET("/:id/outline", novelHandler.GetNovelOutline)
			authorized.PUT("/:id/outline", novelHandler.UpdateNovelOutline)
			authorized.GET("/favorite/:id", novelHandler.CheckFavorite)
			authorized.POST("/favorite/:id", novelHandler.FavoriteNovel)
			authorized.DELETE("/favorite/:id", novelHandler.UnfavoriteNovel)
			authorized.GET("/favorites", novelHandler.GetFavorites)
		}
	}

	// 章节相关路由
	chapters := r.Group("/api/v1/chapters")
	chapters.Use(middleware.JWTAuth())
	{
		chapters.POST("", chapterHandler.CreateChapter)
		chapters.GET("/:id", chapterHandler.GetChapter)
		chapters.PUT("/:id", chapterHandler.UpdateChapter)
		chapters.DELETE("/:id", chapterHandler.DeleteChapter)
		chapters.PUT("/:id/move", chapterHandler.MoveChapter)
		chapters.GET("/novel/:novelId", chapterHandler.ListNovelChapters)
		chapters.PUT("/:id/status", chapterHandler.UpdateChapterStatus)
	}

	// 阅读进度相关路由
	progress := r.Group("/api/v1/reading-progress")
	progress.Use(middleware.JWTAuth())
	{
		progress.POST("", readProgressHandler.UpdateReadProgress)
		progress.GET("/novel/:novelId", readProgressHandler.GetReadProgress)
	}

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

package router

import (
	"x-novel/internal/api/handler"
	"x-novel/internal/api/middleware"
	"x-novel/internal/repository"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	r *gin.Engine,
	deviceRepo *repository.DeviceRepository,
	deviceHandler *handler.DeviceHandler,
	projectHandler *handler.ProjectHandler,
	chapterHandler *handler.ChapterHandler,
	modelConfigHandler *handler.ModelConfigHandler,
	chatHandler *handler.ChatHandler,
	writingAssistantHandler *handler.WritingAssistantHandler,
	graphHandler *handler.GraphHandler,
	reviewHandler *handler.ReviewHandler,
	backupHandler *handler.BackupHandler,
) {
	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.Device(deviceRepo))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 设备相关
		device := v1.Group("/device")
		{
			device.GET("/info", deviceHandler.GetInfo)
			device.GET("/settings", deviceHandler.GetSettings)
			device.PUT("/settings", deviceHandler.UpdateSettings)
		}

		// 模型配置
		models := v1.Group("/models")
		{
			models.GET("", modelConfigHandler.List)
			models.POST("", modelConfigHandler.Create)
			models.GET("/providers", modelConfigHandler.ListProviders)
			models.POST("/validate", modelConfigHandler.Validate)
			models.GET("/:id", modelConfigHandler.GetByID)
			models.PUT("/:id", modelConfigHandler.Update)
			models.DELETE("/:id", modelConfigHandler.Delete)
		}

		// 项目
		projects := v1.Group("/projects")
		{
			projects.GET("", projectHandler.List)
			projects.POST("", projectHandler.Create)
			projects.GET("/:id", projectHandler.GetByID)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)

			// 项目子资源
			projects.GET("/:id/chapters", chapterHandler.List)
			projects.POST("/:id/chapters", chapterHandler.Create)
			projects.GET("/:id/chapters/:chapterNumber", chapterHandler.GetByNumber)
			projects.PUT("/:id/chapters/:chapterNumber", chapterHandler.Update)
			projects.POST("/:id/chapters/:chapterNumber/generate", chapterHandler.GenerateContent)
			projects.POST("/:id/chapters/:chapterNumber/finalize", chapterHandler.Finalize)
			projects.POST("/:id/chapters/:chapterNumber/enrich", chapterHandler.Enrich)

			// 架构生成
			projects.POST("/:id/architecture/generate", projectHandler.GenerateArchitecture)

			// 大纲生成
			projects.POST("/:id/blueprint/generate", projectHandler.GenerateBlueprint)

			// 导出
			projects.GET("/:id/export/:format", projectHandler.ExportProject)

			// 关系图谱
			projects.GET("/:id/graph", graphHandler.GetGraph)
			projects.POST("/:id/graph/generate", graphHandler.GenerateGraph)
			projects.POST("/:id/graph/chapters/:chapterNumber", graphHandler.UpdateFromChapter)
			projects.GET("/:id/graph/chapters/:chapterNumber", graphHandler.GetChapterSnapshot)
		}

		// 写作助手
		v1.POST("/writing/assist", writingAssistantHandler.Assist)

		// 错误检测 & AI 审阅 & 市场预测
		v1.POST("/review/detect", reviewHandler.DetectErrors)
		v1.POST("/projects/:id/review", reviewHandler.ReviewProject)
		v1.POST("/projects/:id/review/chapters/:chapterNumber", reviewHandler.ReviewChapter)
		v1.POST("/projects/:id/market-predict", reviewHandler.MarketPredict)

		// 数据备份
		backup := v1.Group("/backup")
		{
			backup.GET("/preview", backupHandler.ExportPreview)
			backup.GET("/export", backupHandler.Export)
			backup.POST("/import", backupHandler.Import)
		}

		// 对话
		chat := v1.Group("/conversations")
		{
			chat.GET("", chatHandler.ListConversations)
			chat.POST("", chatHandler.CreateConversation)
			chat.GET("/:id", chatHandler.GetConversation)
			chat.PUT("/:id", chatHandler.UpdateConversation)
			chat.DELETE("/:id", chatHandler.DeleteConversation)
			chat.POST("/:id/messages", chatHandler.SendMessage)
		}
	}
}

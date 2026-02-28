package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"x-novel/internal/api/handler"
	"x-novel/internal/api/router"
	"x-novel/internal/config"
	"x-novel/internal/llm"
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/internal/service"
	"x-novel/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync()

	logger.Info("启动 X-Novel 服务器",
		zap.String("version", "1.0.0"),
		zap.String("mode", cfg.Server.Mode),
	)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("连接数据库失败", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("获取数据库连接失败", zap.Error(err))
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("数据库连接成功",
		zap.String("host", cfg.Database.Host),
		zap.String("dbname", cfg.Database.DBName),
	)

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		logger.Fatal("数据库迁移失败", zap.Error(err))
	}

	// 初始化仓储
	deviceRepo := repository.NewDeviceRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	chapterRepo := repository.NewChapterRepository(db)
	modelConfigRepo := repository.NewModelConfigRepository(db)

	// 初始化 LLM 管理器
	llmManager := llm.NewManager()
	// 注册默认适配器
	llmManager.Register("openai", llm.NewOpenAIAdapter("", "gpt-3.5-turbo"))

	// 初始化服务
	deviceService := service.NewDeviceService(deviceRepo)
	exportService := service.NewExportService(projectRepo, chapterRepo)
	projectService := service.NewProjectService(projectRepo, chapterRepo, modelConfigRepo, llmManager, exportService)
	chapterService := service.NewChapterService(projectRepo, chapterRepo, modelConfigRepo, llmManager)

	// 初始化处理器
	deviceHandler := handler.NewDeviceHandler(deviceService)
	projectHandler := handler.NewProjectHandler(projectService)
	chapterHandler := handler.NewChapterHandler(chapterService, projectService)

	// 设置 Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 设置路由
	router.SetupRouter(r, deviceRepo, deviceHandler, projectHandler, chapterHandler)

	// 启动服务器
	srv := &http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		logger.Info("服务器启动",
			zap.String("addr", srv.Addr),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	// 关闭数据库连接
	if err := sqlDB.Close(); err != nil {
		logger.Error("数据库连接关闭失败", zap.Error(err))
	}

	logger.Info("服务器已关闭")
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	logger.Info("开始数据库迁移...")

	// 迁移所有模型
	err := db.AutoMigrate(
		&model.Device{},
		&model.DeviceSetting{},
		&model.Project{},
		&model.Chapter{},
		&model.ModelProvider{},
		&model.ModelConfig{},
	)

	if err != nil {
		return fmt.Errorf("数据库迁移失败: %w", err)
	}

	logger.Info("数据库迁移完成")
	return nil
}

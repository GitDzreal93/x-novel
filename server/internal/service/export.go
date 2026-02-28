package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"go.uber.org/zap"
)

// ExportService 导出服务
type ExportService struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
}

// NewExportService 创建导出服务
func NewExportService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
) *ExportService {
	return &ExportService{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
	}
}

// ExportFormat 导出格式
type ExportFormat string

const (
	FormatTXT      ExportFormat = "txt"
	FormatMarkdown ExportFormat = "md"
)

// ExportProject 导出项目
func (s *ExportService) ExportProject(ctx context.Context, projectID string, format ExportFormat) (string, error) {
	// 获取项目信息
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		logger.Error("获取项目失败", zap.Error(err))
		return "", err
	}

	// 获取所有章节
	chapters, err := s.chapterRepo.ListByProject(ctx, projectID)
	if err != nil {
		logger.Error("获取章节列表失败", zap.Error(err))
		return "", err
	}

	logger.Info("开始导出项目",
		zap.String("project_id", projectID),
		zap.String("format", string(format)),
		zap.Int("chapter_count", len(chapters)),
	)

	// 根据格式导出
	switch format {
	case FormatTXT:
		return s.exportToTXT(project, chapters)
	case FormatMarkdown:
		return s.exportToMarkdown(project, chapters)
	default:
		return "", fmt.Errorf("不支持的导出格式: %s", format)
	}
}

// exportToTXT 导出为纯文本格式
func (s *ExportService) exportToTXT(project *model.Project, chapters []*model.Chapter) (string, error) {
	var builder strings.Builder
	totalWords := 0

	// 计算总字数
	for _, chapter := range chapters {
		totalWords += chapter.WordCount
	}

	// 写入标题和分隔线
	builder.WriteString(project.Title + "\n")
	builder.WriteString(strings.Repeat("=", len(project.Title)) + "\n\n")

	// 写入项目信息
	builder.WriteString(fmt.Sprintf("主题：%s\n", project.Topic))
	builder.WriteString(fmt.Sprintf("类型：%s\n", project.Genre))
	builder.WriteString(fmt.Sprintf("章节总数：%d\n", len(chapters)))
	builder.WriteString(fmt.Sprintf("总字数：%d\n", totalWords))
	builder.WriteString(fmt.Sprintf("导出时间：%s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	builder.WriteString(strings.Repeat("-", 50) + "\n\n")

	// 写入章节
	for _, chapter := range chapters {
		builder.WriteString(fmt.Sprintf("第%d章 %s\n\n", chapter.ChapterNumber, chapter.Title))
		if chapter.Content != "" {
			builder.WriteString(chapter.Content)
			builder.WriteString("\n\n")
		}
		builder.WriteString(strings.Repeat("-", 30) + "\n\n")
	}

	return builder.String(), nil
}

// exportToMarkdown 导出为 Markdown 格式
func (s *ExportService) exportToMarkdown(project *model.Project, chapters []*model.Chapter) (string, error) {
	var builder strings.Builder
	totalWords := 0

	// 计算总字数
	for _, chapter := range chapters {
		totalWords += chapter.WordCount
	}

	// 写入标题
	builder.WriteString("# " + project.Title + "\n\n")

	// 写入元数据
	builder.WriteString("## 元数据\n\n")
	builder.WriteString(fmt.Sprintf("- **主题**：%s\n", project.Topic))
	builder.WriteString(fmt.Sprintf("- **类型**：%s\n", project.Genre))
	builder.WriteString(fmt.Sprintf("- **章节总数**：%d\n", len(chapters)))
	builder.WriteString(fmt.Sprintf("- **总字数**：%d\n", totalWords))
	builder.WriteString(fmt.Sprintf("- **导出时间**：%s\n", time.Now().Format("2006-01-02 15:04:05")))
	builder.WriteString("\n")

	// 写入章节目录
	builder.WriteString("## 目录\n\n")
	for _, chapter := range chapters {
		anchor := fmt.Sprintf("第%d章", chapter.ChapterNumber)
		builder.WriteString(fmt.Sprintf("%d. [第%d章 %s](#%s)\n", chapter.ChapterNumber, chapter.ChapterNumber, chapter.Title, anchor))
	}
	builder.WriteString("\n")

	// 写入章节内容
	for _, chapter := range chapters {
		builder.WriteString(fmt.Sprintf("## 第%d章 %s\n\n", chapter.ChapterNumber, chapter.Title))
		if chapter.Content != "" {
			builder.WriteString(chapter.Content)
			builder.WriteString("\n\n")
		}
	}

	return builder.String(), nil
}

package service

import (
	"context"
	"errors"
	"unicode/utf8"

	"x-novel/internal/dto"
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"go.uber.org/zap"
)

// ChapterService 章节服务
type ChapterService struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
}

// NewChapterService 创建章节服务
func NewChapterService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
) *ChapterService {
	return &ChapterService{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
	}
}

// Create 创建章节
func (s *ChapterService) Create(ctx context.Context, projectID string, req *dto.CreateChapterRequest) (*model.Chapter, error) {
	// 检查项目是否存在
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 检查章节号是否已存在
	existingChapter, _ := s.chapterRepo.GetByProjectAndNumber(ctx, projectID, req.ChapterNumber)
	if existingChapter != nil {
		return nil, errors.New("章节号已存在")
	}

	chapter := &model.Chapter{
		ProjectID:        project.ID,
		ChapterNumber:    req.ChapterNumber,
		Title:            req.Title,
		BlueprintSummary: req.BlueprintSummary,
		Status:           "not_started",
	}

	if err := s.chapterRepo.Create(ctx, chapter); err != nil {
		logger.Error("创建章节失败", zap.Error(err))
		return nil, err
	}

	return chapter, nil
}

// GetByID 根据 ID 获取章节
func (s *ChapterService) GetByID(ctx context.Context, id string) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("获取章节失败",
			zap.String("chapter_id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return chapter, nil
}

// GetByProjectAndNumber 根据项目 ID 和章节号获取章节
func (s *ChapterService) GetByProjectAndNumber(ctx context.Context, projectID string, chapterNumber int) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByProjectAndNumber(ctx, projectID, chapterNumber)
	if err != nil {
		logger.Error("获取章节失败",
			zap.String("project_id", projectID),
			zap.Int("chapter_number", chapterNumber),
			zap.Error(err),
		)
		return nil, err
	}
	return chapter, nil
}

// List 获取章节列表
func (s *ChapterService) List(ctx context.Context, projectID string, page, pageSize int) ([]*model.Chapter, int64, error) {
	offset := (page - 1) * pageSize
	chapters, total, err := s.chapterRepo.List(ctx, projectID, offset, pageSize)
	if err != nil {
		logger.Error("获取章节列表失败",
			zap.String("project_id", projectID),
			zap.Error(err),
		)
		return nil, 0, err
	}
	return chapters, total, nil
}

// Update 更新章节
func (s *ChapterService) Update(ctx context.Context, id string, req *dto.UpdateChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		chapter.Title = *req.Title
	}
	if req.Content != nil {
		chapter.Content = *req.Content
		chapter.WordCount = utf8.RuneCountInString(*req.Content)
	}
	if req.Status != nil {
		chapter.Status = *req.Status
	}
	if req.IsFinalized != nil {
		chapter.IsFinalized = *req.IsFinalized
	}

	// 更新大纲信息
	if req.BlueprintPosition != nil {
		chapter.BlueprintPosition = *req.BlueprintPosition
	}
	if req.BlueprintPurpose != nil {
		chapter.BlueprintPurpose = *req.BlueprintPurpose
	}
	if req.BlueprintSuspense != nil {
		chapter.BlueprintSuspense = *req.BlueprintSuspense
	}
	if req.BlueprintForeshadowing != nil {
		chapter.BlueprintForeshadowing = *req.BlueprintForeshadowing
	}
	if req.BlueprintTwistLevel != nil {
		chapter.BlueprintTwistLevel = *req.BlueprintTwistLevel
	}
	if req.BlueprintSummary != nil {
		chapter.BlueprintSummary = *req.BlueprintSummary
	}

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		logger.Error("更新章节失败",
			zap.String("chapter_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return chapter, nil
}

// Delete 删除章节
func (s *ChapterService) Delete(ctx context.Context, id string) error {
	if err := s.chapterRepo.Delete(ctx, id); err != nil {
		logger.Error("删除章节失败",
			zap.String("chapter_id", id),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GenerateChapterContent 生成章节内容
func (s *ChapterService) GenerateChapterContent(ctx context.Context, projectID, chapterID string, req *dto.GenerateChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 检查是否已有内容
	if chapter.Content != "" && !req.Overwrite {
		return nil, errors.New("章节已有内容，如需重新生成请设置 overwrite=true")
	}

	// TODO: 实现 LLM 调用生成章节内容
	logger.Info("开始生成章节内容",
		zap.String("project_id", projectID),
		zap.String("chapter_id", chapterID),
		zap.Int("chapter_number", req.ChapterNumber),
	)

	return chapter, nil
}

// FinalizeChapter 定稿章节
func (s *ChapterService) FinalizeChapter(ctx context.Context, projectID, chapterID string, req *dto.FinalizeChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 检查章节是否有内容
	if chapter.Content == "" {
		return nil, errors.New("章节内容为空，无法定稿")
	}

	// 设置为定稿状态
	chapter.IsFinalized = true
	chapter.Status = "completed"

	if err := s.chapterRepo.SetFinalized(ctx, chapterID, true); err != nil {
		return nil, err
	}

	// TODO: 如果需要更新全局摘要，在这里实现
	if req.UpdateSummary {
		// 更新项目的全局摘要
	}

	return chapter, nil
}

// EnrichChapter 扩写章节
func (s *ChapterService) EnrichChapter(ctx context.Context, projectID, chapterID string, req *dto.EnrichChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 检查章节是否有内容
	if chapter.Content == "" {
		return nil, errors.New("章节内容为空，无法扩写")
	}

	// TODO: 实现 LLM 调用扩写章节内容
	logger.Info("开始扩写章节内容",
		zap.String("project_id", projectID),
		zap.String("chapter_id", chapterID),
		zap.Int("target_words", req.TargetWords),
	)

	return chapter, nil
}

// GetPreviousChapters 获取前面已完成的所有章节（用于生成前文摘要）
func (s *ChapterService) GetPreviousChapters(ctx context.Context, projectID string, chapterNumber int) ([]*model.Chapter, error) {
	chapters, err := s.chapterRepo.ListCompleted(ctx, projectID, chapterNumber)
	if err != nil {
		logger.Error("获取前面章节失败",
			zap.String("project_id", projectID),
			zap.Int("chapter_number", chapterNumber),
			zap.Error(err),
		)
		return nil, err
	}
	return chapters, nil
}

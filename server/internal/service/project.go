package service

import (
	"context"
	"errors"

	"x-novel/internal/dto"
	"x-novel/internal/llm"
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ProjectService 项目服务
type ProjectService struct {
	projectRepo  *repository.ProjectRepository
	chapterRepo  *repository.ChapterRepository
	modelRepo    *repository.ModelConfigRepository
	llmManager   *llm.Manager
}

// NewProjectService 创建项目服务
func NewProjectService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
		modelRepo:   modelRepo,
		llmManager:  llmManager,
	}
}

// Create 创建项目
func (s *ProjectService) Create(ctx context.Context, deviceID uuid.UUID, req *dto.CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		DeviceID:         deviceID,
		Title:            req.Title,
		Topic:            req.Topic,
		Genre:            req.Genre,
		ChapterCount:     req.ChapterCount,
		WordsPerChapter:  req.WordsPerChapter,
		UserGuidance:     req.UserGuidance,
		Status:           "draft",
	}

	// 设置默认值
	if project.ChapterCount == 0 {
		project.ChapterCount = 100
	}
	if project.WordsPerChapter == 0 {
		project.WordsPerChapter = 3000
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		logger.Error("创建项目失败", zap.Error(err))
		return nil, err
	}

	return project, nil
}

// GetByID 根据 ID 获取项目
func (s *ProjectService) GetByID(ctx context.Context, id string) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("获取项目失败",
			zap.String("project_id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return project, nil
}

// List 获取项目列表
func (s *ProjectService) List(ctx context.Context, deviceID uuid.UUID, page, pageSize int) ([]*model.Project, int64, error) {
	offset := (page - 1) * pageSize
	projects, total, err := s.projectRepo.List(ctx, deviceID.String(), offset, pageSize)
	if err != nil {
		logger.Error("获取项目列表失败",
			zap.String("device_id", deviceID.String()),
			zap.Error(err),
		)
		return nil, 0, err
	}
	return projects, total, nil
}

// Update 更新项目
func (s *ProjectService) Update(ctx context.Context, id string, req *dto.UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		project.Title = *req.Title
	}
	if req.Topic != nil {
		project.Topic = *req.Topic
	}
	if req.Genre != nil {
		project.Genre = req.Genre
	}
	if req.ChapterCount != nil {
		project.ChapterCount = *req.ChapterCount
	}
	if req.WordsPerChapter != nil {
		project.WordsPerChapter = *req.WordsPerChapter
	}
	if req.UserGuidance != nil {
		project.UserGuidance = *req.UserGuidance
	}
	if req.Status != nil {
		project.Status = *req.Status
	}

	// 更新架构数据
	if req.CoreSeed != nil {
		project.CoreSeed = *req.CoreSeed
	}
	if req.CharacterDynamics != nil {
		project.CharacterDynamics = *req.CharacterDynamics
	}
	if req.WorldBuilding != nil {
		project.WorldBuilding = *req.WorldBuilding
	}
	if req.PlotArchitecture != nil {
		project.PlotArchitecture = *req.PlotArchitecture
	}
	if req.CharacterState != nil {
		project.CharacterState = *req.CharacterState
	}

	// 更新大纲数据
	if req.ChapterBlueprint != nil {
		project.ChapterBlueprint = *req.ChapterBlueprint
	}

	// 更新全局摘要
	if req.GlobalSummary != nil {
		project.GlobalSummary = *req.GlobalSummary
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		logger.Error("更新项目失败",
			zap.String("project_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return project, nil
}

// Delete 删除项目
func (s *ProjectService) Delete(ctx context.Context, id string) error {
	if err := s.projectRepo.Delete(ctx, id); err != nil {
		logger.Error("删除项目失败",
			zap.String("project_id", id),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GenerateArchitecture 生成小说架构
func (s *ProjectService) GenerateArchitecture(ctx context.Context, deviceID uuid.UUID, projectID string, req *dto.GenerateArchitectureRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 检查是否已经生成过
	if project.ArchitectureGenerated && !req.Overwrite {
		return nil, errors.New("架构已生成，如需重新生成请设置 overwrite=true")
	}

	// TODO: 实现 LLM 调用生成架构
	// 这里需要：
	// 1. 获取用于架构生成的模型配置
	// 2. 构建 5 个步骤的提示词
	// 3. 调用 LLM 生成
	// 4. 保存结果

	logger.Info("开始生成小说架构",
		zap.String("project_id", projectID),
	)

	// 暂时返回项目本身
	return project, nil
}

// GenerateBlueprint 生成章节大纲
func (s *ProjectService) GenerateBlueprint(ctx context.Context, deviceID uuid.UUID, projectID string, req *dto.GenerateBlueprintRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 检查架构是否已生成
	if !project.ArchitectureGenerated {
		return nil, errors.New("请先生成小说架构")
	}

	// 检查是否已经生成过
	if project.BlueprintGenerated && !req.Overwrite {
		return nil, errors.New("大纲已生成，如需重新生成请设置 overwrite=true")
	}

	// TODO: 实现 LLM 调用生成大纲
	logger.Info("开始生成章节大纲",
		zap.String("project_id", projectID),
	)

	return project, nil
}

// ExportProject 导出项目
func (s *ProjectService) ExportProject(ctx context.Context, projectID, format string) (string, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return "", err
	}

	// TODO: 实现导出逻辑
	_ = project
	_ = format

	return "", errors.New("导出功能待实现")
}

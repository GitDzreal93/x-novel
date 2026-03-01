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

// ModelConfigService 模型配置服务
type ModelConfigService struct {
	modelRepo  *repository.ModelConfigRepository
	llmManager *llm.Manager
}

// NewModelConfigService 创建模型配置服务
func NewModelConfigService(
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
) *ModelConfigService {
	return &ModelConfigService{
		modelRepo:  modelRepo,
		llmManager: llmManager,
	}
}

// Create 创建模型配置
func (s *ModelConfigService) Create(ctx context.Context, deviceID uuid.UUID, req *dto.CreateModelConfigRequest) (*model.ModelConfig, error) {
	// 验证提供商是否存在
	provider, err := s.modelRepo.GetProviderByID(ctx, req.ProviderID)
	if err != nil {
		logger.Error("提供商不存在", zap.Int("provider_id", req.ProviderID), zap.Error(err))
		return nil, errors.New("提供商不存在")
	}

	config := &model.ModelConfig{
		DeviceID:   deviceID,
		ProviderID: req.ProviderID,
		ModelName:  req.ModelName,
		Purpose:    req.Purpose,
		APIKey:     req.APIKey,
		BaseURL:    req.BaseURL,
		IsActive:   true,
	}

	if err := s.modelRepo.Create(ctx, config); err != nil {
		logger.Error("创建模型配置失败", zap.Error(err))
		return nil, err
	}

	// 设置 Provider
	config.Provider = provider

	return config, nil
}

// GetByID 根据 ID 获取配置
func (s *ModelConfigService) GetByID(ctx context.Context, id string) (*model.ModelConfig, error) {
	config, err := s.modelRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("获取模型配置失败", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return config, nil
}

// List 获取配置列表
func (s *ModelConfigService) List(ctx context.Context, deviceID uuid.UUID, page, pageSize int) ([]*model.ModelConfig, int64, error) {
	offset := (page - 1) * pageSize
	configs, total, err := s.modelRepo.List(ctx, deviceID.String(), offset, pageSize)
	if err != nil {
		logger.Error("获取模型配置列表失败", zap.Error(err))
		return nil, 0, err
	}
	return configs, total, nil
}

// Update 更新配置
func (s *ModelConfigService) Update(ctx context.Context, id string, req *dto.UpdateModelConfigRequest) (*model.ModelConfig, error) {
	config, err := s.modelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ModelName != nil {
		config.ModelName = *req.ModelName
	}
	if req.Purpose != nil {
		config.Purpose = *req.Purpose
	}
	if req.APIKey != nil {
		config.APIKey = *req.APIKey
	}
	if req.BaseURL != nil {
		config.BaseURL = *req.BaseURL
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	}

	if err := s.modelRepo.Update(ctx, config); err != nil {
		logger.Error("更新模型配置失败", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return config, nil
}

// Delete 删除配置
func (s *ModelConfigService) Delete(ctx context.Context, id string) error {
	if err := s.modelRepo.Delete(ctx, id); err != nil {
		logger.Error("删除模型配置失败", zap.String("id", id), zap.Error(err))
		return err
	}
	return nil
}

// ListProviders 获取提供商列表
func (s *ModelConfigService) ListProviders(ctx context.Context) ([]*model.ModelProvider, error) {
	providers, err := s.modelRepo.ListActiveProviders(ctx)
	if err != nil {
		logger.Error("获取提供商列表失败", zap.Error(err))
		return nil, err
	}
	return providers, nil
}

// Validate 验证模型配置
func (s *ModelConfigService) Validate(ctx context.Context, req *dto.ValidateModelConfigRequest) error {
	// 获取提供商
	provider, err := s.modelRepo.GetProviderByID(ctx, req.ProviderID)
	if err != nil {
		return errors.New("提供商不存在")
	}

	// 确定 BaseURL
	baseURL := req.BaseURL
	if baseURL == "" {
		baseURL = provider.BaseURL
	}

	// 使用 OpenAI 适配器进行验证
	adapter := llm.NewOpenAIAdapter(baseURL, "gpt-3.5-turbo")
	if err := adapter.ValidateConfig(req.APIKey, baseURL); err != nil {
		return err
	}

	// 尝试发送一个简单请求验证
	messages := []llm.ChatMessage{
		{Role: "user", Content: "Hello"},
	}
	options := llm.ChatOptions{
		Temperature: 0.7,
		MaxTokens:   10,
		APIKey:      req.APIKey,
	}

	_, err = adapter.ChatCompletion(ctx, messages, options)
	if err != nil {
		logger.Error("验证模型配置失败", zap.Error(err))
		return errors.New("API 验证失败: " + err.Error())
	}

	return nil
}

// GetByPurpose 根据用途获取配置
func (s *ModelConfigService) GetByPurpose(ctx context.Context, deviceID uuid.UUID, purpose string) (*model.ModelConfig, error) {
	config, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), purpose)
	if err != nil {
		logger.Error("获取模型配置失败",
			zap.String("device_id", deviceID.String()),
			zap.String("purpose", purpose),
			zap.Error(err),
		)
		return nil, err
	}
	return config, nil
}

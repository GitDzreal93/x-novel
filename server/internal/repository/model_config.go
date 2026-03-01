package repository

import (
	"context"
	"x-novel/internal/model"

	"gorm.io/gorm"
)

// ModelConfigRepository 模型配置仓储
type ModelConfigRepository struct {
	db *gorm.DB
}

// NewModelConfigRepository 创建模型配置仓储
func NewModelConfigRepository(db *gorm.DB) *ModelConfigRepository {
	return &ModelConfigRepository{db: db}
}

// Create 创建模型配置
func (r *ModelConfigRepository) Create(ctx context.Context, config *model.ModelConfig) error {
	return r.db.WithContext(ctx).Create(config).Error
}

// GetByID 根据 ID 获取配置
func (r *ModelConfigRepository) GetByID(ctx context.Context, id string) (*model.ModelConfig, error) {
	var config model.ModelConfig
	err := r.db.WithContext(ctx).
		Preload("Provider").
		Where("id = ?", id).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// List 获取配置列表
func (r *ModelConfigRepository) List(ctx context.Context, deviceID string, offset, limit int) ([]*model.ModelConfig, int64, error) {
	var configs []*model.ModelConfig
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ModelConfig{}).
		Preload("Provider").
		Where("device_id = ?", deviceID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&configs).Error

	return configs, total, err
}

// Update 更新配置
func (r *ModelConfigRepository) Update(ctx context.Context, config *model.ModelConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}

// Delete 删除配置
func (r *ModelConfigRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ModelConfig{}).Error
}

// GetByPurpose 根据用途获取配置（通过 model_bindings 表间接查找）
func (r *ModelConfigRepository) GetByPurpose(ctx context.Context, deviceID string, purpose string) (*model.ModelConfig, error) {
	var binding model.ModelBinding
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND purpose = ?", deviceID, purpose).
		First(&binding).Error
	if err != nil {
		// fallback: 尝试查找 general 绑定
		if purpose != "general" {
			err = r.db.WithContext(ctx).
				Where("device_id = ? AND purpose = ?", deviceID, "general").
				First(&binding).Error
		}
		if err != nil {
			return nil, err
		}
	}

	var config model.ModelConfig
	err = r.db.WithContext(ctx).
		Preload("Provider").
		Where("id = ? AND is_active = ?", binding.ModelConfigID, true).
		First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// ========== ModelBinding CRUD ==========

// ListBindings 获取设备的所有功能绑定
func (r *ModelConfigRepository) ListBindings(ctx context.Context, deviceID string) ([]*model.ModelBinding, error) {
	var bindings []*model.ModelBinding
	err := r.db.WithContext(ctx).
		Preload("ModelConfig").
		Preload("ModelConfig.Provider").
		Where("device_id = ?", deviceID).
		Order("purpose ASC").
		Find(&bindings).Error
	return bindings, err
}

// UpsertBinding 创建或更新功能绑定（同一 device + purpose 只能有一条）
func (r *ModelConfigRepository) UpsertBinding(ctx context.Context, binding *model.ModelBinding) error {
	var existing model.ModelBinding
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND purpose = ?", binding.DeviceID, binding.Purpose).
		First(&existing).Error

	if err == nil {
		return r.db.WithContext(ctx).
			Model(&existing).
			Update("model_config_id", binding.ModelConfigID).Error
	}
	return r.db.WithContext(ctx).Create(binding).Error
}

// DeleteBinding 删除功能绑定
func (r *ModelConfigRepository) DeleteBinding(ctx context.Context, deviceID string, purpose string) error {
	return r.db.WithContext(ctx).
		Where("device_id = ? AND purpose = ?", deviceID, purpose).
		Delete(&model.ModelBinding{}).Error
}

// DeleteBindingsByConfigID 删除某个配置关联的所有绑定
func (r *ModelConfigRepository) DeleteBindingsByConfigID(ctx context.Context, configID string) error {
	return r.db.WithContext(ctx).
		Where("model_config_id = ?", configID).
		Delete(&model.ModelBinding{}).Error
}

// GetProviderByID 根据 ID 获取提供商
func (r *ModelConfigRepository) GetProviderByID(ctx context.Context, id int) (*model.ModelProvider, error) {
	var provider model.ModelProvider
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// ListActiveProviders 获取活跃的提供商列表
func (r *ModelConfigRepository) ListActiveProviders(ctx context.Context) ([]*model.ModelProvider, error) {
	var providers []*model.ModelProvider
	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&providers).Error
	return providers, err
}

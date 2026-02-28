package repository

import (
	"context"
	"x-novel/internal/model"

	"gorm.io/gorm"
)

// ProjectRepository 项目仓储
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建项目仓储
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create 创建项目
func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

// GetByID 根据 ID 获取项目
func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
	var project model.Project
	err := r.db.WithContext(ctx).
		Preload("Chapters").
		Where("id = ?", id).
		First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// List 获取项目列表
func (r *ProjectRepository) List(ctx context.Context, deviceID string, offset, limit int) ([]*model.Project, int64, error) {
	var projects []*model.Project
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Project{}).Where("device_id = ?", deviceID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	err := query.Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&projects).Error

	return projects, total, err
}

// Update 更新项目
func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

// Delete 删除项目
func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Project{}).Error
}

// UpdateArchitecture 更新架构数据
func (r *ProjectRepository) UpdateArchitecture(ctx context.Context, id string, architecture map[string]string) error {
	return r.db.WithContext(ctx).Model(&model.Project{}).
		Where("id = ?", id).
		Updates(architecture).Error
}

// UpdateBlueprint 更新大纲数据
func (r *ProjectRepository) UpdateBlueprint(ctx context.Context, id string, blueprint string) error {
	return r.db.WithContext(ctx).Model(&model.Project{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"chapter_blueprint":   blueprint,
			"blueprint_generated": true,
		}).Error
}

// UpdateGlobalSummary 更新全局摘要
func (r *ProjectRepository) UpdateGlobalSummary(ctx context.Context, id string, summary string) error {
	return r.db.WithContext(ctx).Model(&model.Project{}).
		Where("id = ?", id).
		Update("global_summary", summary).Error
}

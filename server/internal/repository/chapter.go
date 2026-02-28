package repository

import (
	"context"
	"x-novel/internal/model"

	"gorm.io/gorm"
)

// ChapterRepository 章节仓储
type ChapterRepository struct {
	db *gorm.DB
}

// NewChapterRepository 创建章节仓储
func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

// Create 创建章节
func (r *ChapterRepository) Create(ctx context.Context, chapter *model.Chapter) error {
	return r.db.WithContext(ctx).Create(chapter).Error
}

// GetByID 根据 ID 获取章节
func (r *ChapterRepository) GetByID(ctx context.Context, id string) (*model.Chapter, error) {
	var chapter model.Chapter
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

// GetByProjectAndNumber 根据项目 ID 和章节号获取章节
func (r *ChapterRepository) GetByProjectAndNumber(ctx context.Context, projectID string, chapterNumber int) (*model.Chapter, error) {
	var chapter model.Chapter
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND chapter_number = ?", projectID, chapterNumber).
		First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

// List 获取章节列表
func (r *ChapterRepository) List(ctx context.Context, projectID string, offset, limit int) ([]*model.Chapter, int64, error) {
	var chapters []*model.Chapter
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Chapter{}).Where("project_id = ?", projectID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取列表
	err := query.Order("chapter_number ASC").
		Offset(offset).
		Limit(limit).
		Find(&chapters).Error

	return chapters, total, err
}

// Update 更新章节
func (r *ChapterRepository) Update(ctx context.Context, chapter *model.Chapter) error {
	return r.db.WithContext(ctx).Save(chapter).Error
}

// Delete 删除章节
func (r *ChapterRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Chapter{}).Error
}

// UpdateContent 更新章节内容
func (r *ChapterRepository) UpdateContent(ctx context.Context, id string, content string, wordCount int) error {
	return r.db.WithContext(ctx).Model(&model.Chapter{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"content":    content,
			"word_count": wordCount,
		}).Error
}

// UpdateStatus 更新章节状态
func (r *ChapterRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.db.WithContext(ctx).Model(&model.Chapter{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// SetFinalized 设置章节为定稿状态
func (r *ChapterRepository) SetFinalized(ctx context.Context, id string, finalized bool) error {
	return r.db.WithContext(ctx).Model(&model.Chapter{}).
		Where("id = ?", id).
		Update("is_finalized", finalized).Error
}

// ListCompleted 获取已完成的章节列表（用于生成前文摘要）
func (r *ChapterRepository) ListCompleted(ctx context.Context, projectID string, beforeChapterNumber int) ([]*model.Chapter, error) {
	var chapters []*model.Chapter
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND chapter_number < ? AND status = ?", projectID, beforeChapterNumber, "completed").
		Order("chapter_number ASC").
		Find(&chapters).Error
	return chapters, err
}

// CountCompleted 统计已完成章节数
func (r *ChapterRepository) CountCompleted(ctx context.Context, projectID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Chapter{}).
		Where("project_id = ? AND status = ?", projectID, "completed").
		Count(&count).Error
	return count, err
}

// GetTotalWords 获取项目总字数
func (r *ChapterRepository) GetTotalWords(ctx context.Context, projectID string) (int64, error) {
	var result struct {
		Total int64
	}
	err := r.db.WithContext(ctx).Model(&model.Chapter{}).
		Select("COALESCE(SUM(word_count), 0) as total").
		Where("project_id = ?", projectID).
		Scan(&result).Error
	return result.Total, err
}

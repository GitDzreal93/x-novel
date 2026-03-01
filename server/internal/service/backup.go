package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// BackupProject 备份项目数据
type BackupProject struct {
	model.Project
	Chapters []model.Chapter `json:"chapters"`
}

// BackupData 完整备份数据
type BackupData struct {
	Version       string              `json:"version"`
	ExportedAt    time.Time           `json:"exported_at"`
	DeviceID      string              `json:"device_id"`
	Projects      []BackupProject     `json:"projects"`
	Conversations []BackupConversation `json:"conversations"`
}

// BackupConversation 备份对话数据
type BackupConversation struct {
	model.Conversation
	Messages []model.Message `json:"messages"`
}

type BackupService struct {
	db          *gorm.DB
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
	chatRepo    *repository.ChatRepository
}

func NewBackupService(
	db *gorm.DB,
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	chatRepo *repository.ChatRepository,
) *BackupService {
	return &BackupService{
		db:          db,
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
		chatRepo:    chatRepo,
	}
}

// Export 导出设备下所有数据
func (s *BackupService) Export(ctx context.Context, deviceID uuid.UUID) (*BackupData, error) {
	backup := &BackupData{
		Version:    "1.0.0",
		ExportedAt: time.Now(),
		DeviceID:   deviceID.String(),
	}

	projects, _, err := s.projectRepo.List(ctx, deviceID.String(), 0, 10000)
	if err != nil {
		return nil, fmt.Errorf("获取项目列表失败: %w", err)
	}

	for _, p := range projects {
		bp := BackupProject{Project: *p}

		chapters, err := s.chapterRepo.ListByProject(ctx, p.ID.String())
		if err != nil {
			logger.Warn("获取章节列表失败", zap.String("project_id", p.ID.String()), zap.Error(err))
			continue
		}
		for _, ch := range chapters {
			bp.Chapters = append(bp.Chapters, *ch)
		}
		backup.Projects = append(backup.Projects, bp)
	}

	conversations, _, err := s.chatRepo.ListConversations(ctx, deviceID, 0, 10000)
	if err != nil {
		logger.Warn("获取对话列表失败", zap.Error(err))
	} else {
		for _, conv := range conversations {
			bc := BackupConversation{Conversation: *conv}
			messages, err := s.chatRepo.ListMessages(ctx, conv.ID.String())
			if err != nil {
				logger.Warn("获取消息失败", zap.String("conversation_id", conv.ID.String()), zap.Error(err))
				continue
			}
			for _, msg := range messages {
				bc.Messages = append(bc.Messages, *msg)
			}
			backup.Conversations = append(backup.Conversations, bc)
		}
	}

	return backup, nil
}

// Import 导入备份数据到当前设备
func (s *BackupService) Import(ctx context.Context, deviceID uuid.UUID, data []byte) (*ImportResult, error) {
	var backup BackupData
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("解析备份数据失败: %w", err)
	}

	result := &ImportResult{}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, bp := range backup.Projects {
			newProjectID := uuid.New()
			oldProjectID := bp.Project.ID

			bp.Project.ID = newProjectID
			bp.Project.DeviceID = deviceID
			bp.Project.CreatedAt = time.Now()
			bp.Project.UpdatedAt = time.Now()

			if err := tx.Create(&bp.Project).Error; err != nil {
				logger.Warn("导入项目失败", zap.String("title", bp.Project.Title), zap.Error(err))
				result.FailedProjects++
				continue
			}
			result.ImportedProjects++

			logger.Info("导入项目成功",
				zap.String("old_id", oldProjectID.String()),
				zap.String("new_id", newProjectID.String()),
				zap.String("title", bp.Project.Title),
			)

			for _, ch := range bp.Chapters {
				ch.ID = uuid.New()
				ch.ProjectID = newProjectID
				ch.CreatedAt = time.Now()
				ch.UpdatedAt = time.Now()

				if err := tx.Create(&ch).Error; err != nil {
					logger.Warn("导入章节失败",
						zap.Int("chapter_number", ch.ChapterNumber),
						zap.Error(err),
					)
					result.FailedChapters++
					continue
				}
				result.ImportedChapters++
			}
		}

		for _, bc := range backup.Conversations {
			newConvID := uuid.New()

			bc.Conversation.ID = newConvID
			bc.Conversation.DeviceID = deviceID
			bc.Conversation.CreatedAt = time.Now()
			bc.Conversation.UpdatedAt = time.Now()

			if err := tx.Create(&bc.Conversation).Error; err != nil {
				logger.Warn("导入对话失败", zap.String("title", bc.Conversation.Title), zap.Error(err))
				result.FailedConversations++
				continue
			}
			result.ImportedConversations++

			for _, msg := range bc.Messages {
				msg.ID = uuid.New()
				msg.ConversationID = newConvID
				msg.CreatedAt = time.Now()

				if err := tx.Create(&msg).Error; err != nil {
					result.FailedMessages++
					continue
				}
				result.ImportedMessages++
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("导入失败: %w", err)
	}

	return result, nil
}

// ImportResult 导入结果统计
type ImportResult struct {
	ImportedProjects      int `json:"imported_projects"`
	ImportedChapters      int `json:"imported_chapters"`
	ImportedConversations int `json:"imported_conversations"`
	ImportedMessages      int `json:"imported_messages"`
	FailedProjects        int `json:"failed_projects"`
	FailedChapters        int `json:"failed_chapters"`
	FailedConversations   int `json:"failed_conversations"`
	FailedMessages        int `json:"failed_messages"`
}

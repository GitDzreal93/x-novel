package repository

import (
	"context"

	"x-novel/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// ========== Conversation ==========

func (r *ChatRepository) CreateConversation(ctx context.Context, conv *model.Conversation) error {
	return r.db.WithContext(ctx).Create(conv).Error
}

func (r *ChatRepository) GetConversation(ctx context.Context, id string) (*model.Conversation, error) {
	var conv model.Conversation
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&conv).Error
	return &conv, err
}

func (r *ChatRepository) ListConversations(ctx context.Context, deviceID uuid.UUID, offset, limit int) ([]*model.Conversation, int64, error) {
	var conversations []*model.Conversation
	var total int64

	query := r.db.WithContext(ctx).Where("device_id = ?", deviceID)
	query.Model(&model.Conversation{}).Count(&total)

	err := query.Order("updated_at DESC").Offset(offset).Limit(limit).Find(&conversations).Error
	return conversations, total, err
}

func (r *ChatRepository) UpdateConversationTitle(ctx context.Context, id string, title string) error {
	return r.db.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", id).Update("title", title).Error
}

func (r *ChatRepository) DeleteConversation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("conversation_id = ?", id).Delete(&model.Message{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", id).Delete(&model.Conversation{}).Error
	})
}

func (r *ChatRepository) TouchConversation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&model.Conversation{}).Where("id = ?", id).Update("updated_at", gorm.Expr("NOW()")).Error
}

// ========== Message ==========

func (r *ChatRepository) CreateMessage(ctx context.Context, msg *model.Message) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *ChatRepository) ListMessages(ctx context.Context, conversationID string) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).
		Order("created_at ASC").Find(&messages).Error
	return messages, err
}

// RecentMessages 获取最近 N 条消息（用于构建上下文）
func (r *ChatRepository) RecentMessages(ctx context.Context, conversationID string, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).
		Order("created_at DESC").Limit(limit).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	// 反转为时间正序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

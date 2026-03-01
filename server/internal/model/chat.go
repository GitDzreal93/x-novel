package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Conversation 对话会话
type Conversation struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DeviceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Mode      string    `gorm:"size:30;not null;default:general" json:"mode"` // creative, building, character, general
	ProjectID *uuid.UUID `gorm:"type:uuid;index" json:"project_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Messages []Message `gorm:"-" json:"messages,omitempty"`
}

func (Conversation) TableName() string {
	return "conversations"
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// Message 对话消息
type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index" json:"conversation_id"`
	Role           string    `gorm:"size:20;not null" json:"role"` // user, assistant
	Content        string    `gorm:"type:text;not null" json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

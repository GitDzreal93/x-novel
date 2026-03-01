package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelProvider 模型提供商
type ModelProvider struct {
	ID          int       `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"size:50;uniqueIndex;not null" json:"name"` // openai, anthropic, google, etc.
	DisplayName string    `gorm:"size:100;not null" json:"display_name"`
	BaseURL     string    `gorm:"size:500" json:"base_url"`
	AuthType    string    `gorm:"size:20;default:api_key" json:"auth_type"` // api_key, oauth, custom
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	Configs []ModelConfig `gorm:"-" json:"configs,omitempty"`
}

func (ModelProvider) TableName() string {
	return "model_providers"
}

// ModelConfig 模型配置（只管连接信息，不绑定用途）
type ModelConfig struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DeviceID   uuid.UUID `gorm:"type:uuid;not null;index:idx_device_model" json:"device_id"`
	ProviderID int       `gorm:"not null;index:idx_device_model" json:"provider_id"`
	ModelName  string    `gorm:"size:100;not null;index:idx_device_model" json:"model_name"`
	Purpose    string    `gorm:"size:50;index:idx_purpose" json:"purpose,omitempty"` // 已废弃，保留字段避免迁移问题
	APIKey     string    `gorm:"type:text;not null" json:"api_key"`
	BaseURL    string    `gorm:"size:500" json:"base_url,omitempty"`
	IsActive   bool      `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联
	Device   *Device        `gorm:"-" json:"device,omitempty"`
	Provider *ModelProvider `gorm:"-" json:"provider,omitempty"`
}

func (ModelConfig) TableName() string {
	return "model_configs"
}

// BeforeCreate GORM hook
func (mc *ModelConfig) BeforeCreate(tx *gorm.DB) error {
	if mc.ID == uuid.Nil {
		mc.ID = uuid.New()
	}
	return nil
}

// ModelBinding 功能绑定：哪个功能使用哪个模型配置
type ModelBinding struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DeviceID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_device_purpose" json:"device_id"`
	Purpose       string    `gorm:"size:50;not null;uniqueIndex:idx_device_purpose" json:"purpose"` // architecture, chapter, writing, review, general
	ModelConfigID uuid.UUID `gorm:"type:uuid;not null" json:"model_config_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// 关联
	ModelConfig *ModelConfig `gorm:"foreignKey:ModelConfigID" json:"model_config,omitempty"`
}

func (ModelBinding) TableName() string {
	return "model_bindings"
}

func (mb *ModelBinding) BeforeCreate(tx *gorm.DB) error {
	if mb.ID == uuid.Nil {
		mb.ID = uuid.New()
	}
	return nil
}

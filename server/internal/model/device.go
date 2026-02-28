package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Device 设备模型
type Device struct {
	ID        uuid.UUID              `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DeviceID  string                 `gorm:"uniqueIndex;not null" json:"device_id"`
	DeviceInfo map[string]interface{} `gorm:"type:jsonb" json:"device_info,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	LastSeen  time.Time              `json:"last_seen"`
}

func (Device) TableName() string {
	return "devices"
}

// BeforeCreate GORM hook
func (d *Device) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

// DeviceSetting 设备设置模型
type DeviceSetting struct {
	ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DeviceID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"device_id"`
	Theme            string    `gorm:"default:light" json:"theme"`               // light, dark
	Language         string    `gorm:"default:zh-CN" json:"language"`
	AutoSaveEnabled  bool      `gorm:"default:true" json:"auto_save_enabled"`
	AutoSaveInterval int       `gorm:"default:30000" json:"auto_save_interval"` // ms
}

func (DeviceSetting) TableName() string {
	return "device_settings"
}

// BeforeCreate GORM hook
func (ds *DeviceSetting) BeforeCreate(tx *gorm.DB) error {
	if ds.ID == uuid.Nil {
		ds.ID = uuid.New()
	}
	return nil
}

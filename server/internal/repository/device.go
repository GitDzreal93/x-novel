package repository

import (
	"context"
	"errors"
	"time"

	"x-novel/internal/model"

	"gorm.io/gorm"
)

// DeviceRepository 设备仓储
type DeviceRepository struct {
	db *gorm.DB
}

// NewDeviceRepository 创建设备仓储
func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// GetOrCreateByDeviceID 根据设备 ID 获取或创建设备
func (r *DeviceRepository) GetOrCreateByDeviceID(ctx context.Context, deviceID string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&device).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新设备
			device = model.Device{
				DeviceID:  deviceID,
				LastSeen:  time.Now(),
			}
			if err := r.db.WithContext(ctx).Create(&device).Error; err != nil {
				return nil, err
			}

			// 创建默认设置
			setting := &model.DeviceSetting{
				DeviceID: device.ID,
			}
			if err := r.db.WithContext(ctx).Create(setting).Error; err != nil {
				return &device, err // 忽略设置创建错误
			}

			return &device, nil
		}
		return nil, err
	}

	// 更新最后访问时间
	device.LastSeen = time.Now()
	if err := r.db.WithContext(ctx).Save(&device).Error; err != nil {
		return nil, err
	}

	return &device, nil
}

// GetByID 根据 ID 获取设备
func (r *DeviceRepository) GetByID(ctx context.Context, id string) (*model.Device, error) {
	var device model.Device
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

// Update 更新设备
func (r *DeviceRepository) Update(ctx context.Context, device *model.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

// GetSetting 获取设备设置
func (r *DeviceRepository) GetSetting(ctx context.Context, deviceID string) (*model.DeviceSetting, error) {
	var setting model.DeviceSetting
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// UpdateSetting 更新设备设置
func (r *DeviceRepository) UpdateSetting(ctx context.Context, setting *model.DeviceSetting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

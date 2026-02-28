package service

import (
	"context"

	"x-novel/internal/dto"
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"go.uber.org/zap"
)

// DeviceService 设备服务
type DeviceService struct {
	deviceRepo *repository.DeviceRepository
}

// NewDeviceService 创建设备服务
func NewDeviceService(deviceRepo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{
		deviceRepo: deviceRepo,
	}
}

// GetInfo 获取设备信息
func (s *DeviceService) GetInfo(ctx context.Context, deviceUUID string) (*model.Device, error) {
	device, err := s.deviceRepo.GetByID(ctx, deviceUUID)
	if err != nil {
		logger.Error("获取设备信息失败",
			zap.String("device_uuid", deviceUUID),
			zap.Error(err),
		)
		return nil, err
	}
	return device, nil
}

// GetSettings 获取设备设置
func (s *DeviceService) GetSettings(ctx context.Context, deviceUUID string) (*model.DeviceSetting, error) {
	setting, err := s.deviceRepo.GetSetting(ctx, deviceUUID)
	if err != nil {
		logger.Error("获取设备设置失败",
			zap.String("device_uuid", deviceUUID),
			zap.Error(err),
		)
		return nil, err
	}
	return setting, nil
}

// UpdateSettings 更新设备设置
func (s *DeviceService) UpdateSettings(ctx context.Context, deviceUUID string, req *dto.UpdateDeviceSettingsRequest) (*model.DeviceSetting, error) {
	setting, err := s.deviceRepo.GetSetting(ctx, deviceUUID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Theme != nil {
		setting.Theme = *req.Theme
	}
	if req.Language != nil {
		setting.Language = *req.Language
	}
	if req.AutoSaveEnabled != nil {
		setting.AutoSaveEnabled = *req.AutoSaveEnabled
	}
	if req.AutoSaveInterval != nil {
		setting.AutoSaveInterval = *req.AutoSaveInterval
	}

	if err := s.deviceRepo.UpdateSetting(ctx, setting); err != nil {
		logger.Error("更新设备设置失败",
			zap.String("device_uuid", deviceUUID),
			zap.Error(err),
		)
		return nil, err
	}

	return setting, nil
}

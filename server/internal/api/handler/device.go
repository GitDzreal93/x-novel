package handler

import (
	"net/http"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

// DeviceHandler 设备处理器
type DeviceHandler struct {
	deviceService *service.DeviceService
}

// NewDeviceHandler 创建设备处理器
func NewDeviceHandler(deviceService *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// GetInfo 获取设备信息
// @Summary 获取设备信息
// @Description 获取当前设备的信息
// @Tags device
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response{data=dto.DeviceResponse}
// @Router /api/v1/device/info [get]
func (h *DeviceHandler) GetInfo(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	device, err := h.deviceService.GetInfo(c.Request.Context(), deviceUUID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取设备信息失败",
		})
		return
	}

	response := &dto.DeviceResponse{
		ID:        device.ID,
		DeviceID:  device.DeviceID,
		DeviceInfo: device.DeviceInfo,
		CreatedAt: device.CreatedAt,
		LastSeen:  device.LastSeen,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// GetSettings 获取设备设置
// @Summary 获取设备设置
// @Description 获取当前设备的设置
// @Tags device
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response{data=dto.DeviceSettingResponse}
// @Router /api/v1/device/settings [get]
func (h *DeviceHandler) GetSettings(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	setting, err := h.deviceService.GetSettings(c.Request.Context(), deviceUUID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取设备设置失败",
		})
		return
	}

	response := &dto.DeviceSettingResponse{
		ID:               setting.ID,
		DeviceID:         setting.DeviceID,
		Theme:            setting.Theme,
		Language:         setting.Language,
		AutoSaveEnabled:  setting.AutoSaveEnabled,
		AutoSaveInterval: setting.AutoSaveInterval,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// UpdateSettings 更新设备设置
// @Summary 更新设备设置
// @Description 更新当前设备的设置
// @Tags device
// @Accept json
// @Produce json
// @Param request body dto.UpdateDeviceSettingsRequest true "更新请求"
// @Success 200 {object} dto.Response{data=dto.DeviceSettingResponse}
// @Router /api/v1/device/settings [put]
func (h *DeviceHandler) UpdateSettings(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	var req dto.UpdateDeviceSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	setting, err := h.deviceService.UpdateSettings(c.Request.Context(), deviceUUID.String(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "更新设备设置失败",
		})
		return
	}

	response := &dto.DeviceSettingResponse{
		ID:               setting.ID,
		DeviceID:         setting.DeviceID,
		Theme:            setting.Theme,
		Language:         setting.Language,
		AutoSaveEnabled:  setting.AutoSaveEnabled,
		AutoSaveInterval: setting.AutoSaveInterval,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

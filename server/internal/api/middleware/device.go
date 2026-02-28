package middleware

import (
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/device"
	"x-novel/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Device 设备识别中间件
func Device(deviceRepo *repository.DeviceRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取设备 ID
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			// 如果没有设备 ID，生成一个新的
			deviceID = device.GenerateDeviceID()
		}

		// 查找或创建设备
		dev, err := deviceRepo.GetOrCreateByDeviceID(c.Request.Context(), deviceID)
		if err != nil {
			logger.Error("获取设备信息失败",
				zap.String("device_id", deviceID),
				zap.Error(err),
			)
			// 继续处理，但不中断请求
		}

		// 将设备 ID 和设备信息存入上下文
		c.Set("device_id", deviceID)
		if dev != nil {
			c.Set("device", dev)
		}

		// 在响应头中返回设备 ID
		c.Header("X-Device-ID", deviceID)

		c.Next()
	}
}

// GetDevice 从上下文获取设备
func GetDevice(c *gin.Context) (*model.Device, bool) {
	if dev, exists := c.Get("device"); exists {
		if d, ok := dev.(*model.Device); ok {
			return d, true
		}
	}
	return nil, false
}

// GetDeviceID 从上下文获取设备 ID
func GetDeviceID(c *gin.Context) string {
	if deviceID, exists := c.Get("device_id"); exists {
		if id, ok := deviceID.(string); ok {
			return id
		}
	}
	return ""
}

// GetDeviceUUID 从上下文获取设备 UUID
func GetDeviceUUID(c *gin.Context) (uuid.UUID, error) {
	dev, exists := GetDevice(c)
	if !exists {
		return uuid.Nil, ErrDeviceNotFound
	}
	return dev.ID, nil
}

// Errors
var (
	ErrDeviceNotFound = &MiddlewareError{Message: "设备未找到"}
)

// MiddlewareError 中间件错误
type MiddlewareError struct {
	Message string
}

func (e *MiddlewareError) Error() string {
	return e.Message
}

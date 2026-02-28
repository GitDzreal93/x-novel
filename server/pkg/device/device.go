package device

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// GenerateDeviceID 生成设备 ID
func GenerateDeviceID() string {
	// 生成基于时间戳和随机数的设备 ID
	timestamp := time.Now().UnixNano()
	random := uuid.New().String()
	hash := sha256.Sum256([]byte(strconv.FormatInt(timestamp, 10) + random))
	return hex.EncodeToString(hash[:])[:32]
}

// ParseDeviceID 从请求头或 Cookie 解析设备 ID
func ParseDeviceID(header string, cookie string) string {
	if header != "" {
		return header
	}
	if cookie != "" {
		return cookie
	}
	return ""
}

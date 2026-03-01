package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	backupService *service.BackupService
}

func NewBackupHandler(backupService *service.BackupService) *BackupHandler {
	return &BackupHandler{backupService: backupService}
}

// Export 导出所有数据为 JSON
func (h *BackupHandler) Export(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	backup, err := h.backupService.Export(c.Request.Context(), deviceUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "导出失败: " + err.Error()})
		return
	}

	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "序列化失败"})
		return
	}

	filename := fmt.Sprintf("x-novel-backup-%s.json", time.Now().Format("20060102-150405"))

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", data)
}

// ExportPreview 预览导出数据的概要
func (h *BackupHandler) ExportPreview(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	backup, err := h.backupService.Export(c.Request.Context(), deviceUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "获取数据失败: " + err.Error()})
		return
	}

	totalChapters := 0
	totalWords := 0
	for _, p := range backup.Projects {
		totalChapters += len(p.Chapters)
		for _, ch := range p.Chapters {
			totalWords += ch.WordCount
		}
	}

	totalMessages := 0
	for _, conv := range backup.Conversations {
		totalMessages += len(conv.Messages)
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data: map[string]interface{}{
			"projects":      len(backup.Projects),
			"chapters":      totalChapters,
			"total_words":   totalWords,
			"conversations": len(backup.Conversations),
			"messages":      totalMessages,
		},
	})
}

// Import 导入 JSON 备份数据
func (h *BackupHandler) Import(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "请上传备份文件"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "读取文件失败"})
		return
	}

	const maxSize = 100 * 1024 * 1024 // 100MB
	if len(data) > maxSize {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "文件过大，最大支持 100MB"})
		return
	}

	result, err := h.backupService.Import(c.Request.Context(), deviceUUID, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "导入失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "导入完成",
		Data:    result,
	})
}

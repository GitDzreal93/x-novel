package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

type WritingAssistantHandler struct {
	writingService *service.WritingAssistantService
}

func NewWritingAssistantHandler(writingService *service.WritingAssistantService) *WritingAssistantHandler {
	return &WritingAssistantHandler{writingService: writingService}
}

// Assist 写作助手统一入口
func (h *WritingAssistantHandler) Assist(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	var req dto.WritingAssistantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "请求参数错误: " + err.Error()})
		return
	}

	if req.Stream {
		h.handleStream(c, deviceUUID, &req)
		return
	}

	var result string
	switch req.Action {
	case "polish":
		result, err = h.writingService.Polish(c.Request.Context(), deviceUUID, req.Content, req.Style)
	case "continue":
		targetWords := req.TargetWords
		if targetWords <= 0 {
			targetWords = 500
		}
		result, err = h.writingService.Continue(c.Request.Context(), deviceUUID, req.ProjectID, req.Content, targetWords)
	case "suggestion":
		aspect := req.Aspect
		if aspect == "" {
			aspect = "plot"
		}
		result, err = h.writingService.Suggest(c.Request.Context(), deviceUUID, req.ProjectID, req.Content, aspect)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "处理失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data:    map[string]string{"result": result},
	})
}

func (h *WritingAssistantHandler) handleStream(c *gin.Context, deviceUUID interface{ String() string }, req *dto.WritingAssistantRequest) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "不支持流式响应"})
		return
	}

	deviceID, _ := middleware.GetDeviceUUID(c)

	callback := func(chunk string) error {
		data, _ := json.Marshal(map[string]string{"content": chunk})
		_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		if err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	var result string
	var err error

	switch req.Action {
	case "polish":
		result, err = h.writingService.PolishStream(c.Request.Context(), deviceID, req.Content, req.Style, callback)
	case "continue":
		targetWords := req.TargetWords
		if targetWords <= 0 {
			targetWords = 500
		}
		result, err = h.writingService.ContinueStream(c.Request.Context(), deviceID, req.ProjectID, req.Content, targetWords, callback)
	case "suggestion":
		aspect := req.Aspect
		if aspect == "" {
			aspect = "plot"
		}
		result, err = h.writingService.SuggestStream(c.Request.Context(), deviceID, req.ProjectID, req.Content, aspect, callback)
	}

	if err != nil {
		errData, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(c.Writer, "data: %s\n\n", errData)
		flusher.Flush()
		return
	}

	doneData, _ := json.Marshal(map[string]interface{}{"done": true, "result": result})
	fmt.Fprintf(c.Writer, "data: %s\n\n", doneData)
	flusher.Flush()
}

package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// CreateConversation 创建对话
func (h *ChatHandler) CreateConversation(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "请求参数错误"})
		return
	}

	var projectID *uuid.UUID
	if req.ProjectID != "" {
		id, err := uuid.Parse(req.ProjectID)
		if err == nil {
			projectID = &id
		}
	}

	conv, err := h.chatService.CreateConversation(c.Request.Context(), deviceUUID, req.Title, req.Mode, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "创建对话失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: dto.ConversationFromModel(conv)})
}

// ListConversations 获取对话列表
func (h *ChatHandler) ListConversations(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	conversations, total, err := h.chatService.ListConversations(c.Request.Context(), deviceUUID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "获取对话列表失败"})
		return
	}

	items := make([]dto.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		items[i] = *dto.ConversationFromModel(conv)
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data:    dto.ConversationListResponse{Conversations: items, Total: total},
	})
}

// GetConversation 获取对话详情（含消息）
func (h *ChatHandler) GetConversation(c *gin.Context) {
	id := c.Param("id")
	conv, err := h.chatService.GetConversation(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Code: 404, Message: "对话不存在"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: dto.ConversationFromModel(conv)})
}

// UpdateConversation 更新对话标题
func (h *ChatHandler) UpdateConversation(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "请求参数错误"})
		return
	}

	if err := h.chatService.UpdateConversationTitle(c.Request.Context(), id, req.Title); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "更新对话失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success"})
}

// DeleteConversation 删除对话
func (h *ChatHandler) DeleteConversation(c *gin.Context) {
	id := c.Param("id")
	if err := h.chatService.DeleteConversation(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "删除对话失败"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success"})
}

// SendMessage 发送消息（支持 SSE 流式 / 普通 JSON）
func (h *ChatHandler) SendMessage(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	conversationID := c.Param("id")

	var req dto.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "请求参数错误"})
		return
	}

	if req.Stream {
		h.handleStreamMessage(c, deviceUUID, conversationID, req.Content)
		return
	}

	userMsg, assistantMsg, err := h.chatService.SendMessage(c.Request.Context(), deviceUUID, conversationID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: fmt.Sprintf("发送消息失败: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data: dto.SendMessageResponse{
			UserMessage:      *dto.MessageFromModel(userMsg),
			AssistantMessage: *dto.MessageFromModel(assistantMsg),
		},
	})
}

func (h *ChatHandler) handleStreamMessage(c *gin.Context, deviceUUID uuid.UUID, conversationID, content string) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "不支持流式响应"})
		return
	}

	callback := func(chunk string) error {
		data, _ := json.Marshal(map[string]string{"content": chunk})
		_, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		if err != nil {
			return err
		}
		flusher.Flush()
		return nil
	}

	userMsg, assistantMsg, err := h.chatService.SendMessageStream(c.Request.Context(), deviceUUID, conversationID, content, callback)
	if err != nil {
		errData, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(c.Writer, "data: %s\n\n", errData)
		flusher.Flush()
		return
	}

	// 发送完成事件，包含完整的消息数据
	doneData, _ := json.Marshal(map[string]interface{}{
		"done":              true,
		"user_message":      dto.MessageFromModel(userMsg),
		"assistant_message": dto.MessageFromModel(assistantMsg),
	})
	fmt.Fprintf(c.Writer, "data: %s\n\n", doneData)
	flusher.Flush()
}

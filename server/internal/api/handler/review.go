package handler

import (
	"net/http"
	"strconv"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// DetectErrors 检测章节内容中的错误
func (h *ReviewHandler) DetectErrors(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	var req dto.DetectErrorsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "请求参数错误: " + err.Error()})
		return
	}

	result, err := h.reviewService.DetectErrors(c.Request.Context(), deviceUUID, req.Content, req.Types)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "检测失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data:    result,
	})
}

// ReviewChapter 审阅单个章节
func (h *ReviewHandler) ReviewChapter(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	projectID := c.Param("id")
	chapterNumber, err := strconv.Atoi(c.Param("chapterNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "章节编号无效"})
		return
	}

	result, err := h.reviewService.ReviewChapter(c.Request.Context(), deviceUUID, projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "审阅失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data:    result,
	})
}

// ReviewProject 审阅整个项目
func (h *ReviewHandler) ReviewProject(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	projectID := c.Param("id")

	result, err := h.reviewService.ReviewProject(c.Request.Context(), deviceUUID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "审阅失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data:    result,
	})
}

// MarketPredict 市场预测
func (h *ReviewHandler) MarketPredict(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	projectID := c.Param("id")

	result, err := h.reviewService.MarketPredict(c.Request.Context(), deviceUUID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: "市场预测失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    200,
		Message: "success",
		Data:    result,
	})
}

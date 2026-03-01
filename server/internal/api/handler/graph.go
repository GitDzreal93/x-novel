package handler

import (
	"net/http"
	"strconv"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

type GraphHandler struct {
	graphService *service.GraphService
}

func NewGraphHandler(graphService *service.GraphService) *GraphHandler {
	return &GraphHandler{graphService: graphService}
}

// GenerateGraph 从项目架构生成关系图谱
func (h *GraphHandler) GenerateGraph(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	projectID := c.Param("id")

	graphData, err := h.graphService.GenerateGraph(c.Request.Context(), deviceUUID, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: graphData})
}

// GetGraph 获取项目关系图谱
func (h *GraphHandler) GetGraph(c *gin.Context) {
	projectID := c.Param("id")

	graphData, err := h.graphService.GetGraph(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: graphData})
}

// UpdateFromChapter 从章节更新图谱
func (h *GraphHandler) UpdateFromChapter(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Code: 401, Message: "未授权"})
		return
	}

	projectID := c.Param("id")
	chapterNumber, err := strconv.Atoi(c.Param("chapterNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "无效的章节号"})
		return
	}

	graphData, err := h.graphService.UpdateGraphFromChapter(c.Request.Context(), deviceUUID, projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: graphData})
}

// GetChapterSnapshot 获取章节图谱快照
func (h *GraphHandler) GetChapterSnapshot(c *gin.Context) {
	projectID := c.Param("id")
	chapterNumber, err := strconv.Atoi(c.Param("chapterNumber"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: 400, Message: "无效的章节号"})
		return
	}

	graphData, err := h.graphService.GetChapterSnapshot(c.Request.Context(), projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: 500, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Code: 200, Message: "success", Data: graphData})
}

package handler

import (
	"net/http"
	"strconv"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

// ChapterHandler 章节处理器
type ChapterHandler struct {
	chapterService *service.ChapterService
	projectService *service.ProjectService
}

// NewChapterHandler 创建章节处理器
func NewChapterHandler(
	chapterService *service.ChapterService,
	projectService *service.ProjectService,
) *ChapterHandler {
	return &ChapterHandler{
		chapterService: chapterService,
		projectService: projectService,
	}
}

// Create 创建章节
// @Summary 创建章节
// @Description 为项目创建新章节
// @Tags chapter
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body dto.CreateChapterRequest true "创建请求"
// @Success 200 {object} dto.Response{data=dto.ChapterResponse}
// @Router /api/v1/projects/{id}/chapters [post]
func (h *ChapterHandler) Create(c *gin.Context) {
	projectID := c.Param("id")

	var req dto.CreateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	chapter, err := h.chapterService.Create(c.Request.Context(), projectID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.ChapterFromModel(chapter)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// List 获取章节列表
// @Summary 获取章节列表
// @Description 获取项目的章节列表
// @Tags chapter
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} dto.Response{data=dto.ChapterListResponse}
// @Router /api/v1/projects/{id}/chapters [get]
func (h *ChapterHandler) List(c *gin.Context) {
	projectID := c.Param("id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	chapters, total, err := h.chapterService.List(c.Request.Context(), projectID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取章节列表失败",
		})
		return
	}

	chapterResponses := make([]dto.ChapterResponse, 0, len(chapters))
	for _, chapter := range chapters {
		chapterResponses = append(chapterResponses, *dto.ChapterFromModel(chapter))
	}

	response := &dto.ChapterListResponse{
		Chapters: chapterResponses,
		Total:    int(total),
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// GetByNumber 获取章节详情
// @Summary 获取章节详情
// @Description 根据章节号获取章节详情
// @Tags chapter
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterNumber path int true "章节号"
// @Success 200 {object} dto.Response{data=dto.ChapterResponse}
// @Router /api/v1/projects/{id}/chapters/{chapterNumber} [get]
func (h *ChapterHandler) GetByNumber(c *gin.Context) {
	projectID := c.Param("id")
	chapterNumber, _ := strconv.Atoi(c.Param("chapterNumber"))

	chapter, err := h.chapterService.GetByProjectAndNumber(c.Request.Context(), projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "章节不存在",
		})
		return
	}

	response := dto.ChapterFromModel(chapter)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// Update 更新章节
// @Summary 更新章节
// @Description 更新章节信息
// @Tags chapter
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterNumber path int true "章节号"
// @Param request body dto.UpdateChapterRequest true "更新请求"
// @Success 200 {object} dto.Response{data=dto.ChapterResponse}
// @Router /api/v1/projects/{id}/chapters/{chapterNumber} [put]
func (h *ChapterHandler) Update(c *gin.Context) {
	projectID := c.Param("id")
	chapterNumber, _ := strconv.Atoi(c.Param("chapterNumber"))

	// 先获取章节
	chapter, err := h.chapterService.GetByProjectAndNumber(c.Request.Context(), projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "章节不存在",
		})
		return
	}

	var req dto.UpdateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	updatedChapter, err := h.chapterService.Update(c.Request.Context(), chapter.ID.String(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.ChapterFromModel(updatedChapter)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// GenerateContent 生成章节内容
// @Summary 生成章节内容
// @Description 使用AI生成章节内容
// @Tags chapter
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterNumber path int true "章节号"
// @Param request body dto.GenerateChapterRequest true "生成请求"
// @Success 200 {object} dto.Response{data=dto.ChapterResponse}
// @Router /api/v1/projects/{id}/chapters/{chapterNumber}/generate [post]
func (h *ChapterHandler) GenerateContent(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	projectID := c.Param("id")
	chapterNumber, _ := strconv.Atoi(c.Param("chapterNumber"))

	// 先获取章节
	chapter, err := h.chapterService.GetByProjectAndNumber(c.Request.Context(), projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "章节不存在",
		})
		return
	}

	var req dto.GenerateChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	req.ChapterNumber = chapterNumber
	updatedChapter, err := h.chapterService.GenerateChapterContent(c.Request.Context(), deviceUUID, projectID, chapter.ID.String(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.ChapterFromModel(updatedChapter)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// Finalize 定稿章节
// @Summary 定稿章节
// @Description 定稿章节，更新全局摘要
// @Tags chapter
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterNumber path int true "章节号"
// @Param request body dto.FinalizeChapterRequest true "定稿请求"
// @Success 200 {object} dto.Response{data=dto.ChapterResponse}
// @Router /api/v1/projects/{id}/chapters/{chapterNumber}/finalize [post]
func (h *ChapterHandler) Finalize(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	projectID := c.Param("id")
	chapterNumber, _ := strconv.Atoi(c.Param("chapterNumber"))

	// 先获取章节
	chapter, err := h.chapterService.GetByProjectAndNumber(c.Request.Context(), projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "章节不存在",
		})
		return
	}

	var req dto.FinalizeChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	updatedChapter, err := h.chapterService.FinalizeChapter(c.Request.Context(), deviceUUID, projectID, chapter.ID.String(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.ChapterFromModel(updatedChapter)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// Enrich 扩写章节
// @Summary 扩写章节
// @Description 扩写章节内容到目标字数
// @Tags chapter
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param chapterNumber path int true "章节号"
// @Param request body dto.EnrichChapterRequest true "扩写请求"
// @Success 200 {object} dto.Response{data=dto.ChapterResponse}
// @Router /api/v1/projects/{id}/chapters/{chapterNumber}/enrich [post]
func (h *ChapterHandler) Enrich(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	projectID := c.Param("id")
	chapterNumber, _ := strconv.Atoi(c.Param("chapterNumber"))

	// 先获取章节
	chapter, err := h.chapterService.GetByProjectAndNumber(c.Request.Context(), projectID, chapterNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "章节不存在",
		})
		return
	}

	var req dto.EnrichChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	updatedChapter, err := h.chapterService.EnrichChapter(c.Request.Context(), deviceUUID, projectID, chapter.ID.String(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.ChapterFromModel(updatedChapter)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

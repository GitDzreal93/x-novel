package handler

import (
	"net/http"
	"strconv"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

// ProjectHandler 项目处理器
type ProjectHandler struct {
	projectService *service.ProjectService
}

// NewProjectHandler 创建项目处理器
func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// Create 创建项目
// @Summary 创建项目
// @Description 创建新的小说项目
// @Tags project
// @Accept json
// @Produce json
// @Param request body dto.CreateProjectRequest true "创建请求"
// @Success 200 {object} dto.Response{data=dto.ProjectResponse}
// @Router /api/v1/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	project, err := h.projectService.Create(c.Request.Context(), deviceUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "创建项目失败",
		})
		return
	}

	response := dto.FromModel(project)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// GetByID 获取项目详情
// @Summary 获取项目详情
// @Description 根据ID获取项目详情
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} dto.Response{data=dto.ProjectResponse}
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	project, err := h.projectService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "项目不存在",
		})
		return
	}

	response := dto.FromModel(project)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// List 获取项目列表
// @Summary 获取项目列表
// @Description 获取当前设备的项目列表
// @Tags project
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} dto.Response{data=dto.ProjectListResponse}
// @Router /api/v1/projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	projects, total, err := h.projectService.List(c.Request.Context(), deviceUUID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取项目列表失败",
		})
		return
	}

	projectResponses := make([]dto.ProjectResponse, 0, len(projects))
	for _, project := range projects {
		projectResponses = append(projectResponses, *dto.FromModel(project))
	}

	response := &dto.ProjectListResponse{
		Projects: projectResponses,
		Total:    int(total),
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// Update 更新项目
// @Summary 更新项目
// @Description 更新项目信息
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body dto.UpdateProjectRequest true "更新请求"
// @Success 200 {object} dto.Response{data=dto.ProjectResponse}
// @Router /api/v1/projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	project, err := h.projectService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "更新项目失败",
		})
		return
	}

	response := dto.FromModel(project)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// Delete 删除项目
// @Summary 删除项目
// @Description 删除项目
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Success 200 {object} dto.Response
// @Router /api/v1/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.projectService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "删除项目失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
	})
}

// GenerateArchitecture 生成小说架构
// @Summary 生成小说架构
// @Description 生成小说的架构信息
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body dto.GenerateArchitectureRequest true "生成请求"
// @Success 200 {object} dto.Response{data=dto.ProjectResponse}
// @Router /api/v1/projects/{id}/architecture/generate [post]
func (h *ProjectHandler) GenerateArchitecture(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	id := c.Param("id")

	var req dto.GenerateArchitectureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	project, err := h.projectService.GenerateArchitecture(c.Request.Context(), deviceUUID, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.FromModel(project)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// GenerateBlueprint 生成章节大纲
// @Summary 生成章节大纲
// @Description 生成章节大纲
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param request body dto.GenerateBlueprintRequest true "生成请求"
// @Success 200 {object} dto.Response{data=dto.ProjectResponse}
// @Router /api/v1/projects/{id}/blueprint/generate [post]
func (h *ProjectHandler) GenerateBlueprint(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	id := c.Param("id")

	var req dto.GenerateBlueprintRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	project, err := h.projectService.GenerateBlueprint(c.Request.Context(), deviceUUID, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := dto.FromModel(project)

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

// ExportProject 导出项目
// @Summary 导出项目
// @Description 导出项目为指定格式
// @Tags project
// @Accept json
// @Produce json
// @Param id path string true "项目ID"
// @Param format path string true "导出格式" Enums(txt, md)
// @Success 200 {object} dto.Response{data=dto.ExportResponse}
// @Router /api/v1/projects/{id}/export/{format} [get]
func (h *ProjectHandler) ExportProject(c *gin.Context) {
	id := c.Param("id")
	format := c.Param("format")

	downloadURL, err := h.projectService.ExportProject(c.Request.Context(), id, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	response := &dto.ExportResponse{
		DownloadURL: downloadURL,
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    response,
	})
}

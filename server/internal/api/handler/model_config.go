package handler

import (
	"net/http"
	"strconv"

	"x-novel/internal/api/middleware"
	"x-novel/internal/dto"
	"x-novel/internal/service"

	"github.com/gin-gonic/gin"
)

// ModelConfigHandler 模型配置处理器
type ModelConfigHandler struct {
	modelConfigService *service.ModelConfigService
}

// NewModelConfigHandler 创建模型配置处理器
func NewModelConfigHandler(modelConfigService *service.ModelConfigService) *ModelConfigHandler {
	return &ModelConfigHandler{
		modelConfigService: modelConfigService,
	}
}

// Create 创建模型配置
// @Summary 创建模型配置
// @Description 创建新的模型配置
// @Tags model-config
// @Accept json
// @Produce json
// @Param request body dto.CreateModelConfigRequest true "创建请求"
// @Success 200 {object} dto.Response{data=dto.ModelConfigResponse}
// @Router /api/v1/models [post]
func (h *ModelConfigHandler) Create(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	var req dto.CreateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	config, err := h.modelConfigService.Create(c.Request.Context(), deviceUUID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    dto.ModelConfigFromModel(config),
	})
}

// List 获取模型配置列表
// @Summary 获取模型配置列表
// @Description 获取当前设备的模型配置列表
// @Tags model-config
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} dto.Response{data=dto.ModelConfigListResponse}
// @Router /api/v1/models [get]
func (h *ModelConfigHandler) List(c *gin.Context) {
	deviceUUID, err := middleware.GetDeviceUUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "未授权",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	configs, total, err := h.modelConfigService.List(c.Request.Context(), deviceUUID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取模型配置列表失败",
		})
		return
	}

	configResponses := make([]dto.ModelConfigResponse, 0, len(configs))
	for _, config := range configs {
		configResponses = append(configResponses, *dto.ModelConfigFromModel(config))
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: dto.ModelConfigListResponse{
			Configs: configResponses,
			Total:   int(total),
		},
	})
}

// GetByID 获取模型配置详情
// @Summary 获取模型配置详情
// @Description 根据ID获取模型配置详情
// @Tags model-config
// @Accept json
// @Produce json
// @Param id path string true "配置ID"
// @Success 200 {object} dto.Response{data=dto.ModelConfigResponse}
// @Router /api/v1/models/{id} [get]
func (h *ModelConfigHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	config, err := h.modelConfigService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "模型配置不存在",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    dto.ModelConfigFromModel(config),
	})
}

// Update 更新模型配置
// @Summary 更新模型配置
// @Description 更新模型配置信息
// @Tags model-config
// @Accept json
// @Produce json
// @Param id path string true "配置ID"
// @Param request body dto.UpdateModelConfigRequest true "更新请求"
// @Success 200 {object} dto.Response{data=dto.ModelConfigResponse}
// @Router /api/v1/models/{id} [put]
func (h *ModelConfigHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	config, err := h.modelConfigService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    dto.ModelConfigFromModel(config),
	})
}

// Delete 删除模型配置
// @Summary 删除模型配置
// @Description 删除模型配置
// @Tags model-config
// @Accept json
// @Produce json
// @Param id path string true "配置ID"
// @Success 200 {object} dto.Response
// @Router /api/v1/models/{id} [delete]
func (h *ModelConfigHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.modelConfigService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "删除模型配置失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
	})
}

// ListProviders 获取提供商列表
// @Summary 获取提供商列表
// @Description 获取所有可用的模型提供商
// @Tags model-config
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response{data=[]dto.ModelProviderResponse}
// @Router /api/v1/models/providers [get]
func (h *ModelConfigHandler) ListProviders(c *gin.Context) {
	providers, err := h.modelConfigService.ListProviders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "获取提供商列表失败",
		})
		return
	}

	providerResponses := make([]dto.ModelProviderResponse, 0, len(providers))
	for _, provider := range providers {
		providerResponses = append(providerResponses, dto.ModelProviderResponse{
			ID:          provider.ID,
			Name:        provider.Name,
			DisplayName: provider.DisplayName,
			BaseURL:     provider.BaseURL,
			AuthType:    provider.AuthType,
		})
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    providerResponses,
	})
}

// Validate 验证模型配置
// @Summary 验证模型配置
// @Description 验证API Key是否有效
// @Tags model-config
// @Accept json
// @Produce json
// @Param request body dto.ValidateModelConfigRequest true "验证请求"
// @Success 200 {object} dto.Response
// @Router /api/v1/models/validate [post]
func (h *ModelConfigHandler) Validate(c *gin.Context) {
	var req dto.ValidateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "请求参数错误",
		})
		return
	}

	if err := h.modelConfigService.Validate(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Code:    http.StatusOK,
		Message: "验证成功",
	})
}

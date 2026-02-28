package dto

// ========== 项目相关 ==========

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	Title            string   `json:"title" binding:"required"`
	Topic            string   `json:"topic"`
	Genre            []string `json:"genre"`
	ChapterCount     int      `json:"chapter_count"`
	WordsPerChapter  int      `json:"words_per_chapter"`
	UserGuidance     string   `json:"user_guidance"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Title            *string  `json:"title"`
	Topic            *string  `json:"topic"`
	Genre            []string `json:"genre"`
	ChapterCount     *int     `json:"chapter_count"`
	WordsPerChapter  *int     `json:"words_per_chapter"`
	UserGuidance     *string  `json:"user_guidence"`
	Status           *string  `json:"status"`

	// 架构数据
	CoreSeed         *string `json:"core_seed"`
	CharacterDynamics *string `json:"character_dynamics"`
	WorldBuilding    *string `json:"world_building"`
	PlotArchitecture *string `json:"plot_architecture"`
	CharacterState   *string `json:"character_state"`

	// 大纲数据
	ChapterBlueprint *string `json:"chapter_blueprint"`

	// 全局摘要
	GlobalSummary    *string `json:"global_summary"`
}

// GenerateArchitectureRequest 生成架构请求
type GenerateArchitectureRequest struct {
	Overwrite bool `json:"overwrite"` // 是否覆盖已有架构
}

// GenerateBlueprintRequest 生成大纲请求
type GenerateBlueprintRequest struct {
	Overwrite bool `json:"overwrite"` // 是否覆盖已有大纲
}

// ExportProjectRequest 导出项目请求
type ExportProjectRequest struct {
	Format string `json:"format" binding:"required,oneof=txt md markdown"` // 导出格式
}

// ========== 章节相关 ==========

// CreateChapterRequest 创建章节请求
type CreateChapterRequest struct {
	ChapterNumber int    `json:"chapter_number" binding:"required"`
	Title         string `json:"title"`
	BlueprintSummary string `json:"blueprint_summary"`
}

// UpdateChapterRequest 更新章节请求
type UpdateChapterRequest struct {
	Title            *string `json:"title"`
	Content          *string `json:"content"`
	Status           *string `json:"status"`
	IsFinalized      *bool   `json:"is_finalized"`

	// 大纲信息
	BlueprintPosition    *string `json:"blueprint_position"`
	BlueprintPurpose     *string `json:"blueprint_purpose"`
	BlueprintSuspense    *string `json:"blueprint_suspense"`
	BlueprintForeshadowing *string `json:"blueprint_foreshadowing"`
	BlueprintTwistLevel  *string `json:"blueprint_twist_level"`
	BlueprintSummary     *string `json:"blueprint_summary"`
}

// GenerateChapterRequest 生成章节请求
type GenerateChapterRequest struct {
	ChapterNumber int    `json:"chapter_number"` // 从 URL 路径设置，非必填
	Overwrite     bool   `json:"overwrite"` // 是否覆盖已有内容
}

// FinalizeChapterRequest 定稿章节请求
type FinalizeChapterRequest struct {
	UpdateSummary bool `json:"update_summary"` // 是否更新全局摘要
}

// EnrichChapterRequest 扩写章节请求
type EnrichChapterRequest struct {
	TargetWords int    `json:"target_words"` // 目标字数
}

// ========== 模型配置相关 ==========

// CreateModelConfigRequest 创建模型配置请求
type CreateModelConfigRequest struct {
	ProviderID int    `json:"provider_id" binding:"required"`
	ModelName  string `json:"model_name" binding:"required"`
	Purpose    string `json:"purpose" binding:"required,oneof=architecture chapter writing review general"`
	APIKey     string `json:"api_key" binding:"required"`
	BaseURL    string `json:"base_url"`
}

// UpdateModelConfigRequest 更新模型配置请求
type UpdateModelConfigRequest struct {
	ModelName  *string `json:"model_name"`
	Purpose    *string `json:"purpose"`
	APIKey     *string `json:"api_key"`
	BaseURL    *string `json:"base_url"`
	IsActive   *bool   `json:"is_active"`
}

// ValidateModelConfigRequest 验证模型配置请求
type ValidateModelConfigRequest struct {
	ProviderID int    `json:"provider_id" binding:"required"`
	APIKey     string `json:"api_key" binding:"required"`
	BaseURL    string `json:"base_url"`
}

// ========== 设备相关 ==========

// UpdateDeviceSettingsRequest 更新设备设置请求
type UpdateDeviceSettingsRequest struct {
	Theme            *string `json:"theme"`
	Language         *string `json:"language"`
	AutoSaveEnabled  *bool   `json:"auto_save_enabled"`
	AutoSaveInterval *int    `json:"auto_save_interval"`
}

package dto

import (
	"time"

	"x-novel/internal/model"

	"github.com/google/uuid"
)

// ========== 通用响应 ==========

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorDetail 错误详情
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

// ========== 项目响应 ==========

// ProjectResponse 项目响应
type ProjectResponse struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Topic     string     `json:"topic"`
	Genre     []string   `json:"genre"`
	ChapterCount     int    `json:"chapter_count"`
	WordsPerChapter  int    `json:"words_per_chapter"`
	UserGuidance     string `json:"user_guidance"`

	// 架构数据
	CoreSeed            string `json:"core_seed,omitempty"`
	CharacterDynamics   string `json:"character_dynamics,omitempty"`
	WorldBuilding       string `json:"world_building,omitempty"`
	PlotArchitecture    string `json:"plot_architecture,omitempty"`
	CharacterState      string `json:"character_state,omitempty"`
	ArchitectureGenerated bool  `json:"architecture_generated"`

	// 大纲数据
	ChapterBlueprint   string `json:"chapter_blueprint,omitempty"`
	BlueprintGenerated bool   `json:"blueprint_generated"`

	// 统计
	GlobalSummary      string `json:"global_summary,omitempty"`
	TotalChapters      int    `json:"total_chapters"`
	CompletedChapters  int    `json:"completed_chapters"`
	TotalWords         int    `json:"total_words"`

	// 状态
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ProjectListResponse 项目列表响应
type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int               `json:"total"`
}

// ========== 章节响应 ==========

// ChapterResponse 章节响应
type ChapterResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	ChapterNumber int    `json:"chapter_number"`
	Title     string    `json:"title"`

	// 大纲信息
	BlueprintPosition    string `json:"blueprint_position,omitempty"`
	BlueprintPurpose     string `json:"blueprint_purpose,omitempty"`
	BlueprintSuspense    string `json:"blueprint_suspense,omitempty"`
	BlueprintForeshadowing string `json:"blueprint_foreshadowing,omitempty"`
	BlueprintTwistLevel  string `json:"blueprint_twist_level,omitempty"`
	BlueprintSummary     string `json:"blueprint_summary,omitempty"`

	// 内容
	Content   string `json:"content,omitempty"`
	WordCount int    `json:"word_count"`

	// 状态
	Status      string `json:"status"`
	IsFinalized bool   `json:"is_finalized"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ChapterListResponse 章节列表响应
type ChapterListResponse struct {
	Chapters []ChapterResponse `json:"chapters"`
	Total    int               `json:"total"`
}

// ========== 模型配置响应 ==========

// ModelConfigResponse 模型配置响应
type ModelConfigResponse struct {
	ID         uuid.UUID `json:"id"`
	ProviderID int       `json:"provider_id"`
	ModelName  string    `json:"model_name"`
	Purpose    string    `json:"purpose"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联的提供商信息
	Provider   *ModelProviderResponse `json:"provider,omitempty"`
}

// ModelConfigListResponse 模型配置列表响应
type ModelConfigListResponse struct {
	Configs []ModelConfigResponse `json:"configs"`
	Total   int                   `json:"total"`
}

// ModelProviderResponse 模型提供商响应
type ModelProviderResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	IsActive    bool   `json:"is_active"`
}

// ========== 设备响应 ==========

// DeviceResponse 设备响应
type DeviceResponse struct {
	ID        uuid.UUID              `json:"id"`
	DeviceID  string                 `json:"device_id"`
	DeviceInfo map[string]interface{} `json:"device_info,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	LastSeen  time.Time              `json:"last_seen"`
}

// DeviceSettingResponse 设备设置响应
type DeviceSettingResponse struct {
	ID               uuid.UUID `json:"id"`
	DeviceID         uuid.UUID `json:"device_id"`
	Theme            string    `json:"theme"`
	Language         string    `json:"language"`
	AutoSaveEnabled  bool      `json:"auto_save_enabled"`
	AutoSaveInterval int       `json:"auto_save_interval"`
}

// ========== 生成任务响应 ==========

// GenerationTaskResponse 生成任务响应
type GenerationTaskResponse struct {
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"`    // pending, processing, completed, failed
	Progress  int       `json:"progress"`  // 0-100
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ========== 导出响应 ==========

// ExportResponse 导出响应
type ExportResponse struct {
	DownloadURL string `json:"download_url"`
	FileSize    int64  `json:"file_size"`
	WordCount   int    `json:"word_count"`
	ChapterCount int  `json:"chapter_count"`
}

// ========== 辅助函数 ==========

// FromModel 从模型转换为响应
func FromModel(p *model.Project) *ProjectResponse {
	totalChapters := len(p.Chapters)
	completedChapters := 0
	totalWords := 0

	for _, ch := range p.Chapters {
		if ch.Status == "completed" {
			completedChapters++
		}
		totalWords += ch.WordCount
	}

	return &ProjectResponse{
		ID:                   p.ID,
		Title:                p.Title,
		Topic:                p.Topic,
		Genre:                p.Genre,
		ChapterCount:         p.ChapterCount,
		WordsPerChapter:      p.WordsPerChapter,
		UserGuidance:         p.UserGuidance,
		CoreSeed:             p.CoreSeed,
		CharacterDynamics:    p.CharacterDynamics,
		WorldBuilding:        p.WorldBuilding,
		PlotArchitecture:     p.PlotArchitecture,
		CharacterState:       p.CharacterState,
		ArchitectureGenerated: p.ArchitectureGenerated,
		ChapterBlueprint:     p.ChapterBlueprint,
		BlueprintGenerated:   p.BlueprintGenerated,
		GlobalSummary:        p.GlobalSummary,
		TotalChapters:        totalChapters,
		CompletedChapters:    completedChapters,
		TotalWords:           totalWords,
		Status:               p.Status,
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}
}

// ChapterFromModel 从模型转换为章节响应
func ChapterFromModel(c *model.Chapter) *ChapterResponse {
	return &ChapterResponse{
		ID:                   c.ID,
		ProjectID:            c.ProjectID,
		ChapterNumber:        c.ChapterNumber,
		Title:                c.Title,
		BlueprintPosition:    c.BlueprintPosition,
		BlueprintPurpose:     c.BlueprintPurpose,
		BlueprintSuspense:    c.BlueprintSuspense,
		BlueprintForeshadowing: c.BlueprintForeshadowing,
		BlueprintTwistLevel:  c.BlueprintTwistLevel,
		BlueprintSummary:     c.BlueprintSummary,
		Content:              c.Content,
		WordCount:            c.WordCount,
		Status:               c.Status,
		IsFinalized:          c.IsFinalized,
		CreatedAt:            c.CreatedAt,
		UpdatedAt:            c.UpdatedAt,
	}
}

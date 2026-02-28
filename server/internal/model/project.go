package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DeviceID  uuid.UUID `gorm:"type:uuid;not null;index" json:"device_id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Topic     string    `gorm:"type:text" json:"topic,omitempty"`
	Genre     string    `gorm:"type:text" json:"genre,omitempty"` // 存储 JSON 字符串
	ChapterCount     int       `gorm:"default:100" json:"chapter_count"`
	WordsPerChapter  int       `gorm:"default:3000" json:"words_per_chapter"`
	UserGuidance     string    `gorm:"type:text" json:"user_guidance,omitempty"`

	// 架构数据
	CoreSeed            string `gorm:"type:text" json:"core_seed,omitempty"`
	CharacterDynamics   string `gorm:"type:text" json:"character_dynamics,omitempty"`
	WorldBuilding       string `gorm:"type:text" json:"world_building,omitempty"`
	PlotArchitecture    string `gorm:"type:text" json:"plot_architecture,omitempty"`
	CharacterState      string `gorm:"type:text" json:"character_state,omitempty"`
	ArchitectureGenerated bool   `gorm:"default:false" json:"architecture_generated"`

	// 大纲数据
	ChapterBlueprint      string `gorm:"type:text" json:"chapter_blueprint,omitempty"`
	BlueprintGenerated    bool   `gorm:"default:false" json:"blueprint_generated"`

	// 上下文数据
	GlobalSummary string `gorm:"type:text" json:"global_summary,omitempty"`

	// 关系图谱数据
	GraphData     string `gorm:"type:text" json:"graph_data,omitempty"` // 存储 JSON 字符串

	// 状态
	Status     string    `gorm:"size:20;default:draft" json:"status"` // draft, writing, completed, published
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联
	Device      *Device      `gorm:"-" json:"device,omitempty"`
	Chapters    []Chapter    `gorm:"-" json:"chapters,omitempty"`
}

func (Project) TableName() string {
	return "projects"
}

// BeforeCreate GORM hook
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// Chapter 章节模型
type Chapter struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index:idx_project_chapter" json:"project_id"`
	ChapterNumber int    `gorm:"not null;index:idx_project_chapter" json:"chapter_number"`
	Title     string    `gorm:"size:200" json:"title,omitempty"`

	// 大纲信息
	BlueprintPosition    string `gorm:"size:100" json:"blueprint_position,omitempty"`
	BlueprintPurpose     string `gorm:"size:100" json:"blueprint_purpose,omitempty"`
	BlueprintSuspense    string `gorm:"size:100" json:"blueprint_suspense,omitempty"`
	BlueprintForeshadowing  string `gorm:"type:text" json:"blueprint_foreshadowing,omitempty"`
	BlueprintTwistLevel  string `gorm:"size:20" json:"blueprint_twist_level,omitempty"`
	BlueprintSummary     string `gorm:"type:text" json:"blueprint_summary,omitempty"`

	// 章节内容
	Content    string `gorm:"type:text" json:"content,omitempty"`
	WordCount  int    `gorm:"default:0" json:"word_count"`

	// 状态
	Status     string `gorm:"size:20;default:not_started" json:"status"` // not_started, draft, completed
	IsFinalized bool   `gorm:"default:false" json:"is_finalized"`

	// 分析数据
	Analysis     string `gorm:"type:text" json:"analysis,omitempty"` // 存储 JSON 字符串

	// 关系图谱
	ChapterGraph string `gorm:"type:text" json:"chapter_graph,omitempty"` // 存储 JSON 字符串

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联
	Project    *Project `gorm:"-" json:"project,omitempty"`
}

func (Chapter) TableName() string {
	return "chapters"
}

// BeforeCreate GORM hook
func (c *Chapter) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

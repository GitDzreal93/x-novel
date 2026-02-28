package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"x-novel/internal/dto"
	"x-novel/internal/llm"
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ProjectService 项目服务
type ProjectService struct {
	projectRepo  *repository.ProjectRepository
	chapterRepo  *repository.ChapterRepository
	modelRepo    *repository.ModelConfigRepository
	llmManager   *llm.Manager
	exportService *ExportService
}

// NewProjectService 创建项目服务
func NewProjectService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
	exportService *ExportService,
) *ProjectService {
	return &ProjectService{
		projectRepo:  projectRepo,
		chapterRepo:  chapterRepo,
		modelRepo:    modelRepo,
		llmManager:   llmManager,
		exportService: exportService,
	}
}

// Create 创建项目
func (s *ProjectService) Create(ctx context.Context, deviceID uuid.UUID, req *dto.CreateProjectRequest) (*model.Project, error) {
	// 将 Genre 数组转换为 JSON 字符串
	genreJSON, err := json.Marshal(req.Genre)
	if err != nil {
		logger.Error("序列化 Genre 失败", zap.Error(err))
		return nil, err
	}

	project := &model.Project{
		DeviceID:         deviceID,
		Title:            req.Title,
		Topic:            req.Topic,
		Genre:            string(genreJSON),
		ChapterCount:     req.ChapterCount,
		WordsPerChapter:  req.WordsPerChapter,
		UserGuidance:     req.UserGuidance,
		Status:           "draft",
	}

	// 设置默认值
	if project.ChapterCount == 0 {
		project.ChapterCount = 100
	}
	if project.WordsPerChapter == 0 {
		project.WordsPerChapter = 3000
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		logger.Error("创建项目失败", zap.Error(err))
		return nil, err
	}

	return project, nil
}

// GetByID 根据 ID 获取项目
func (s *ProjectService) GetByID(ctx context.Context, id string) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("获取项目失败",
			zap.String("project_id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return project, nil
}

// List 获取项目列表
func (s *ProjectService) List(ctx context.Context, deviceID uuid.UUID, page, pageSize int) ([]*model.Project, int64, error) {
	offset := (page - 1) * pageSize
	projects, total, err := s.projectRepo.List(ctx, deviceID.String(), offset, pageSize)
	if err != nil {
		logger.Error("获取项目列表失败",
			zap.String("device_id", deviceID.String()),
			zap.Error(err),
		)
		return nil, 0, err
	}
	return projects, total, nil
}

// Update 更新项目
func (s *ProjectService) Update(ctx context.Context, id string, req *dto.UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		project.Title = *req.Title
	}
	if req.Topic != nil {
		project.Topic = *req.Topic
	}
	if req.Genre != nil {
		genreJSON, err := json.Marshal(req.Genre)
		if err != nil {
			logger.Error("序列化 Genre 失败", zap.Error(err))
			return nil, err
		}
		project.Genre = string(genreJSON)
	}
	if req.ChapterCount != nil {
		project.ChapterCount = *req.ChapterCount
	}
	if req.WordsPerChapter != nil {
		project.WordsPerChapter = *req.WordsPerChapter
	}
	if req.UserGuidance != nil {
		project.UserGuidance = *req.UserGuidance
	}
	if req.Status != nil {
		project.Status = *req.Status
	}

	// 更新架构数据
	if req.CoreSeed != nil {
		project.CoreSeed = *req.CoreSeed
	}
	if req.CharacterDynamics != nil {
		project.CharacterDynamics = *req.CharacterDynamics
	}
	if req.WorldBuilding != nil {
		project.WorldBuilding = *req.WorldBuilding
	}
	if req.PlotArchitecture != nil {
		project.PlotArchitecture = *req.PlotArchitecture
	}
	if req.CharacterState != nil {
		project.CharacterState = *req.CharacterState
	}

	// 更新大纲数据
	if req.ChapterBlueprint != nil {
		project.ChapterBlueprint = *req.ChapterBlueprint
	}

	// 更新全局摘要
	if req.GlobalSummary != nil {
		project.GlobalSummary = *req.GlobalSummary
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		logger.Error("更新项目失败",
			zap.String("project_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return project, nil
}

// Delete 删除项目
func (s *ProjectService) Delete(ctx context.Context, id string) error {
	if err := s.projectRepo.Delete(ctx, id); err != nil {
		logger.Error("删除项目失败",
			zap.String("project_id", id),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GenerateArchitecture 生成小说架构
func (s *ProjectService) GenerateArchitecture(ctx context.Context, deviceID uuid.UUID, projectID string, req *dto.GenerateArchitectureRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 检查是否已经生成过
	if project.ArchitectureGenerated && !req.Overwrite {
		return nil, errors.New("架构已生成，如需重新生成请设置 overwrite=true")
	}

	logger.Info("开始生成小说架构",
		zap.String("project_id", projectID),
	)

	// 解析 Genre
	var genres []string
	if project.Genre != "" {
		if err := json.Unmarshal([]byte(project.Genre), &genres); err != nil {
			logger.Error("解析 Genre 失败", zap.Error(err))
			return nil, err
		}
	}

	// 获取用于架构生成的模型配置
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "architecture")
	if err != nil {
		logger.Error("获取架构生成模型配置失败", zap.Error(err))

		// 开发模式：返回模拟数据
		logger.Info("使用模拟模式生成架构")
		return s.generateMockArchitecture(ctx, project, genres)
	}

	// 准备提示词参数
	promptParams := ArchitecturePromptParams{
		Topic:          project.Topic,
		Genre:          genres,
		ChapterCount:   project.ChapterCount,
		WordsPerChapter: project.WordsPerChapter,
		UserGuidance:   project.UserGuidance,
	}

	// 步骤1: 生成核心种子
	logger.Info("步骤1: 生成核心种子")
	coreSeed, err := s.generateArchitectureStep(ctx, modelConfig, "core_seed", promptParams, "")
	if err != nil {
		logger.Error("生成核心种子失败", zap.Error(err))
		return nil, err
	}
	project.CoreSeed = coreSeed
	promptParams.CoreSeed = coreSeed
	logger.Info("核心种子生成成功", zap.String("core_seed", coreSeed))

	// 步骤2: 生成角色动力学
	logger.Info("步骤2: 生成角色动力学")
	characterDynamics, err := s.generateArchitectureStep(ctx, modelConfig, "character_dynamics", promptParams, "")
	if err != nil {
		logger.Error("生成角色动力学失败", zap.Error(err))
		return nil, err
	}
	project.CharacterDynamics = characterDynamics
	promptParams.CharacterDynamics = characterDynamics
	logger.Info("角色动力学生成成功")

	// 步骤3: 生成世界观
	logger.Info("步骤3: 生成世界观")
	worldBuilding, err := s.generateArchitectureStep(ctx, modelConfig, "world_building", promptParams, "")
	if err != nil {
		logger.Error("生成世界观失败", zap.Error(err))
		return nil, err
	}
	project.WorldBuilding = worldBuilding
	promptParams.WorldBuilding = worldBuilding
	logger.Info("世界观生成成功")

	// 步骤4: 生成情节架构
	logger.Info("步骤4: 生成情节架构")
	plotArchitecture, err := s.generateArchitectureStep(ctx, modelConfig, "plot_architecture", promptParams, "")
	if err != nil {
		logger.Error("生成情节架构失败", zap.Error(err))
		return nil, err
	}
	project.PlotArchitecture = plotArchitecture
	logger.Info("情节架构生成成功")

	// 步骤5: 生成角色状态
	logger.Info("步骤5: 生成角色状态")
	characterState, err := s.generateArchitectureStep(ctx, modelConfig, "character_state", promptParams, "")
	if err != nil {
		logger.Error("生成角色状态失败", zap.Error(err))
		return nil, err
	}
	project.CharacterState = characterState
	logger.Info("角色状态生成成功")

	// 标记架构已生成
	project.ArchitectureGenerated = true

	// 保存更新
	if err := s.projectRepo.Update(ctx, project); err != nil {
		logger.Error("保存架构数据失败", zap.Error(err))
		return nil, err
	}

	logger.Info("小说架构生成完成",
		zap.String("project_id", projectID),
	)

	return project, nil
}

// generateArchitectureStep 执行单个架构生成步骤
func (s *ProjectService) generateArchitectureStep(ctx context.Context, modelConfig *model.ModelConfig, step string, params ArchitecturePromptParams, systemPrompt string) (string, error) {
	// 构建用户提示词
	userPrompt := GetArchitecturePrompt(step, params)

	// 构建消息
	messages := []llm.ChatMessage{}
	if systemPrompt != "" {
		messages = append(messages, llm.ChatMessage{Role: "system", Content: systemPrompt})
	}
	messages = append(messages, llm.ChatMessage{Role: "user", Content: userPrompt})

	// 调用 LLM
	provider := "openai"
	if modelConfig.Provider != nil {
		provider = modelConfig.Provider.Name
	}

	options := llm.ChatOptions{
		Temperature: 0.7,
		MaxTokens:   4000,
		APIKey:      modelConfig.APIKey,
	}

	// 如果有自定义 BaseURL，需要创建临时适配器
	var result string
	var err error
	if modelConfig.BaseURL != "" {
		adapter := llm.NewOpenAIAdapter(modelConfig.BaseURL, modelConfig.ModelName)
		result, err = adapter.ChatCompletion(ctx, messages, options)
	} else {
		result, err = s.llmManager.ChatCompletion(ctx, provider, messages, options)
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

// generateMockArchitecture 生成模拟架构数据（用于开发测试）
func (s *ProjectService) generateMockArchitecture(ctx context.Context, project *model.Project, genres []string) (*model.Project, error) {
	genreStr := "通用"
	if len(genres) > 0 {
		genreStr = genres[0]
	}

	// 根据不同类型生成不同的模拟内容
	project.CoreSeed = fmt.Sprintf("当程序员李明在代码中发现一个神秘的AI程序，获得了超能力，必须阻止隐藏在暗处的黑暗组织，否则整个世界将被数字化控制。")

	project.CharacterDynamics = fmt.Sprintf(`【%s】核心角色设计：

**李明（主角）**
- 特征：28岁男性，程序员，戴着黑框眼镜，身材瘦弱
- 核心驱动力三角：
  * 表面追求：升职加薪，买学区房
  * 深层渴望：获得认可，证明自己的价值
  * 灵魂需求：寻找生命的意义和责任
- 角色弧线：普通程序员 → 发现超能力 → 被迫卷入斗争 → 接受使命 → 成为守护者

**苏小雨（女主角）**
- 特征：26岁女性，神秘组织的前成员，短发干练，眼神锐利
- 核心驱动力三角：
  * 表面追求：逃离组织控制
  * 深层渴望：信任与被信任
  * 灵魂需求：救赎过去的错误
- 角色弧线：组织杀手 → 遇见主角 → 内心转变 → 选择正义 → 与主角并肩作战

**黑暗组织首领（反派）**
- 特征：50岁男性，科技巨头CEO，外表儒雅但心机深沉
- 核心驱动力三角：
  * 表面追求：掌控全球科技
  * 深层渴望：永生与不朽
  * 灵魂需求：超越人类的极限
`, genreStr)

	project.WorldBuilding = fmt.Sprintf(`【%s】世界观设定：

**物理维度**
- 空间结构：近未来的都市设定，充满高科技感与传统城市的对比
- 时间背景：2035年，AI技术已经渗透到生活的方方面面
- 规则体系：超能力来源于一种神秘的量子算法，可以通过代码操控现实

**社会维度**
- 社会结构：科技巨头掌控社会资源，普通人成为数字时代的"无产阶级"
- 文化氛围：虚拟与现实界限模糊，人们更愿意生活在虚拟世界中
- 生活方式：高度依赖AI助手，人际关系逐渐数字化

**情感维度**
- 核心意象：代码、数据流、虚拟与现实的重叠
- 环境氛围：赛博朋克式的压抑与希望并存
- 情感体验：科技冷漠中的温情，数字时代的真实情感
`, genreStr)

	project.PlotArchitecture = fmt.Sprintf(`【%s】三幕式情节架构：

**第一幕（开端）**
1. 日常状态：李明在公司加班，处理bug，生活平淡无奇（3处场景铺垫）
2. 引出故事：
   - 主线：发现AI程序中的异常代码
   - 感情线：与苏小雨初次相遇
   - 副线：公司的秘密项目露出端倪
3. 契机事件：李明意外激活AI程序，获得超能力
4. 初步反应：震惊、恐惧，试图隐藏这个秘密

**第二幕（发展）**
1. 剧情深入：
   - 黑暗组织察觉到李明的存在
   - 苏小雨接近李明，身份成谜
   - 李明逐渐掌握超能力
2. 挑战与成长：
   - 李明被组织追杀
   - 第一次使用超能力战斗
   - 与苏小雨建立信任
3. 情感升温：两人在逃亡中产生感情
4. 重要转折：苏小雨揭露身份，李明陷入信任危机

**第三幕（高潮与结局）**
1. 核心冲突爆发：黑暗组织启动"数字永生"计划
2. 角色抉择：李明选择牺牲自己的超能力拯救世界
3. 结局收尾：组织被瓦解，李明回归普通生活，但内心已成长
`, genreStr)

	project.CharacterState = fmt.Sprintf(`【角色状态文档】

**李明**
├── 物品:
│  ├── 黑色笔记本：记录着神秘的AI代码
│  └── 破旧眼镜：实际上是一副AR眼镜，可以显示数据流
├── 能力
│  ├── 技能1：代码操控现实：可以通过编写代码改变物理世界
│  └── 技能2：数据感知：能够看到数字世界的流动
├── 状态
│  ├── 身体状态：普通程序员体格，经常加班熬夜
│  └── 心理状态：从平凡到震撼，正在慢慢接受新的身份
├── 主要角色间关系网
│  ├── 苏小雨：从怀疑到信任，再到相互依靠
│  └── 黑暗组织：被追杀的对象，也是最终要对抗的敌人
├── 触发或加深的事件
│  ├── 发现AI程序：改变命运的关键时刻
│  └── 苏小雨揭露身份：信任的考验，感情的转折
`)

	// 标记架构已生成
	project.ArchitectureGenerated = true

	// 保存更新
	if err := s.projectRepo.Update(ctx, project); err != nil {
		logger.Error("保存模拟架构数据失败", zap.Error(err))
		return nil, err
	}

	logger.Info("模拟架构生成完成",
		zap.String("project_id", project.ID.String()),
	)

	return project, nil
}

// GenerateBlueprint 生成章节大纲
func (s *ProjectService) GenerateBlueprint(ctx context.Context, deviceID uuid.UUID, projectID string, req *dto.GenerateBlueprintRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 检查架构是否已生成
	if !project.ArchitectureGenerated {
		return nil, errors.New("请先生成小说架构")
	}

	// 检查是否已经生成过
	if project.BlueprintGenerated && !req.Overwrite {
		return nil, errors.New("大纲已生成，如需重新生成请设置 overwrite=true")
	}

	logger.Info("开始生成章节大纲",
		zap.String("project_id", projectID),
	)

	// 获取用于大纲生成的模型配置
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "blueprint")
	if err != nil {
		logger.Error("获取大纲生成模型配置失败", zap.Error(err))

		// 开发模式：返回模拟数据
		logger.Info("使用模拟模式生成章节大纲")
		return s.generateMockBlueprint(ctx, project)
	}

	// TODO: 实现真实的 LLM 调用生成大纲
	_ = modelConfig
	return nil, errors.New("大纲生成功能待实现")
}

// generateMockBlueprint 生成模拟章节大纲（用于开发测试）
func (s *ProjectService) generateMockBlueprint(ctx context.Context, project *model.Project) (*model.Project, error) {
	// 准备参数
	params := BlueprintPromptParams{
		UserGuidance:      project.UserGuidance,
		CoreSeed:          project.CoreSeed,
		CharacterDynamics: project.CharacterDynamics,
		WorldBuilding:     project.WorldBuilding,
		PlotArchitecture:  project.PlotArchitecture,
		ChapterCount:      project.ChapterCount,
	}

	// 生成模拟大纲
	project.ChapterBlueprint = GenerateMockBlueprint(params)
	project.BlueprintGenerated = true

	// 保存更新
	if err := s.projectRepo.Update(ctx, project); err != nil {
		logger.Error("保存模拟大纲数据失败", zap.Error(err))
		return nil, err
	}

	logger.Info("模拟章节大纲生成完成",
		zap.String("project_id", project.ID.String()),
		zap.Int("chapter_count", project.ChapterCount),
	)

	return project, nil
}

// ExportProject 导出项目
func (s *ProjectService) ExportProject(ctx context.Context, projectID, format string) (string, error) {
	return s.exportService.ExportProject(ctx, projectID, ExportFormat(format))
}

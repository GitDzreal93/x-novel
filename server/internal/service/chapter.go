package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"x-novel/internal/dto"
	"x-novel/internal/llm"
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ChapterService 章节服务
type ChapterService struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
	modelRepo   *repository.ModelConfigRepository
	llmManager  *llm.Manager
}

// NewChapterService 创建章节服务
func NewChapterService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
) *ChapterService {
	return &ChapterService{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
		modelRepo:   modelRepo,
		llmManager:  llmManager,
	}
}

// Create 创建章节
func (s *ChapterService) Create(ctx context.Context, projectID string, req *dto.CreateChapterRequest) (*model.Chapter, error) {
	// 检查项目是否存在
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, errors.New("项目不存在")
	}

	// 检查章节号是否已存在
	existingChapter, _ := s.chapterRepo.GetByProjectAndNumber(ctx, projectID, req.ChapterNumber)
	if existingChapter != nil {
		return nil, errors.New("章节号已存在")
	}

	chapter := &model.Chapter{
		ProjectID:        project.ID,
		ChapterNumber:    req.ChapterNumber,
		Title:            req.Title,
		BlueprintSummary: req.BlueprintSummary,
		Status:           "not_started",
	}

	if err := s.chapterRepo.Create(ctx, chapter); err != nil {
		logger.Error("创建章节失败", zap.Error(err))
		return nil, err
	}

	return chapter, nil
}

// GetByID 根据 ID 获取章节
func (s *ChapterService) GetByID(ctx context.Context, id string) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		logger.Error("获取章节失败",
			zap.String("chapter_id", id),
			zap.Error(err),
		)
		return nil, err
	}
	return chapter, nil
}

// GetByProjectAndNumber 根据项目 ID 和章节号获取章节
func (s *ChapterService) GetByProjectAndNumber(ctx context.Context, projectID string, chapterNumber int) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByProjectAndNumber(ctx, projectID, chapterNumber)
	if err != nil {
		logger.Error("获取章节失败",
			zap.String("project_id", projectID),
			zap.Int("chapter_number", chapterNumber),
			zap.Error(err),
		)
		return nil, err
	}
	return chapter, nil
}

// List 获取章节列表
func (s *ChapterService) List(ctx context.Context, projectID string, page, pageSize int) ([]*model.Chapter, int64, error) {
	offset := (page - 1) * pageSize
	chapters, total, err := s.chapterRepo.List(ctx, projectID, offset, pageSize)
	if err != nil {
		logger.Error("获取章节列表失败",
			zap.String("project_id", projectID),
			zap.Error(err),
		)
		return nil, 0, err
	}
	return chapters, total, nil
}

// Update 更新章节
func (s *ChapterService) Update(ctx context.Context, id string, req *dto.UpdateChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		chapter.Title = *req.Title
	}
	if req.Content != nil {
		chapter.Content = *req.Content
		chapter.WordCount = utf8.RuneCountInString(*req.Content)
	}
	if req.Status != nil {
		chapter.Status = *req.Status
	}
	if req.IsFinalized != nil {
		chapter.IsFinalized = *req.IsFinalized
	}

	// 更新大纲信息
	if req.BlueprintPosition != nil {
		chapter.BlueprintPosition = *req.BlueprintPosition
	}
	if req.BlueprintPurpose != nil {
		chapter.BlueprintPurpose = *req.BlueprintPurpose
	}
	if req.BlueprintSuspense != nil {
		chapter.BlueprintSuspense = *req.BlueprintSuspense
	}
	if req.BlueprintForeshadowing != nil {
		chapter.BlueprintForeshadowing = *req.BlueprintForeshadowing
	}
	if req.BlueprintTwistLevel != nil {
		chapter.BlueprintTwistLevel = *req.BlueprintTwistLevel
	}
	if req.BlueprintSummary != nil {
		chapter.BlueprintSummary = *req.BlueprintSummary
	}

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		logger.Error("更新章节失败",
			zap.String("chapter_id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return chapter, nil
}

// Delete 删除章节
func (s *ChapterService) Delete(ctx context.Context, id string) error {
	if err := s.chapterRepo.Delete(ctx, id); err != nil {
		logger.Error("删除章节失败",
			zap.String("chapter_id", id),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GenerateChapterContent 生成章节内容
func (s *ChapterService) GenerateChapterContent(ctx context.Context, deviceID uuid.UUID, projectID, chapterID string, req *dto.GenerateChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 检查是否已有内容
	if chapter.Content != "" && !req.Overwrite {
		return nil, errors.New("章节已有内容，如需重新生成请设置 overwrite=true")
	}

	logger.Info("开始生成章节内容",
		zap.String("project_id", projectID),
		zap.String("chapter_id", chapterID),
		zap.Int("chapter_number", req.ChapterNumber),
	)

	// 获取项目信息
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 获取用于章节生成的模型配置
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "chapter")
	if err != nil {
		logger.Error("获取章节生成模型配置失败", zap.Error(err))

		// 开发模式：返回模拟数据
		logger.Info("使用模拟模式生成章节内容")
		return s.generateMockChapterContent(ctx, project, chapter, req)
	}

	// 解析 Genre
	var genres []string
	if project.Genre != "" {
		if err := json.Unmarshal([]byte(project.Genre), &genres); err != nil {
			logger.Error("解析 Genre 失败", zap.Error(err))
		}
	}

	// 构建提示词参数
	params := ChapterPromptParams{
		Title:             project.Title,
		Topic:             project.Topic,
		Genre:             genres,
		UserGuidance:      project.UserGuidance,
		WordsPerChapter:   project.WordsPerChapter,
		CoreSeed:          project.CoreSeed,
		CharacterDynamics: project.CharacterDynamics,
		WorldBuilding:     project.WorldBuilding,
		PlotArchitecture:  project.PlotArchitecture,
		CharacterState:    project.CharacterState,
		ChapterNumber:     chapter.ChapterNumber,
		ChapterTitle:      chapter.Title,
		BlueprintSummary:  chapter.BlueprintSummary,
		GlobalSummary:     project.GlobalSummary,
	}

	// 获取提示词
	prompt := GetChapterPrompt(chapter.ChapterNumber, params)

	// 构建消息
	messages := []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}

	// 调用 LLM
	options := llm.ChatOptions{
		Temperature: 0.8,
		MaxTokens:   8000,
		APIKey:      modelConfig.APIKey,
	}

	var content string
	if modelConfig.BaseURL != "" {
		adapter := llm.NewOpenAIAdapter(modelConfig.BaseURL, modelConfig.ModelName)
		content, err = adapter.ChatCompletion(ctx, messages, options)
	} else {
		provider := "openai"
		if modelConfig.Provider != nil {
			provider = modelConfig.Provider.Name
		}
		content, err = s.llmManager.ChatCompletion(ctx, provider, messages, options)
	}

	if err != nil {
		logger.Error("LLM 调用失败", zap.Error(err))
		return nil, fmt.Errorf("生成章节内容失败: %w", err)
	}

	// 更新章节
	chapter.Content = content
	chapter.WordCount = utf8.RuneCountInString(content)
	chapter.Status = "draft"

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		logger.Error("保存章节内容失败", zap.Error(err))
		return nil, err
	}

	logger.Info("章节内容生成完成",
		zap.String("chapter_id", chapterID),
		zap.Int("word_count", chapter.WordCount),
	)

	return chapter, nil
}

// generateMockChapterContent 生成模拟章节内容
func (s *ChapterService) generateMockChapterContent(ctx context.Context, project *model.Project, chapter *model.Chapter, req *dto.GenerateChapterRequest) (*model.Chapter, error) {
	// 生成模拟章节内容
	chapter.Content = s.generateMockChapterText(project, chapter.ChapterNumber)
	chapter.WordCount = utf8.RuneCountInString(chapter.Content)
	chapter.Status = "draft"

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		logger.Error("保存模拟章节内容失败", zap.Error(err))
		return nil, err
	}

	logger.Info("模拟章节内容生成完成",
		zap.String("chapter_id", chapter.ID.String()),
		zap.Int("word_count", chapter.WordCount),
	)

	return chapter, nil
}

// generateMockChapterText 生成模拟章节文本
func (s *ChapterService) generateMockChapterText(project *model.Project, chapterNumber int) string {
	// 根据章节号生成不同内容
	if chapterNumber == 1 {
		return `深夜的办公室里，只剩下李明一个人还在加班。显示器的蓝光映在他疲惫的脸上，键盘的敲击声在空荡的办公室里回荡。

"终于找到这个bug了。"李明长舒一口气，靠在椅背上。

就在他准备关电脑下班的时候，一段奇怪的代码引起了他的注意。这段代码像是凭空出现的，隐藏在一个不起眼的工具函数里。李明好奇地把这段代码复制到编辑器里，仔细查看起来。

代码的逻辑很简单，但却透着一股诡异的美感。它像是一个神经网络，又像是一种加密算法。李明鬼使神差地按下了运行键。

就在那一瞬间，世界变了。

李明眼前的景象开始扭曲，所有的物体都变成了流动的数据流。他看到电脑里的数据像潮水一样涌出，在空中编织成复杂的图案。桌上的咖啡杯分解成无数个粒子，然后又重新组合。

"这...这是怎么回事？"李明惊恐地看着自己的双手，它们也在散发着淡淡的光芒。

就在这时，办公室的门被推开了。一个穿着黑色西装的男人走了进来，他的眼神冰冷，仿佛在看一个死人。

"你动了不该动的东西，"男人说道，声音里没有一丝感情。`
	}

	if chapterNumber == 2 {
		return `李明本能地向后退去，但发现自己的身体像是被什么东西束缚住了。那个黑衣男人抬起手，李明就感觉一股强大的力量将自己提了起来，悬在半空中。

"放下我！"李明挣扎着，但那股力量纹丝不动。

"你以为这是什么？儿童玩具吗？"黑衣男人冷笑着，"这是改变世界的钥匙，而你，刚好打开了它。"

李明的心跳加速，他意识到自己可能卷入了一个巨大的麻烦。他闭上眼睛，试图理解刚才发生的一切。那些流动的数据，那些变化的景象...难道自己获得了超能力？

就在这个念头闪过的瞬间，李明感觉到周围的一切都慢了下来。他看到了空气中流动的数据粒子，看到了黑衣男人身上发出的能量波动。

"放开他！"

一个清脆的女声打破了沉默。李明睁开眼睛，看到一个年轻女子冲了进来。她留着短发，眼神锐利，动作敏捷。

黑衣男人皱了皱眉，"苏小雨，你背叛组织了？"

"我只是不想看到无辜的人被伤害，"苏小雨说道，同时向李明使了个眼色，"快走！"

李明不知道自己哪来的勇气，但他感觉到一股力量在体内涌动。他猛地挣脱了束缚，向门口冲去。黑衣男人想要追赶，但被苏小雨拦住了。

"走！"苏小雨喊道。

李明冲出了办公室，在走廊里狂奔。他的心脏剧烈地跳动着，脑海里充满了问号。这到底是怎么回事？自己为什么会卷入这种事情？那个叫苏小雨的女人又是谁？`
	}

	// 通用章节内容
	return fmt.Sprintf(`第%d章的内容正在生成中...

李明走在空荡的街道上，夜风吹过，带来一丝凉意。他还在思考着发生的一切，那些超能力的真相，那个神秘的组织，还有苏小雨...

"这一切都太不真实了，"李明自言自语道。

但他知道，这一切都是真实的。从那一刻他运行了那段代码，他的人生就彻底改变了。现在，他必须面对这个全新的世界，学会使用自己的能力，并且找到活下去的方法。

街道尽头的路灯闪烁着，像是某种信号。李明停下脚步，感觉到了什么。有人在跟踪他。

"出来吧，"李明转过身，"我知道你在那里。"

黑暗中，一个身影缓缓走了出来...`, chapterNumber)
}

// FinalizeChapter 定稿章节
func (s *ChapterService) FinalizeChapter(ctx context.Context, deviceID uuid.UUID, projectID, chapterID string, req *dto.FinalizeChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 检查章节是否有内容
	if chapter.Content == "" {
		return nil, errors.New("章节内容为空，无法定稿")
	}

	// 设置为定稿状态
	chapter.IsFinalized = true
	chapter.Status = "completed"

	if err := s.chapterRepo.SetFinalized(ctx, chapterID, true); err != nil {
		return nil, err
	}

	logger.Info("章节定稿完成",
		zap.String("chapter_id", chapterID),
		zap.Int("word_count", chapter.WordCount),
	)

	// 如果需要更新全局摘要
	if req.UpdateSummary {
		// 获取所有已定稿的章节
		chapters, err := s.chapterRepo.ListCompleted(ctx, projectID, 0)
		if err == nil && len(chapters) > 0 {
			// 生成全局摘要
			globalSummary := s.generateGlobalSummary(chapters)
			if err := s.projectRepo.UpdateGlobalSummary(ctx, projectID, globalSummary); err != nil {
				logger.Warn("更新全局摘要失败", zap.Error(err))
			}
		}
	}

	return chapter, nil
}

// generateGlobalSummary 生成全局摘要
func (s *ChapterService) generateGlobalSummary(chapters []*model.Chapter) string {
	var summary strings.Builder
	summary.WriteString("// 全局剧情摘要\n\n")

	for _, chapter := range chapters {
		summary.WriteString(fmt.Sprintf("第%d章：%s\n", chapter.ChapterNumber, chapter.Title))
		if chapter.BlueprintSummary != "" {
			summary.WriteString(fmt.Sprintf("摘要：%s\n\n", chapter.BlueprintSummary))
		}
	}

	return summary.String()
}

// EnrichChapter 扩写章节
func (s *ChapterService) EnrichChapter(ctx context.Context, deviceID uuid.UUID, projectID, chapterID string, req *dto.EnrichChapterRequest) (*model.Chapter, error) {
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return nil, err
	}

	// 检查章节是否有内容
	if chapter.Content == "" {
		return nil, errors.New("章节内容为空，无法扩写")
	}

	logger.Info("开始扩写章节内容",
		zap.String("project_id", projectID),
		zap.String("chapter_id", chapterID),
		zap.Int("target_words", req.TargetWords),
	)

	// 获取用于章节生成的模型配置
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "chapter")
	if err != nil {
		logger.Error("获取章节生成模型配置失败", zap.Error(err))

		// 开发模式：返回模拟数据
		logger.Info("使用模拟模式扩写章节内容")
		return s.enrichMockChapter(ctx, chapter, req.TargetWords)
	}

	// 获取项目信息
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// 解析 Genre
	var genres []string
	if project.Genre != "" {
		if err := json.Unmarshal([]byte(project.Genre), &genres); err != nil {
			logger.Error("解析 Genre 失败", zap.Error(err))
		}
	}

	// 计算目标字数
	targetWords := req.TargetWords
	if targetWords == 0 {
		targetWords = chapter.WordCount + 500 // 默认增加 500 字
	}

	// 构建提示词参数
	params := ChapterPromptParams{
		Genre:          genres,
		CurrentContent: chapter.Content,
		TargetWords:    targetWords,
	}

	// 获取扩写提示词
	prompt := GetEnrichPrompt(params)

	// 构建消息
	messages := []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}

	// 调用 LLM
	options := llm.ChatOptions{
		Temperature: 0.7,
		MaxTokens:   10000,
		APIKey:      modelConfig.APIKey,
	}

	var content string
	if modelConfig.BaseURL != "" {
		adapter := llm.NewOpenAIAdapter(modelConfig.BaseURL, modelConfig.ModelName)
		content, err = adapter.ChatCompletion(ctx, messages, options)
	} else {
		provider := "openai"
		if modelConfig.Provider != nil {
			provider = modelConfig.Provider.Name
		}
		content, err = s.llmManager.ChatCompletion(ctx, provider, messages, options)
	}

	if err != nil {
		logger.Error("LLM 调用失败", zap.Error(err))
		return nil, fmt.Errorf("扩写章节内容失败: %w", err)
	}

	// 更新章节
	chapter.Content = content
	chapter.WordCount = utf8.RuneCountInString(content)

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		logger.Error("保存扩写章节内容失败", zap.Error(err))
		return nil, err
	}

	logger.Info("章节扩写完成",
		zap.String("chapter_id", chapterID),
		zap.Int("word_count", chapter.WordCount),
	)

	return chapter, nil
}

// enrichMockChapter 模拟扩写章节
func (s *ChapterService) enrichMockChapter(ctx context.Context, chapter *model.Chapter, targetWords int) (*model.Chapter, error) {
	// 简单模拟：在原文基础上添加一些描述性文字
	enrichedContent := chapter.Content + `

李明深深吸了一口气，空气中弥漫着夜晚特有的清冷。街道两旁的路灯投下昏黄的光影，将他的影子拉得很长。他能听到远处传来的车流声，那是城市在这个时间点唯一的声响。

这一刻，李明突然意识到，自己再也无法回到过去那个平凡的生活了。那些数据流，那些超能力，都已经成为了他生命的一部分。无论是好是坏，他都必须接受这个事实。

他握紧了拳头，感受着体内涌动的力量。这力量既让他感到强大，也让他感到恐惧。他不知道该如何控制它，也不知道它最终会将他带向何方。

但有一点是确定的：他不能放弃。无论前方等待他的是什么，他都必须勇敢地面对。`

	chapter.Content = enrichedContent
	chapter.WordCount = utf8.RuneCountInString(enrichedContent)

	if err := s.chapterRepo.Update(ctx, chapter); err != nil {
		logger.Error("保存扩写章节内容失败", zap.Error(err))
		return nil, err
	}

	logger.Info("模拟章节扩写完成",
		zap.String("chapter_id", chapter.ID.String()),
		zap.Int("new_word_count", chapter.WordCount),
	)

	return chapter, nil
}

// GetPreviousChapters 获取前面已完成的所有章节（用于生成前文摘要）
func (s *ChapterService) GetPreviousChapters(ctx context.Context, projectID string, chapterNumber int) ([]*model.Chapter, error) {
	chapters, err := s.chapterRepo.ListCompleted(ctx, projectID, chapterNumber)
	if err != nil {
		logger.Error("获取前面章节失败",
			zap.String("project_id", projectID),
			zap.Int("chapter_number", chapterNumber),
			zap.Error(err),
		)
		return nil, err
	}
	return chapters, nil
}

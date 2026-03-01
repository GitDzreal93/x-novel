package service

import (
	"context"
	"fmt"

	"x-novel/internal/llm"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WritingAssistantService struct {
	projectRepo *repository.ProjectRepository
	chapterRepo *repository.ChapterRepository
	modelRepo   *repository.ModelConfigRepository
	llmManager  *llm.Manager
}

func NewWritingAssistantService(
	projectRepo *repository.ProjectRepository,
	chapterRepo *repository.ChapterRepository,
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
) *WritingAssistantService {
	return &WritingAssistantService{
		projectRepo: projectRepo,
		chapterRepo: chapterRepo,
		modelRepo:   modelRepo,
		llmManager:  llmManager,
	}
}

// Polish 润色文本
func (s *WritingAssistantService) Polish(ctx context.Context, deviceID uuid.UUID, content, style string) (string, error) {
	prompt := GetPolishPrompt(content, style)
	result, err := s.callLLM(ctx, deviceID, prompt, 0.7)
	if err != nil {
		logger.Error("润色失败", zap.Error(err))
		return s.mockPolish(content, style), nil
	}
	return result, nil
}

// Continue 续写文本
func (s *WritingAssistantService) Continue(ctx context.Context, deviceID uuid.UUID, projectID, content string, targetWords int) (string, error) {
	projectContext := s.getProjectContext(ctx, projectID)
	prompt := GetContinuePrompt(content, targetWords, projectContext)
	result, err := s.callLLM(ctx, deviceID, prompt, 0.85)
	if err != nil {
		logger.Error("续写失败", zap.Error(err))
		return s.mockContinue(content, targetWords), nil
	}
	return result, nil
}

// Suggest 提供灵感建议
func (s *WritingAssistantService) Suggest(ctx context.Context, deviceID uuid.UUID, projectID, content, aspect string) (string, error) {
	projectContext := s.getProjectContext(ctx, projectID)
	prompt := GetSuggestionPrompt(content, aspect, projectContext)
	result, err := s.callLLM(ctx, deviceID, prompt, 0.9)
	if err != nil {
		logger.Error("灵感建议失败", zap.Error(err))
		return s.mockSuggestion(aspect), nil
	}
	return result, nil
}

// PolishStream 流式润色
func (s *WritingAssistantService) PolishStream(ctx context.Context, deviceID uuid.UUID, content, style string, callback llm.StreamCallback) (string, error) {
	prompt := GetPolishPrompt(content, style)
	result, err := s.callLLMStream(ctx, deviceID, prompt, 0.7, callback)
	if err != nil {
		mock := s.mockPolish(content, style)
		if cbErr := callback(mock); cbErr != nil {
			return "", cbErr
		}
		return mock, nil
	}
	return result, nil
}

// ContinueStream 流式续写
func (s *WritingAssistantService) ContinueStream(ctx context.Context, deviceID uuid.UUID, projectID, content string, targetWords int, callback llm.StreamCallback) (string, error) {
	projectContext := s.getProjectContext(ctx, projectID)
	prompt := GetContinuePrompt(content, targetWords, projectContext)
	result, err := s.callLLMStream(ctx, deviceID, prompt, 0.85, callback)
	if err != nil {
		mock := s.mockContinue(content, targetWords)
		if cbErr := callback(mock); cbErr != nil {
			return "", cbErr
		}
		return mock, nil
	}
	return result, nil
}

// SuggestStream 流式灵感建议
func (s *WritingAssistantService) SuggestStream(ctx context.Context, deviceID uuid.UUID, projectID, content, aspect string, callback llm.StreamCallback) (string, error) {
	projectContext := s.getProjectContext(ctx, projectID)
	prompt := GetSuggestionPrompt(content, aspect, projectContext)
	result, err := s.callLLMStream(ctx, deviceID, prompt, 0.9, callback)
	if err != nil {
		mock := s.mockSuggestion(aspect)
		if cbErr := callback(mock); cbErr != nil {
			return "", cbErr
		}
		return mock, nil
	}
	return result, nil
}

func (s *WritingAssistantService) callLLM(ctx context.Context, deviceID uuid.UUID, prompt string, temperature float32) (string, error) {
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "writing")
	if err != nil {
		modelConfig, err = s.modelRepo.GetByPurpose(ctx, deviceID.String(), "general")
		if err != nil {
			return "", fmt.Errorf("未配置写作/通用模型: %w", err)
		}
	}

	messages := []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}

	options := llm.ChatOptions{
		Temperature: temperature,
		MaxTokens:   4096,
		APIKey:      modelConfig.APIKey,
	}

	if modelConfig.BaseURL != "" {
		adapter := llm.NewOpenAIAdapter(modelConfig.BaseURL, modelConfig.ModelName)
		return adapter.ChatCompletion(ctx, messages, options)
	}

	provider := "openai"
	if modelConfig.Provider != nil {
		provider = modelConfig.Provider.Name
	}
	return s.llmManager.ChatCompletion(ctx, provider, messages, options)
}

func (s *WritingAssistantService) callLLMStream(ctx context.Context, deviceID uuid.UUID, prompt string, temperature float32, callback llm.StreamCallback) (string, error) {
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "writing")
	if err != nil {
		modelConfig, err = s.modelRepo.GetByPurpose(ctx, deviceID.String(), "general")
		if err != nil {
			return "", fmt.Errorf("未配置写作/通用模型: %w", err)
		}
	}

	messages := []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}

	options := llm.ChatOptions{
		Temperature: temperature,
		MaxTokens:   4096,
		APIKey:      modelConfig.APIKey,
	}

	if modelConfig.BaseURL != "" {
		adapter := llm.NewOpenAIAdapter(modelConfig.BaseURL, modelConfig.ModelName)
		return adapter.StreamChatCompletion(ctx, messages, options, callback)
	}

	provider := "openai"
	if modelConfig.Provider != nil {
		provider = modelConfig.Provider.Name
	}
	return s.llmManager.StreamChatCompletion(ctx, provider, messages, options, callback)
}

func (s *WritingAssistantService) getProjectContext(ctx context.Context, projectID string) string {
	if projectID == "" {
		return ""
	}
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return ""
	}
	context := fmt.Sprintf("标题：%s", project.Title)
	if project.Topic != "" {
		context += fmt.Sprintf("，主题：%s", project.Topic)
	}
	if project.CoreSeed != "" {
		seed := project.CoreSeed
		runes := []rune(seed)
		if len(runes) > 300 {
			seed = string(runes[:300]) + "..."
		}
		context += fmt.Sprintf("\n核心设定：%s", seed)
	}
	if project.GlobalSummary != "" {
		summary := project.GlobalSummary
		runes := []rune(summary)
		if len(runes) > 500 {
			summary = string(runes[:500]) + "..."
		}
		context += fmt.Sprintf("\n前文摘要：%s", summary)
	}
	return context
}

func (s *WritingAssistantService) mockPolish(content, style string) string {
	runes := []rune(content)
	preview := string(runes)
	if len(runes) > 100 {
		preview = string(runes[:100]) + "..."
	}
	return fmt.Sprintf("【润色结果 - %s风格】\n\n%s\n\n（以上为模拟润色结果，配置 LLM 后将获得真实 AI 润色）", style, preview)
}

func (s *WritingAssistantService) mockContinue(content string, targetWords int) string {
	return fmt.Sprintf(`他抬起头，目光穿过昏暗的灯光，落在远处那扇紧闭的门上。一种莫名的不安涌上心头，像是有什么重要的事情即将发生。

空气中弥漫着一股淡淡的潮湿气味，雨后的城市总是带着这种特殊的味道。街道上行人稀少，偶尔有一辆车驶过，轮胎碾过水洼的声音在寂静中格外清晰。

"也许我应该回去，"他自言自语，但脚步却不由自主地向前迈去。

（以上为模拟续写约 %d 字，配置 LLM 后将基于上下文生成连贯内容）`, targetWords)
}

func (s *WritingAssistantService) mockSuggestion(aspect string) string {
	return fmt.Sprintf(`## 创作建议

### 建议一：强化内心独白
通过第一人称的心理描写，让读者更深入地了解角色的内心世界。可以在关键决策点加入角色的犹豫、挣扎和最终的决断过程。

### 建议二：增加感官细节
不要仅仅描述视觉画面，尝试加入听觉、嗅觉、触觉等多维度的感官体验，让场景更加立体生动。

### 建议三：制造悬念钩子
在章节结尾处留下一个未解之谜或意外转折，激发读者继续阅读的欲望。

### 建议四：对比与反差
通过人物性格、环境氛围或情节走向的反差，增强故事的戏剧性和吸引力。

（以上为模拟建议 - %s方向，配置 LLM 后将基于实际内容生成针对性建议）`, aspect)
}

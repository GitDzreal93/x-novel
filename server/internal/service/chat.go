package service

import (
	"context"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"x-novel/internal/llm"
	"x-novel/internal/model"
	"x-novel/internal/repository"
	"x-novel/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChatService struct {
	chatRepo    *repository.ChatRepository
	projectRepo *repository.ProjectRepository
	modelRepo   *repository.ModelConfigRepository
	llmManager  *llm.Manager
}

func NewChatService(
	chatRepo *repository.ChatRepository,
	projectRepo *repository.ProjectRepository,
	modelRepo *repository.ModelConfigRepository,
	llmManager *llm.Manager,
) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		projectRepo: projectRepo,
		modelRepo:   modelRepo,
		llmManager:  llmManager,
	}
}

// CreateConversation 创建对话
func (s *ChatService) CreateConversation(ctx context.Context, deviceID uuid.UUID, title, mode string, projectID *uuid.UUID) (*model.Conversation, error) {
	if title == "" {
		title = "新对话"
	}
	if mode == "" {
		mode = string(ChatModeGeneral)
	}

	conv := &model.Conversation{
		DeviceID:  deviceID,
		Title:     title,
		Mode:      mode,
		ProjectID: projectID,
	}

	if err := s.chatRepo.CreateConversation(ctx, conv); err != nil {
		logger.Error("创建对话失败", zap.Error(err))
		return nil, err
	}

	return conv, nil
}

// ListConversations 获取对话列表
func (s *ChatService) ListConversations(ctx context.Context, deviceID uuid.UUID, page, pageSize int) ([]*model.Conversation, int64, error) {
	offset := (page - 1) * pageSize
	return s.chatRepo.ListConversations(ctx, deviceID, offset, pageSize)
}

// GetConversation 获取对话详情（含消息）
func (s *ChatService) GetConversation(ctx context.Context, id string) (*model.Conversation, error) {
	conv, err := s.chatRepo.GetConversation(ctx, id)
	if err != nil {
		return nil, err
	}

	messages, err := s.chatRepo.ListMessages(ctx, id)
	if err != nil {
		return nil, err
	}

	conv.Messages = make([]model.Message, len(messages))
	for i, m := range messages {
		conv.Messages[i] = *m
	}

	return conv, nil
}

// DeleteConversation 删除对话
func (s *ChatService) DeleteConversation(ctx context.Context, id string) error {
	return s.chatRepo.DeleteConversation(ctx, id)
}

// UpdateConversationTitle 更新对话标题
func (s *ChatService) UpdateConversationTitle(ctx context.Context, id, title string) error {
	return s.chatRepo.UpdateConversationTitle(ctx, id, title)
}

// SendMessage 发送消息并获取 AI 回复（非流式）
func (s *ChatService) SendMessage(ctx context.Context, deviceID uuid.UUID, conversationID, content string) (*model.Message, *model.Message, error) {
	conv, err := s.chatRepo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("对话不存在: %w", err)
	}

	// 保存用户消息
	userMsg := &model.Message{
		ConversationID: conv.ID,
		Role:           "user",
		Content:        content,
	}
	if err := s.chatRepo.CreateMessage(ctx, userMsg); err != nil {
		return nil, nil, err
	}

	// 构建 LLM 消息列表
	llmMessages, err := s.buildLLMMessages(ctx, conv, content)
	if err != nil {
		return nil, nil, err
	}

	// 调用 LLM
	replyContent, err := s.callLLM(ctx, deviceID, llmMessages)
	if err != nil {
		logger.Error("LLM 对话失败", zap.Error(err))
		replyContent = s.getMockReply(ChatMode(conv.Mode), content)
	}

	// 保存 AI 回复
	assistantMsg := &model.Message{
		ConversationID: conv.ID,
		Role:           "assistant",
		Content:        replyContent,
	}
	if err := s.chatRepo.CreateMessage(ctx, assistantMsg); err != nil {
		return nil, nil, err
	}

	_ = s.chatRepo.TouchConversation(ctx, conversationID)

	// 如果是第一条消息，自动更新对话标题
	messages, _ := s.chatRepo.ListMessages(ctx, conversationID)
	if len(messages) == 2 {
		newTitle := s.generateTitle(content)
		_ = s.chatRepo.UpdateConversationTitle(ctx, conversationID, newTitle)
	}

	return userMsg, assistantMsg, nil
}

// SendMessageStream 发送消息并流式获取 AI 回复
func (s *ChatService) SendMessageStream(ctx context.Context, deviceID uuid.UUID, conversationID, content string, callback llm.StreamCallback) (*model.Message, *model.Message, error) {
	conv, err := s.chatRepo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("对话不存在: %w", err)
	}

	userMsg := &model.Message{
		ConversationID: conv.ID,
		Role:           "user",
		Content:        content,
	}
	if err := s.chatRepo.CreateMessage(ctx, userMsg); err != nil {
		return nil, nil, err
	}

	llmMessages, err := s.buildLLMMessages(ctx, conv, content)
	if err != nil {
		return nil, nil, err
	}

	replyContent, err := s.callLLMStream(ctx, deviceID, llmMessages, callback)
	if err != nil {
		logger.Error("LLM 流式对话失败", zap.Error(err))
		replyContent = s.getMockReply(ChatMode(conv.Mode), content)
		if cbErr := callback(replyContent); cbErr != nil {
			return nil, nil, cbErr
		}
	}

	assistantMsg := &model.Message{
		ConversationID: conv.ID,
		Role:           "assistant",
		Content:        replyContent,
	}
	if err := s.chatRepo.CreateMessage(ctx, assistantMsg); err != nil {
		return nil, nil, err
	}

	_ = s.chatRepo.TouchConversation(ctx, conversationID)

	messages, _ := s.chatRepo.ListMessages(ctx, conversationID)
	if len(messages) == 2 {
		newTitle := s.generateTitle(content)
		_ = s.chatRepo.UpdateConversationTitle(ctx, conversationID, newTitle)
	}

	return userMsg, assistantMsg, nil
}

func (s *ChatService) buildLLMMessages(ctx context.Context, conv *model.Conversation, userContent string) ([]llm.ChatMessage, error) {
	var messages []llm.ChatMessage

	// 构建项目上下文
	var projectContext string
	if conv.ProjectID != nil {
		project, err := s.projectRepo.GetByID(ctx, conv.ProjectID.String())
		if err == nil {
			var genres []string
			if project.Genre != "" {
				json.Unmarshal([]byte(project.Genre), &genres)
			}
			projectContext = fmt.Sprintf("- 小说标题：%s\n- 主题：%s\n- 类型：%v\n- 每章字数：%d",
				project.Title, project.Topic, genres, project.WordsPerChapter)
			if project.CoreSeed != "" {
				projectContext += fmt.Sprintf("\n- 核心设定：%s", truncate(project.CoreSeed, 500))
			}
		}
	}

	// System prompt
	systemPrompt := GetChatSystemPrompt(ChatMode(conv.Mode), projectContext)
	messages = append(messages, llm.ChatMessage{Role: "system", Content: systemPrompt})

	// 历史消息（最近 20 条）
	history, err := s.chatRepo.RecentMessages(ctx, conv.ID.String(), 20)
	if err != nil {
		return nil, err
	}
	for _, m := range history {
		messages = append(messages, llm.ChatMessage{Role: m.Role, Content: m.Content})
	}

	// 当前用户消息
	messages = append(messages, llm.ChatMessage{Role: "user", Content: userContent})

	return messages, nil
}

func (s *ChatService) callLLM(ctx context.Context, deviceID uuid.UUID, messages []llm.ChatMessage) (string, error) {
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "general")
	if err != nil {
		return "", fmt.Errorf("未配置通用模型: %w", err)
	}

	options := llm.ChatOptions{
		Temperature: 0.85,
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

func (s *ChatService) callLLMStream(ctx context.Context, deviceID uuid.UUID, messages []llm.ChatMessage, callback llm.StreamCallback) (string, error) {
	modelConfig, err := s.modelRepo.GetByPurpose(ctx, deviceID.String(), "general")
	if err != nil {
		return "", fmt.Errorf("未配置通用模型: %w", err)
	}

	options := llm.ChatOptions{
		Temperature: 0.85,
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

func (s *ChatService) generateTitle(firstMessage string) string {
	runes := []rune(firstMessage)
	if len(runes) > 20 {
		return string(runes[:20]) + "..."
	}
	return firstMessage
}

func (s *ChatService) getMockReply(mode ChatMode, content string) string {
	wordCount := utf8.RuneCountInString(content)
	switch mode {
	case ChatModeCreative:
		return fmt.Sprintf("## 创意灵感\n\n基于你的想法「%s」，我提供以下几个创意方向：\n\n"+
			"**方向一：反转设定**\n可以尝试将主角的核心特质反转，制造戏剧冲突...\n\n"+
			"**方向二：多重视角**\n从配角或反派的视角重新审视故事，可能会发现新的叙事空间...\n\n"+
			"**方向三：时间错位**\n尝试打破线性叙事，通过时间跳跃增加悬念...\n\n"+
			"你觉得哪个方向更感兴趣？我可以深入展开。", truncateStr(content, 50))
	case ChatModeBuilding:
		return fmt.Sprintf("## 设定完善建议\n\n关于「%s」(%d字)，从世界观角度分析：\n\n"+
			"1. **内部一致性**：需要确保这个设定与已有体系不冲突\n"+
			"2. **延伸可能**：这个设定可以自然引申出更多有趣的子系统\n"+
			"3. **读者体验**：建议通过具体场景来展示设定，而非直接说明\n\n"+
			"需要我帮你细化其中某个方面吗？", truncateStr(content, 50), wordCount)
	case ChatModeCharacter:
		return fmt.Sprintf("## 角色分析\n\n关于「%s」，从角色塑造角度：\n\n"+
			"**性格层次**\n- 表层特征：外在表现和社交面具\n- 深层动机：真正驱动角色的内在需求\n- 内在矛盾：让角色立体的关键冲突\n\n"+
			"**成长弧线建议**\n角色需要经历一个从「不自知」到「自我觉醒」的过程，这会让读者产生共鸣。\n\n"+
			"你想具体讨论哪个方面？", truncateStr(content, 50))
	default:
		return fmt.Sprintf("感谢你的提问！关于「%s」，这是一个很好的创作方向。\n\n"+
			"从专业角度来看，我建议：\n"+
			"1. 先梳理核心冲突和主题\n"+
			"2. 确保角色动机清晰合理\n"+
			"3. 注意节奏把控，张弛有度\n\n"+
			"你可以选择以下模式获得更专业的帮助：\n"+
			"- **创意启发**：激发灵感和新想法\n"+
			"- **设定完善**：完善世界观和设定\n"+
			"- **角色塑造**：深入人物刻画\n\n"+
			"需要我从哪个角度深入分析？", truncateStr(content, 50))
	}
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return s
}

func truncateStr(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return s
}

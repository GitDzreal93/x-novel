package llm

import (
	"context"
)

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

// ChatOptions 聊天选项
type ChatOptions struct {
	Temperature float32 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
	APIKey      string  `json:"-"` // API Key，不序列化到 JSON
}

// StreamCallback 流式响应回调
type StreamCallback func(chunk string) error

// LLMAdapter LLM 适配器接口
type LLMAdapter interface {
	// ChatCompletion 聊天补全
	ChatCompletion(ctx context.Context, messages []ChatMessage, options ChatOptions) (string, error)

	// StreamChatCompletion 流式聊天补全
	StreamChatCompletion(ctx context.Context, messages []ChatMessage, options ChatOptions, callback StreamCallback) (string, error)

	// ValidateConfig 验证配置
	ValidateConfig(apiKey, baseURL string) error

	// GetDefaultModel 获取默认模型名称
	GetDefaultModel() string
}

// Manager LLM 管理器
type Manager struct {
	adapters map[string]LLMAdapter
}

// NewManager 创建 LLM 管理器
func NewManager() *Manager {
	return &Manager{
		adapters: make(map[string]LLMAdapter),
	}
}

// Register 注册适配器
func (m *Manager) Register(provider string, adapter LLMAdapter) {
	m.adapters[provider] = adapter
}

// Get 获取适配器
func (m *Manager) Get(provider string) (LLMAdapter, bool) {
	adapter, ok := m.adapters[provider]
	return adapter, ok
}

// GetDefault 获取默认适配器
func (m *Manager) GetDefault() (LLMAdapter, bool) {
	return m.Get("openai")
}

// ChatCompletion 聊天补全（使用指定提供商）
func (m *Manager) ChatCompletion(ctx context.Context, provider string, messages []ChatMessage, options ChatOptions) (string, error) {
	adapter, ok := m.Get(provider)
	if !ok {
		adapter, ok = m.GetDefault()
		if !ok {
			return "", ErrNoAdapterAvailable
		}
	}

	return adapter.ChatCompletion(ctx, messages, options)
}

// StreamChatCompletion 流式聊天补全
func (m *Manager) StreamChatCompletion(ctx context.Context, provider string, messages []ChatMessage, options ChatOptions, callback StreamCallback) (string, error) {
	adapter, ok := m.Get(provider)
	if !ok {
		adapter, ok = m.GetDefault()
		if !ok {
			return "", ErrNoAdapterAvailable
		}
	}

	return adapter.StreamChatCompletion(ctx, messages, options, callback)
}

// Errors
var (
	ErrNoAdapterAvailable = &LLMError{Message: "没有可用的 LLM 适配器"}
	ErrInvalidAPIKey      = &LLMError{Message: "无效的 API Key"}
	ErrRateLimitExceeded  = &LLMError{Message: "API 请求频率超限"}
	ErrInvalidResponse    = &LLMError{Message: "无效的响应"}
)

// LLMError LLM 错误
type LLMError struct {
	Message string
	Err     error
}

func (e *LLMError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *LLMError) Unwrap() error {
	return e.Err
}

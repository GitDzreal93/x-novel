package llm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIAdapter 基于 go-openai 库的适配器，兼容所有 OpenAI 协议的服务
type OpenAIAdapter struct {
	baseURL string
	model   string
}

// NewOpenAIAdapter 创建 OpenAI 适配器
func NewOpenAIAdapter(baseURL, model string) *OpenAIAdapter {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIAdapter{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		model:   model,
	}
}

func (a *OpenAIAdapter) newClient(apiKey string) *openai.Client {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = a.baseURL
	cfg.HTTPClient = &http.Client{Timeout: 120 * time.Second}
	return openai.NewClientWithConfig(cfg)
}

func toOpenAIMessages(messages []ChatMessage) []openai.ChatCompletionMessage {
	out := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		out[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	return out
}

// ChatCompletion 聊天补全
func (a *OpenAIAdapter) ChatCompletion(ctx context.Context, messages []ChatMessage, options ChatOptions) (string, error) {
	client := a.newClient(options.APIKey)

	req := openai.ChatCompletionRequest{
		Model:       a.model,
		Messages:    toOpenAIMessages(messages),
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("LLM 请求失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", ErrInvalidResponse
	}

	return resp.Choices[0].Message.Content, nil
}

// StreamChatCompletion 流式聊天补全
func (a *OpenAIAdapter) StreamChatCompletion(ctx context.Context, messages []ChatMessage, options ChatOptions, callback StreamCallback) (string, error) {
	client := a.newClient(options.APIKey)

	req := openai.ChatCompletionRequest{
		Model:       a.model,
		Messages:    toOpenAIMessages(messages),
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Stream:      true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", fmt.Errorf("LLM 流式请求失败: %w", err)
	}
	defer stream.Close()

	var fullContent strings.Builder

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fullContent.String(), fmt.Errorf("流式读取失败: %w", err)
		}

		if len(resp.Choices) > 0 {
			content := resp.Choices[0].Delta.Content
			if content != "" {
				fullContent.WriteString(content)
				if cbErr := callback(content); cbErr != nil {
					return fullContent.String(), cbErr
				}
			}
		}
	}

	return fullContent.String(), nil
}

// ValidateConfig 验证模型配置（发送真实测试请求）
func (a *OpenAIAdapter) ValidateConfig(apiKey, baseURL string) error {
	if apiKey == "" {
		return ErrInvalidAPIKey
	}

	adapter := NewOpenAIAdapter(baseURL, a.model)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := adapter.ChatCompletion(ctx, []ChatMessage{
		{Role: "user", Content: "Hi"},
	}, ChatOptions{
		APIKey:      apiKey,
		Temperature: 0.1,
		MaxTokens:   5,
	})

	if err != nil {
		return &LLMError{Message: "连接验证失败", Err: err}
	}
	return nil
}

// GetDefaultModel 获取默认模型
func (a *OpenAIAdapter) GetDefaultModel() string {
	return a.model
}

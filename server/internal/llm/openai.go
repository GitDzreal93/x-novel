package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OpenAIAdapter OpenAI 适配器
type OpenAIAdapter struct {
	client  *http.Client
	baseURL string
	model   string
}

// NewOpenAIAdapter 创建 OpenAI 适配器
func NewOpenAIAdapter(baseURL, model string) *OpenAIAdapter {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIAdapter{
		client:  &http.Client{},
		baseURL: strings.TrimSuffix(baseURL, "/"),
		model:   model,
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float32       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 选择
type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message,omitempty"`
	Delta        *Delta      `json:"delta,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

// Delta 增量（流式）
type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// Usage 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletion 聊天补全
func (a *OpenAIAdapter) ChatCompletion(ctx context.Context, messages []ChatMessage, options ChatOptions) (string, error) {
	req := ChatRequest{
		Model:       a.model,
		Messages:    messages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Stream:      false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if options.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+options.APIKey)
	}

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API 请求失败: %s, %s", resp.Status, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", ErrInvalidResponse
	}

	return chatResp.Choices[0].Message.Content, nil
}

// StreamChatCompletion 流式聊天补全
func (a *OpenAIAdapter) StreamChatCompletion(ctx context.Context, messages []ChatMessage, options ChatOptions, callback StreamCallback) (string, error) {
	req := ChatRequest{
		Model:       a.model,
		Messages:    messages,
		Temperature: options.Temperature,
		MaxTokens:   options.MaxTokens,
		Stream:      true,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if options.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+options.APIKey)
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API 请求失败: %s, %s", resp.Status, string(body))
	}

	var fullContent strings.Builder

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var streamResp struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.NewDecoder(bytes.NewReader([]byte(data))).Decode(&streamResp); err != nil {
			continue
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			fullContent.WriteString(content)
			if err := callback(content); err != nil {
				return "", err
			}
		}
	}

	return fullContent.String(), nil
}

// ValidateConfig 验证配置
func (a *OpenAIAdapter) ValidateConfig(apiKey, baseURL string) error {
	// 简单的验证：检查 API Key 是否非空
	if apiKey == "" {
		return ErrInvalidAPIKey
	}

	// TODO: 可以添加实际的 API 测试请求
	return nil
}

// GetDefaultModel 获取默认模型
func (a *OpenAIAdapter) GetDefaultModel() string {
	return "gpt-3.5-turbo"
}

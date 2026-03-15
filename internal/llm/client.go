package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"research-ability-assessment/internal/config"
	"time"
)

type Client struct {
	config *config.LLMConfig
	client *http.Client
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewClient(cfg *config.LLMConfig) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *Client) Chat(ctx context.Context, messages []Message) (string, error) {
	log.Printf("LLM Client: 开始调用API - BaseURL=%s, Model=%s, API Key长度=%d",
		c.config.BaseURL, c.config.Model, len(c.config.APIKey))

	if c.config.APIKey == "" || c.config.APIKey == "${DEEPSEEK_API_KEY}" {
		log.Printf("LLM Client: 警告 - API Key未正确配置")
		return "", fmt.Errorf("API Key未配置")
	}

	reqBody := ChatRequest{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("LLM Client: 序列化请求失败: %v", err)
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	reqURL := c.config.BaseURL + "/chat/completions"
	log.Printf("LLM Client: 请求URL: %s", reqURL)

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(reqJSON))
	if err != nil {
		log.Printf("LLM Client: 创建请求失败: %v", err)
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	log.Printf("LLM Client: 发送请求...")
	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("LLM Client: API调用失败: %v", err)
		return "", fmt.Errorf("API调用失败: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("LLM Client: 收到响应 - 状态码: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Printf("LLM Client: API返回错误状态码: %d", resp.StatusCode)
		return "", fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		log.Printf("LLM Client: 解析API响应失败: %v", err)
		return "", fmt.Errorf("解析API响应失败: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		log.Printf("LLM Client: API响应中没有内容")
		return "", fmt.Errorf("API响应中没有内容")
	}

	log.Printf("LLM Client: API调用成功 - Choice数量=%d", len(chatResp.Choices))
	return chatResp.Choices[0].Message.Content, nil
}

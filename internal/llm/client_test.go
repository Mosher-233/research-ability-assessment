package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mosher-233/research-ability-assessment/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, responseBody string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/chat/completions", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-api-key", r.Header.Get("Authorization"))

		var req ChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		assert.NotEmpty(t, req.Model)
		assert.NotEmpty(t, req.Messages)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(responseBody))
	}))
}

func TestClient_Chat_Success(t *testing.T) {
	expectedContent := "{\"kbm_name\": \"文献检索策略\", \"level\": 3, \"credibility\": 0.85}"
	resp := ChatResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Model:   "deepseek-chat",
		Choices: []Choice{{Index: 0, Message: Message{Role: "assistant", Content: expectedContent}, FinishReason: "stop"}},
		Usage:   Usage{PromptTokens: 100, CompletionTokens: 50, TotalTokens: 150},
	}
	body, _ := json.Marshal(resp)

	srv := newTestServer(t, string(body), http.StatusOK)
	defer srv.Close()

	client := NewClient(&config.LLMConfig{
		Provider:  "deepseek",
		APIKey:    "test-api-key",
		BaseURL:   srv.URL,
		Model:     "deepseek-chat",
		MaxTokens: 2048,
	})

	content, err := client.Chat(context.Background(), []Message{
		{Role: "system", Content: "You are an expert."},
		{Role: "user", Content: "Classify this evidence."},
	})
	require.NoError(t, err)
	assert.Equal(t, expectedContent, content)
}

func TestClient_Chat_HTTPError(t *testing.T) {
	srv := newTestServer(t, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
	defer srv.Close()

	client := NewClient(&config.LLMConfig{
		APIKey:  "test-api-key",
		BaseURL: srv.URL,
		Model:   "deepseek-chat",
	})

	content, err := client.Chat(context.Background(), []Message{
		{Role: "user", Content: "test"},
	})
	assert.Error(t, err)
	assert.Empty(t, content)
	assert.Contains(t, err.Error(), "500")
}

func TestClient_Chat_EmptyChoices(t *testing.T) {
	resp := ChatResponse{
		ID:      "chatcmpl-456",
		Object:  "chat.completion",
		Model:   "deepseek-chat",
		Choices: []Choice{},
		Usage:   Usage{},
	}
	body, _ := json.Marshal(resp)

	srv := newTestServer(t, string(body), http.StatusOK)
	defer srv.Close()

	client := NewClient(&config.LLMConfig{
		APIKey:  "test-api-key",
		BaseURL: srv.URL,
		Model:   "deepseek-chat",
	})

	content, err := client.Chat(context.Background(), []Message{
		{Role: "user", Content: "test"},
	})
	assert.Error(t, err)
	assert.Empty(t, content)
	assert.Contains(t, err.Error(), "没有内容")
}

func TestClient_Chat_NoAPIKey(t *testing.T) {
	client := NewClient(&config.LLMConfig{
		APIKey:  "",
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})

	content, err := client.Chat(context.Background(), []Message{
		{Role: "user", Content: "test"},
	})
	assert.Error(t, err)
	assert.Empty(t, content)
	assert.Contains(t, err.Error(), "API Key未配置")
}

func TestClient_Chat_UnconfiguredEnvVar(t *testing.T) {
	client := NewClient(&config.LLMConfig{
		APIKey:  "${DEEPSEEK_API_KEY}",
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})

	content, err := client.Chat(context.Background(), []Message{
		{Role: "user", Content: "test"},
	})
	assert.Error(t, err)
	assert.Empty(t, content)
	assert.Contains(t, err.Error(), "API Key未配置")
}

func TestClient_Chat_InvalidJSON(t *testing.T) {
	srv := newTestServer(t, `not json at all {{{`, http.StatusOK)
	defer srv.Close()

	client := NewClient(&config.LLMConfig{
		APIKey:  "test-api-key",
		BaseURL: srv.URL,
		Model:   "deepseek-chat",
	})

	content, err := client.Chat(context.Background(), []Message{
		{Role: "user", Content: "test"},
	})
	assert.Error(t, err)
	assert.Empty(t, content)
}

func TestClient_Chat_MultipleChoices(t *testing.T) {
	resp := ChatResponse{
		ID:     "chatcmpl-789",
		Object: "chat.completion",
		Model:  "deepseek-chat",
		Choices: []Choice{
			{Index: 0, Message: Message{Role: "assistant", Content: "first choice"}, FinishReason: "stop"},
			{Index: 1, Message: Message{Role: "assistant", Content: "second choice"}, FinishReason: "length"},
		},
		Usage: Usage{PromptTokens: 50, CompletionTokens: 20, TotalTokens: 70},
	}
	body, _ := json.Marshal(resp)

	srv := newTestServer(t, string(body), http.StatusOK)
	defer srv.Close()

	client := NewClient(&config.LLMConfig{
		APIKey:  "test-api-key",
		BaseURL: srv.URL,
		Model:   "deepseek-chat",
	})

	content, err := client.Chat(context.Background(), []Message{
		{Role: "user", Content: "test"},
	})
	require.NoError(t, err)
	assert.Equal(t, "first choice", content) // should use first choice
}

func TestChatRequest_Serialization(t *testing.T) {
	req := ChatRequest{
		Model:    "deepseek-chat",
		Messages: []Message{{Role: "user", Content: "hello"}},
		MaxTokens: 2048,
	}

	data, err := json.Marshal(req)
	require.NoError(t, err)
	assert.Contains(t, string(data), "deepseek-chat")
	assert.Contains(t, string(data), "hello")
	assert.Contains(t, string(data), "2048")
}

func TestNewClient_Defaults(t *testing.T) {
	client := NewClient(&config.LLMConfig{
		APIKey:  "sk-test",
		BaseURL: "https://api.deepseek.com",
		Model:   "deepseek-chat",
	})

	assert.NotNil(t, client)
	assert.NotNil(t, client.config)
	assert.NotNil(t, client.client)
	assert.Equal(t, "sk-test", client.config.APIKey)
}

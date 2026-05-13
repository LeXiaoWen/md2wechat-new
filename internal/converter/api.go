package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lexiaowenn/md2wechat-new/internal/action"
	"go.uber.org/zap"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// APIRequest is an OpenAI-compatible chat completions request for API mode.
type APIRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type,omitempty"`
		Code    any    `json:"code,omitempty"`
	} `json:"error,omitempty"`
}

type textProviderConfig struct {
	Provider    string
	APIKey      string
	BaseURL     string
	Model       string
	Temperature float64
	Timeout     time.Duration
}

// apiConverter converts Markdown through an OpenAI-compatible chat completions API.
type apiConverter struct {
	log     *zap.Logger
	baseURL string
	timeout time.Duration
}

// NewAPIConverter creates an API converter with OpenAI-compatible defaults.
func NewAPIConverter(log *zap.Logger) *apiConverter {
	return &apiConverter{
		log:     log,
		baseURL: "https://api.openai.com/v1/chat/completions",
		timeout: 30 * time.Second,
	}
}

// NewAPIConverterWithURL creates an API converter with a specific base URL or endpoint.
func NewAPIConverterWithURL(log *zap.Logger, baseURL string) *apiConverter {
	return &apiConverter{
		log:     log,
		baseURL: normalizeChatCompletionsURL(baseURL),
		timeout: 30 * time.Second,
	}
}

// convertViaAPI executes conversion through the configured text provider.
func (c *converter) convertViaAPI(req *ConvertRequest) *ConvertResult {
	result := &ConvertResult{
		Mode:      ModeAPI,
		Theme:     req.Theme,
		Status:    action.StatusFailed,
		Action:    action.ActionConvert,
		Retryable: true,
		Success:   false,
	}

	provider, err := c.resolveTextProvider(req)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	apiConv := NewAPIConverterWithURL(c.log, provider.BaseURL)
	apiConv.SetTimeout(provider.Timeout)

	html, err := apiConv.Convert(&APIRequest{
		Model: provider.Model,
		Messages: []chatMessage{
			{Role: "system", Content: apiSystemPrompt()},
			{Role: "user", Content: c.buildAPIConversionPrompt(req)},
		},
		Temperature: provider.Temperature,
	}, provider.APIKey)
	if err != nil {
		result.Error = fmt.Sprintf("API call failed: %s", err.Error())
		c.log.Error("API conversion failed",
			zap.String("provider", provider.Provider),
			zap.String("model", provider.Model),
			zap.String("theme", req.Theme),
			zap.Error(err))
		return result
	}

	images := c.ExtractImages(req.Markdown)

	result.HTML = html
	result.Images = images
	result.Status = action.StatusCompleted
	result.Retryable = false
	result.Success = true

	c.log.Info("API conversion succeeded",
		zap.String("provider", provider.Provider),
		zap.String("model", provider.Model),
		zap.String("theme", req.Theme),
		zap.Int("image_count", len(images)))

	return result
}

func (c *converter) resolveTextProvider(req *ConvertRequest) (*textProviderConfig, error) {
	cfg := c.cfg
	provider := strings.ToLower(strings.TrimSpace(cfg.TextProvider))
	if provider == "" {
		provider = "openai"
	}

	baseURL := strings.TrimSpace(cfg.TextAPIBase)
	model := strings.TrimSpace(cfg.TextModel)
	switch provider {
	case "openai":
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		if model == "" {
			model = "gpt-4.1-mini"
		}
	case "deepseek":
		if baseURL == "" {
			baseURL = "https://api.deepseek.com"
		}
		if model == "" {
			model = "deepseek-chat"
		}
	case "siliconflow", "silicon", "sf":
		provider = "siliconflow"
		if baseURL == "" {
			baseURL = "https://api.siliconflow.cn/v1"
		}
		if model == "" {
			model = "Qwen/Qwen2.5-72B-Instruct"
		}
	case "custom":
	default:
		return nil, &ConvertError{Code: "INVALID_PROVIDER", Message: fmt.Sprintf("unsupported text provider: %s", provider)}
	}

	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		apiKey = strings.TrimSpace(cfg.TextAPIKey)
	}
	if apiKey == "" {
		return nil, ErrMissingAPIKey
	}
	if baseURL == "" {
		return nil, &ConvertError{Code: "MISSING_API_BASE", Message: "TEXT_API_BASE is required for API mode"}
	}
	if model == "" {
		return nil, &ConvertError{Code: "MISSING_MODEL", Message: "TEXT_MODEL is required for API mode"}
	}

	temperature := cfg.TextTemperature
	if temperature == 0 {
		temperature = 0.2
	}
	timeout := time.Duration(cfg.HTTPTimeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	return &textProviderConfig{
		Provider:    provider,
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Model:       model,
		Temperature: temperature,
		Timeout:     timeout,
	}, nil
}

func apiSystemPrompt() string {
	return strings.Join([]string{
		"You are a professional WeChat Official Account HTML typesetter.",
		"Return only HTML that can be pasted into WeChat editor.",
		"Use inline CSS only. Do not use scripts, external stylesheets, markdown fences, or explanations.",
	}, "\n")
}

func (c *converter) buildAPIConversionPrompt(req *ConvertRequest) string {
	metadata := req.Metadata
	return fmt.Sprintf(`Convert this Markdown article into WeChat Official Account compatible HTML.

Requirements:
- Return only the final HTML fragment.
- Use inline CSS on elements.
- Preserve image references from the Markdown as img tags or recognizable placeholders so the publishing pipeline can process them.
- Theme: %s
- Font size: %s
- Background: %s
- Title: %s
- Author: %s
- Digest: %s

Markdown:
%s`, req.Theme, req.FontSize, req.BackgroundType, metadata.Title, metadata.Author, metadata.Digest, req.Markdown)
}

// Convert calls an OpenAI-compatible chat completions API and returns normalized HTML.
func (a *apiConverter) Convert(req *APIRequest, apiKey string) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", a.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: a.timeout}
	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var apiResp chatCompletionResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("parse response: %w (body: %s)", err, string(body))
	}
	if apiResp.Error != nil {
		return "", &ConvertError{Code: "API_ERROR", Message: apiResp.Error.Message}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", &ConvertError{Code: "API_ERROR", Message: fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body))}
	}
	if len(apiResp.Choices) == 0 {
		return "", &ConvertError{Code: "EMPTY_RESPONSE", Message: "API returned no choices"}
	}

	html := normalizeProviderHTML(apiResp.Choices[0].Message.Content)
	if html == "" {
		return "", &ConvertError{Code: "EMPTY_RESPONSE", Message: "API returned empty content"}
	}
	return html, nil
}

func normalizeProviderHTML(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```") && strings.HasSuffix(content, "```") {
		lines := strings.Split(content, "\n")
		if len(lines) >= 2 {
			lines = lines[1 : len(lines)-1]
			content = strings.TrimSpace(strings.Join(lines, "\n"))
		}
	}
	return content
}

func normalizeChatCompletionsURL(baseURL string) string {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return "https://api.openai.com/v1/chat/completions"
	}
	baseURL = strings.TrimRight(baseURL, "/")
	if strings.HasSuffix(baseURL, "/chat/completions") {
		return baseURL
	}
	return baseURL + "/chat/completions"
}

// SetBaseURL sets API base URL or full chat completions endpoint for tests.
func (a *apiConverter) SetBaseURL(url string) {
	a.baseURL = normalizeChatCompletionsURL(url)
}

// SetTimeout sets request timeout.
func (a *apiConverter) SetTimeout(timeout time.Duration) {
	a.timeout = timeout
}

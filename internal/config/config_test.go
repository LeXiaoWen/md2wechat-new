package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadWithDefaultsPreservesCompressDefaultWhenOmitted(t *testing.T) {
	t.Setenv("WECHAT_APPID", "")
	t.Setenv("WECHAT_SECRET", "")
	t.Setenv("MD2WECHAT_API_KEY", "")
	t.Setenv("MD2WECHAT_BASE_URL", "")
	t.Setenv("CONVERT_MODE", "")
	t.Setenv("DEFAULT_THEME", "")
	t.Setenv("DEFAULT_BACKGROUND_TYPE", "")
	t.Setenv("IMAGE_API_KEY", "")
	t.Setenv("IMAGE_API_BASE", "")
	t.Setenv("IMAGE_PROVIDER", "")
	t.Setenv("IMAGE_MODEL", "")
	t.Setenv("IMAGE_SIZE", "")
	t.Setenv("COMPRESS_IMAGES", "")
	t.Setenv("MAX_IMAGE_WIDTH", "")
	t.Setenv("MAX_IMAGE_SIZE", "")
	t.Setenv("HTTP_TIMEOUT", "")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`
wechat:
  appid: appid
  secret: secret
api:
  convert_mode: api
`)
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadWithDefaults(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.CompressImages {
		t.Fatalf("expected CompressImages default to remain true when field is omitted")
	}
}

func TestLoadWithDefaultsRespectsExplicitCompressFalse(t *testing.T) {
	t.Setenv("WECHAT_APPID", "")
	t.Setenv("WECHAT_SECRET", "")
	t.Setenv("MD2WECHAT_API_KEY", "")
	t.Setenv("MD2WECHAT_BASE_URL", "")
	t.Setenv("CONVERT_MODE", "")
	t.Setenv("DEFAULT_THEME", "")
	t.Setenv("DEFAULT_BACKGROUND_TYPE", "")
	t.Setenv("IMAGE_API_KEY", "")
	t.Setenv("IMAGE_API_BASE", "")
	t.Setenv("IMAGE_PROVIDER", "")
	t.Setenv("IMAGE_MODEL", "")
	t.Setenv("IMAGE_SIZE", "")
	t.Setenv("COMPRESS_IMAGES", "")
	t.Setenv("MAX_IMAGE_WIDTH", "")
	t.Setenv("MAX_IMAGE_SIZE", "")
	t.Setenv("HTTP_TIMEOUT", "")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`
wechat:
  appid: appid
  secret: secret
api:
  convert_mode: api
image:
  compress: false
`)
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadWithDefaults(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.CompressImages {
		t.Fatalf("expected CompressImages to respect explicit false")
	}
}

func TestLoadWithDefaultsEnvOverridesFileCompressValue(t *testing.T) {
	t.Setenv("WECHAT_APPID", "")
	t.Setenv("WECHAT_SECRET", "")
	t.Setenv("MD2WECHAT_API_KEY", "")
	t.Setenv("MD2WECHAT_BASE_URL", "")
	t.Setenv("CONVERT_MODE", "")
	t.Setenv("DEFAULT_THEME", "")
	t.Setenv("DEFAULT_BACKGROUND_TYPE", "")
	t.Setenv("IMAGE_API_KEY", "")
	t.Setenv("IMAGE_API_BASE", "")
	t.Setenv("IMAGE_PROVIDER", "")
	t.Setenv("IMAGE_MODEL", "")
	t.Setenv("IMAGE_SIZE", "")
	t.Setenv("COMPRESS_IMAGES", "true")
	t.Setenv("MAX_IMAGE_WIDTH", "")
	t.Setenv("MAX_IMAGE_SIZE", "")
	t.Setenv("HTTP_TIMEOUT", "")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`
wechat:
  appid: appid
  secret: secret
api:
  convert_mode: api
image:
  compress: false
`)
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadWithDefaults(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.CompressImages {
		t.Fatalf("expected environment variable to override file value")
	}
}

func TestLoadWithDefaultsJSONUsesSameMergeRules(t *testing.T) {
	t.Setenv("WECHAT_APPID", "")
	t.Setenv("WECHAT_SECRET", "")
	t.Setenv("MD2WECHAT_API_KEY", "")
	t.Setenv("MD2WECHAT_BASE_URL", "")
	t.Setenv("CONVERT_MODE", "")
	t.Setenv("DEFAULT_THEME", "")
	t.Setenv("DEFAULT_BACKGROUND_TYPE", "")
	t.Setenv("IMAGE_API_KEY", "")
	t.Setenv("IMAGE_API_BASE", "")
	t.Setenv("IMAGE_PROVIDER", "")
	t.Setenv("IMAGE_MODEL", "")
	t.Setenv("IMAGE_SIZE", "")
	t.Setenv("COMPRESS_IMAGES", "")
	t.Setenv("MAX_IMAGE_WIDTH", "")
	t.Setenv("MAX_IMAGE_SIZE", "")
	t.Setenv("HTTP_TIMEOUT", "")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := []byte(`{
  "wechat": {
    "appid": "appid",
    "secret": "secret"
  },
  "api": {
    "convert_mode": "api"
  }
}`)
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadWithDefaults(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.CompressImages {
		t.Fatalf("expected JSON loader to preserve CompressImages default when field is omitted")
	}
}

func TestValidateForWeChatRequiresCredentials(t *testing.T) {
	cfg := &Config{}
	if err := cfg.ValidateForWeChat(); err == nil {
		t.Fatal("expected missing appid error")
	}

	cfg.WechatAppID = "appid"
	if err := cfg.ValidateForWeChat(); err == nil {
		t.Fatal("expected missing secret error")
	}
}

func TestValidateForImageGenerationRequiresImageKey(t *testing.T) {
	cfg := &Config{
		WechatAppID:  "appid",
		WechatSecret: "secret",
	}

	if err := cfg.ValidateForImageGeneration(); err == nil {
		t.Fatal("expected missing image key error")
	}

	cfg.ImageAPIKey = "image-key"
	if err := cfg.ValidateForImageGeneration(); err != nil {
		t.Fatalf("ValidateForImageGeneration() error = %v", err)
	}
}

func TestLoadWithDefaultsAppliesTextProviderDefaults(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		baseURL  string
		model    string
	}{
		{name: "openai", provider: "openai", baseURL: "https://api.openai.com/v1", model: "gpt-4.1-mini"},
		{name: "deepseek", provider: "deepseek", baseURL: "https://api.deepseek.com", model: "deepseek-chat"},
		{name: "siliconflow", provider: "silicon", baseURL: "https://api.siliconflow.cn/v1", model: "Qwen/Qwen2.5-72B-Instruct"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "config.yaml")
			content := "api:\n  text_provider: " + tt.provider + "\n"
			if err := os.WriteFile(path, []byte(content), 0600); err != nil {
				t.Fatalf("write config: %v", err)
			}

			cfg, err := LoadWithDefaults(path)
			if err != nil {
				t.Fatalf("LoadWithDefaults() error = %v", err)
			}
			if cfg.TextAPIBase != tt.baseURL {
				t.Fatalf("TextAPIBase = %q, want %q", cfg.TextAPIBase, tt.baseURL)
			}
			if cfg.TextModel != tt.model {
				t.Fatalf("TextModel = %q, want %q", cfg.TextModel, tt.model)
			}
		})
	}
}

func TestLoadWithDefaultsTextEnvOverridesFile(t *testing.T) {
	t.Setenv("TEXT_PROVIDER", "deepseek")
	t.Setenv("TEXT_API_KEY", "env-key")
	t.Setenv("TEXT_API_BASE", "https://env.example.com/v1")
	t.Setenv("TEXT_MODEL", "env-model")
	t.Setenv("TEXT_TEMPERATURE", "0.7")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := strings.TrimSpace(`
api:
  text_provider: openai
  text_key: file-key
  text_base_url: https://file.example.com/v1
  text_model: file-model
  text_temperature: 0.1
`)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadWithDefaults(path)
	if err != nil {
		t.Fatalf("LoadWithDefaults() error = %v", err)
	}
	if cfg.TextProvider != "deepseek" || cfg.TextAPIKey != "env-key" || cfg.TextAPIBase != "https://env.example.com/v1" || cfg.TextModel != "env-model" {
		t.Fatalf("loaded text config = %#v", cfg)
	}
	if cfg.TextTemperature != 0.7 {
		t.Fatalf("TextTemperature = %v", cfg.TextTemperature)
	}
}

func TestValidateForAPIConversionRequiresTextConfig(t *testing.T) {
	cfg := &Config{
		DefaultConvertMode: "api",
		TextProvider:       "custom",
	}
	if err := cfg.ValidateForAPIConversion(); err == nil || !strings.Contains(err.Error(), "TEXT_API_KEY") {
		t.Fatalf("ValidateForAPIConversion() error = %v", err)
	}

	cfg.TextAPIKey = "key"
	if err := cfg.ValidateForAPIConversion(); err == nil || !strings.Contains(err.Error(), "TEXT_API_BASE") {
		t.Fatalf("ValidateForAPIConversion() error = %v", err)
	}

	cfg.TextAPIBase = "https://api.example.com/v1"
	if err := cfg.ValidateForAPIConversion(); err == nil || !strings.Contains(err.Error(), "TEXT_MODEL") {
		t.Fatalf("ValidateForAPIConversion() error = %v", err)
	}

	cfg.TextModel = "model"
	if err := cfg.ValidateForAPIConversion(); err != nil {
		t.Fatalf("ValidateForAPIConversion() error = %v", err)
	}
}

func TestValidateCommonRejectsOutOfRangeValues(t *testing.T) {
	cfg := &Config{
		DefaultConvertMode: "invalid",
		MaxImageWidth:      1920,
		MaxImageSize:       5 * 1024 * 1024,
		HTTPTimeout:        30,
	}
	if err := cfg.validateCommon(); err == nil {
		t.Fatal("expected invalid convert mode error")
	}

	cfg.DefaultConvertMode = "api"
	cfg.MaxImageWidth = 10
	if err := cfg.validateCommon(); err == nil {
		t.Fatal("expected invalid max width error")
	}

	cfg.MaxImageWidth = 1920
	cfg.MaxImageSize = 10
	if err := cfg.validateCommon(); err == nil {
		t.Fatal("expected invalid max image size error")
	}

	cfg.MaxImageSize = 5 * 1024 * 1024
	cfg.HTTPTimeout = 0
	if err := cfg.validateCommon(); err == nil {
		t.Fatal("expected invalid http timeout error")
	}

	cfg.HTTPTimeout = 30
	cfg.TextProvider = "unknown"
	if err := cfg.validateCommon(); err == nil {
		t.Fatal("expected invalid text provider error")
	}

	cfg.TextProvider = "openai"
	cfg.TextTemperature = 3
	if err := cfg.validateCommon(); err == nil {
		t.Fatal("expected invalid text temperature error")
	}
}

func TestToMapMasksSecrets(t *testing.T) {
	cfg := &Config{
		WechatAppID:        "appid",
		WechatSecret:       "secret-value",
		MD2WechatAPIKey:    "api-key-value",
		TextAPIKey:         "text-key-value",
		ImageAPIKey:        "image-key-value",
		CompressImages:     true,
		MaxImageWidth:      1920,
		MaxImageSize:       5 * 1024 * 1024,
		HTTPTimeout:        30,
		configFile:         "/tmp/config.yaml",
		DefaultTheme:       "default",
		DefaultConvertMode: "api",
	}

	result := cfg.ToMap(true)
	if result["wechat_secret"] == "secret-value" || result["md2wechat_api_key"] == "api-key-value" || result["text_api_key"] == "text-key-value" || result["image_api_key"] == "image-key-value" {
		t.Fatalf("expected secrets to be masked: %#v", result)
	}
}

func TestSaveConfigAndLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	cfg := &Config{
		WechatAppID:           "appid",
		WechatSecret:          "secret",
		MD2WechatAPIKey:       "api-key",
		MD2WechatBaseURL:      "https://example.com",
		TextProvider:          "custom",
		TextAPIKey:            "text-key",
		TextAPIBase:           "https://text.example.com/v1",
		TextModel:             "text-model",
		TextTemperature:       0.4,
		DefaultConvertMode:    "api",
		DefaultTheme:          "default",
		DefaultBackgroundType: "none",
		ImageProvider:         "openai",
		ImageAPIKey:           "image-key",
		ImageAPIBase:          "https://api.example.com",
		ImageModel:            "model",
		ImageSize:             "1024x1024",
		CompressImages:        false,
		MaxImageWidth:         1600,
		MaxImageSize:          3 * 1024 * 1024,
		HTTPTimeout:           45,
	}

	if err := SaveConfig(path, cfg); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	loaded, err := LoadWithDefaults(path)
	if err != nil {
		t.Fatalf("LoadWithDefaults() error = %v", err)
	}
	if loaded.WechatAppID != "appid" || loaded.ImageAPIKey != "image-key" || loaded.TextAPIKey != "text-key" || loaded.TextModel != "text-model" || loaded.CompressImages != false {
		t.Fatalf("loaded config = %#v", loaded)
	}
	if loaded.TextTemperature != 0.4 {
		t.Fatalf("TextTemperature = %v", loaded.TextTemperature)
	}
}

func TestLoadWithDefaultsAppliesVolcengineImageDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := strings.TrimSpace(`
api:
  image_provider: "volcengine"
`)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadWithDefaults(path)
	if err != nil {
		t.Fatalf("LoadWithDefaults() error = %v", err)
	}
	if cfg.ImageAPIBase != "https://ark.cn-beijing.volces.com/api/v3" {
		t.Fatalf("ImageAPIBase = %q", cfg.ImageAPIBase)
	}
	if cfg.ImageModel != "doubao-seedream-5-0-260128" {
		t.Fatalf("ImageModel = %q", cfg.ImageModel)
	}
	if cfg.ImageSize != "2K" {
		t.Fatalf("ImageSize = %q", cfg.ImageSize)
	}
}

func TestConfigErrorFormatting(t *testing.T) {
	err := (&ConfigError{
		Field:   "WechatSecret",
		Message: "missing",
		Hint:    "set it",
	}).Error()
	if !strings.Contains(err, "WechatSecret") || !strings.Contains(err, "set it") {
		t.Fatalf("ConfigError.Error() = %q", err)
	}
}

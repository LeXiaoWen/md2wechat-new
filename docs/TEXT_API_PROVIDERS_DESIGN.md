# Text API Provider Design

## Scope

`md2wechat-new` keeps the existing `--mode ai` behavior unchanged. API mode is rebuilt as a configurable text-generation API mode for mainstream OpenAI-compatible providers.

The old `md2wechat.cn` conversion API is not supported in this fork. API mode must not call `/api/convert`, must not send the old `markdown/theme/fontSize/backgroundType` JSON payload, and must not require `MD2WECHAT_API_KEY`.

## Provider Contract

API mode uses an OpenAI-compatible chat completions contract:

- Method: `POST`
- Endpoint: `{text_base_url}/chat/completions`
- Auth: `Authorization: Bearer {text_key}`
- Request body:
  - `model`
  - `messages`
  - `temperature`
  - optional future generation parameters
- Response body:
  - `choices[0].message.content`

The converter builds a WeChat-safe HTML conversion prompt, sends it to the selected provider, extracts the returned HTML, and continues through the existing publish pipeline.

## Configuration

Configuration file keys live under `api`:

```yaml
api:
  convert_mode: "api"
  text_provider: "deepseek"
  text_key: "your-api-key"
  text_base_url: "https://api.deepseek.com"
  text_model: "deepseek-chat"
  text_temperature: 0.2
  default_theme: "default"
  background_type: "none"
  http_timeout: 30
```

Environment variables override config file values:

```bash
TEXT_PROVIDER=deepseek
TEXT_API_KEY=your-api-key
TEXT_API_BASE=https://api.deepseek.com
TEXT_MODEL=deepseek-chat
TEXT_TEMPERATURE=0.2
```

The CLI flag `--api-key` remains as a one-off override for `text_key`, but its help text should describe provider API keys rather than md2wechat.cn keys.

## Built-In Providers

| Provider | Aliases | Default base URL | Default model |
|---|---|---|---|
| `openai` | none | `https://api.openai.com/v1` | `gpt-4.1-mini` |
| `deepseek` | none | `https://api.deepseek.com` | `deepseek-chat` |
| `siliconflow` | `silicon`, `sf` | `https://api.siliconflow.cn/v1` | `Qwen/Qwen2.5-72B-Instruct` |
| `custom` | none | none | none |

Provider names are normalized to lowercase. `custom` requires both `text_base_url` and `text_model`.

## Validation

API mode is ready when:

- `text_key` is configured or supplied through `--api-key`
- provider is known
- resolved base URL is non-empty
- resolved model is non-empty
- `text_temperature`, when configured, is in a conservative range such as `0` to `2`

AI mode remains ready without these fields.

## Prompt Requirements

The API conversion prompt must instruct the model to:

- Convert the given Markdown into WeChat Official Account compatible HTML.
- Use inline CSS only.
- Avoid external stylesheets, scripts, markdown fences, and explanatory text.
- Preserve image references so the existing asset pipeline can process local, remote, and AI-generated images.
- Respect article metadata where useful.
- Apply the selected theme, font size, and background intent as style guidance.

The first implementation can keep the prompt deterministic and compact. Provider-specific prompt tuning should be avoided unless a provider requires compatibility work.

## HTML Extraction

Provider output should be normalized by:

- Trimming surrounding whitespace.
- Removing a single surrounding fenced code block when present.
- Accepting content that contains HTML tags.
- Returning a clear conversion error when the response is empty or no assistant content exists.

## Implementation Targets

- `internal/config/config.go`
  - Add text provider fields to `Config` and `configFile.API`.
  - Load environment variables.
  - Apply provider defaults.
  - Save and show the new fields.

- `internal/converter/api.go`
  - Replace the md2wechat.cn client with an OpenAI-compatible chat completions client.
  - Build prompt from Markdown, metadata, theme, font size, and background type.
  - Parse chat completions responses.

- `internal/converter/converter.go`
  - Validate provider-aware API mode configuration.
  - Continue to leave AI mode untouched.

- `cmd/md2wechat/convert.go`
  - Update `--api-key` help text.
  - Update missing-key error messaging.

- `internal/inspect/inspect.go`
  - Update API readiness checks and suggested fixes.

- Tests
  - Config loading/defaults for `openai`, `deepseek`, `siliconflow`, and `custom`.
  - API converter request shape and response parsing with `httptest`.
  - CLI validation for missing provider key.
  - Inspect readiness messaging.

## Non-Goals

- No compatibility layer for the old md2wechat.cn conversion API.
- No provider-specific SDK dependency.
- No streaming support in the first implementation.
- No changes to image generation providers in this step.

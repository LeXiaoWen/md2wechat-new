# 配置指南

`md2wechat-new` 的配置由三层组成：

1. 环境变量，优先级最高
2. 配置文件
3. 默认值

默认推荐使用全局配置文件：

```text
~/.config/md2wechat-new/config.yaml
```

## 初始化配置

```bash
md2wechat config init
```

如果你在源码目录里工作：

```bash
./md2wechat config init
```

如果文件已存在，命令会拒绝覆盖。直接编辑现有文件即可：

```bash
open ~/.config/md2wechat-new/config.yaml
```

## 查看生效配置

```bash
md2wechat config show --format json
```

如果需要确认密钥是否真的读到，可以临时显示 secret：

```bash
md2wechat config show --format json --show-secret
```

注意不要把带 secret 的输出发到公开渠道。

## 配置文件搜索顺序

当前 CLI 会按顺序查找：

1. `~/.config/md2wechat-new/config.yaml`
2. `~/.md2wechat.yaml`
3. `./md2wechat-new.yaml`

推荐新项目只使用第一项。旧路径保留是为了兼容历史配置。

## 文本 API 配置

API 模式使用 OpenAI-compatible `/chat/completions` 协议。你可以配置 OpenAI、DeepSeek、硅基流动或自定义中转。

### DeepSeek

```yaml
api:
  convert_mode: "api"
  default_theme: "default"
  text_provider: "deepseek"
  text_key: "sk-..."
  text_base_url: "https://api.deepseek.com"
  text_model: "deepseek-chat"
  text_temperature: 0.2
```

### 硅基流动

```yaml
api:
  convert_mode: "api"
  default_theme: "default"
  text_provider: "siliconflow"
  text_key: "sk-..."
  text_base_url: "https://api.siliconflow.cn/v1"
  text_model: "Qwen/Qwen2.5-72B-Instruct"
  text_temperature: 0.2
```

`text_model` 必须使用你在硅基流动控制台有权限调用的模型 ID。

### OpenAI

```yaml
api:
  convert_mode: "api"
  default_theme: "default"
  text_provider: "openai"
  text_key: "sk-..."
  text_base_url: "https://api.openai.com/v1"
  text_model: "gpt-4.1-mini"
  text_temperature: 0.2
```

### 自定义 OpenAI-compatible API

```yaml
api:
  convert_mode: "api"
  default_theme: "default"
  text_provider: "custom"
  text_key: "your-api-key"
  text_base_url: "https://your-api.example.com/v1"
  text_model: "your-model"
  text_temperature: 0.2
```

## 环境变量覆盖

环境变量优先级高于配置文件，适合临时切换或 CI。

```bash
export TEXT_PROVIDER="deepseek"
export TEXT_API_KEY="sk-..."
export TEXT_API_BASE="https://api.deepseek.com"
export TEXT_MODEL="deepseek-chat"
export TEXT_TEMPERATURE="0.2"
```

一次性运行：

```bash
TEXT_PROVIDER=deepseek \
TEXT_API_KEY=sk-... \
TEXT_API_BASE=https://api.deepseek.com \
TEXT_MODEL=deepseek-chat \
md2wechat convert article.md --preview
```

## 微信公众号配置

上传图片、创建草稿、上传已有 HTML 都需要：

```yaml
wechat:
  appid: "你的微信公众号 AppID"
  secret: "你的微信公众号 AppSecret"
```

还需要在微信公众号后台把当前机器公网 IP 加入白名单。

详见 [WECHAT-CREDENTIALS.md](WECHAT-CREDENTIALS.md)。

## 图片生成配置

图片功能使用另一组配置：

```yaml
api:
  image_provider: "openai"
  image_key: "sk-..."
  image_base_url: "https://api.openai.com/v1"
  image_model: "gpt-image-1.5"
  image_size: "1024x1024"
```

查看可用图片 provider：

```bash
md2wechat providers list --json
```

常用环境变量：

```bash
export IMAGE_PROVIDER="openai"
export IMAGE_API_KEY="sk-..."
export IMAGE_API_BASE="https://api.openai.com/v1"
export IMAGE_MODEL="gpt-image-1.5"
export IMAGE_SIZE="1024x1024"
```

## 常见配置问题

### `TEXT_API_KEY is required for API mode`

当前命令走的是 API 模式，但没有读到文本 API Key。

排查：

```bash
md2wechat config show --format json
```

确认：

- `config_file` 不是空
- `text_api_key` 不是空
- 环境变量没有覆盖成空值

### 硅基流动保存后不起作用

优先检查实际生效配置：

```bash
md2wechat config show --format json --show-secret
```

确认：

- `text_provider` 是 `siliconflow`
- `text_api_base` 是 `https://api.siliconflow.cn/v1`
- `text_model` 是控制台可用模型
- `text_api_key` 有值

### `context deadline exceeded`

请求 provider 超时。常见原因：

- provider 服务慢或网络不通
- 模型名不可用导致服务端迟迟不返回
- 本机代理或 DNS 问题

可以临时增大超时：

```yaml
api:
  http_timeout: 120
```

或者用环境变量：

```bash
HTTP_TIMEOUT=120 md2wechat convert article.md -o output.html
```

### 还看到 `https://www.md2wechat.cn`

新版 API 转换不再调用旧接口。若 `config show` 里仍看到旧字段，通常说明你调用的是旧二进制。请确认：

```bash
which md2wechat
md2wechat version --json
```

源码目录建议使用：

```bash
./md2wechat version --json
```

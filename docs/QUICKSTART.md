# 新手快速开始

这份指南帮你在 5 分钟内跑通 `md2wechat-new` 的主路径：配置文本 API，转换 Markdown，生成 HTML。

## 1. 安装

```bash
npm install -g @lexiaowen/md2wechat-new
md2wechat-new version --json
```

## 2. 初始化配置

```bash
md2wechat-new config init
```

默认配置文件：

```text
~/.config/md2wechat-new/config.yaml
```

如果命令提示文件已存在，说明已经初始化过，直接编辑这个文件即可。

## 3. 选择文本 API

API 模式默认使用 OpenAI-compatible Chat Completions 协议。

DeepSeek：

```yaml
api:
  convert_mode: "api"
  text_provider: "deepseek"
  text_key: "sk-..."
  text_base_url: "https://api.deepseek.com"
  text_model: "deepseek-chat"
  text_temperature: 0.2
```

硅基流动：

```yaml
api:
  convert_mode: "api"
  text_provider: "siliconflow"
  text_key: "sk-..."
  text_base_url: "https://api.siliconflow.cn/v1"
  text_model: "Qwen/Qwen2.5-72B-Instruct"
  text_temperature: 0.2
```

OpenAI：

```yaml
api:
  convert_mode: "api"
  text_provider: "openai"
  text_key: "sk-..."
  text_base_url: "https://api.openai.com/v1"
  text_model: "gpt-4.1-mini"
  text_temperature: 0.2
```

检查是否生效：

```bash
md2wechat-new config show --format json
```

重点看：

- `config_file`
- `text_provider`
- `text_api_base`
- `text_model`
- `text_api_key`

## 4. 转换 Markdown

准备 `article.md`：

```markdown
---
title: "文章标题"
author: "作者名"
digest: "文章摘要，最多 120 个字符"
---

# 文章标题

正文内容。
```

执行：

```bash
md2wechat-new inspect article.md --json
md2wechat-new preview article.md
md2wechat-new convert article.md --preview
md2wechat-new convert article.md -o output.html
```

安装到 PATH 后，等价通用写法：

```bash
md2wechat-new inspect article.md
md2wechat-new preview article.md
md2wechat-new convert article.md --preview
```

如果 `convert` 报 `TEXT_API_KEY is required for API mode`，说明 API 模式没有读到文本 API Key。请检查配置文件或临时使用：

```bash
TEXT_API_KEY="sk-..." md2wechat-new convert article.md --preview
```

## 5. 创建微信公众号草稿

配置微信凭证：

```yaml
wechat:
  appid: "你的微信公众号 AppID"
  secret: "你的微信公众号 AppSecret"
```

还需要在微信公众号后台配置当前公网 IP 白名单。详见 [WECHAT-CREDENTIALS.md](WECHAT-CREDENTIALS.md)。

转换并创建草稿：

```bash
md2wechat-new convert article.md --draft --cover cover.jpg
```

如果已经有转换好的 HTML：

```bash
md2wechat-new upload_html output.html --title "文章标题" --cover cover.jpg
```

## 下一步

- 需要完整命令说明：看 [USAGE.md](USAGE.md)
- 需要完整配置说明：看 [CONFIG.md](CONFIG.md)
- 需要从安装配置到草稿箱全流程：看 [LOCAL_TO_WECHAT_GUIDE.md](LOCAL_TO_WECHAT_GUIDE.md)
- 遇到微信凭证或 IP 白名单问题：看 [WECHAT-CREDENTIALS.md](WECHAT-CREDENTIALS.md)

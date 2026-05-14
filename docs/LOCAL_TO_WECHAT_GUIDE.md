# 本地初始化到微信公众号草稿全流程

这份指南从本地源码开始，走到把文章创建为微信公众号草稿。适合你现在这个项目的真实使用方式。

## 0. 准备材料

你需要：

- `article.md`
- `cover.jpg`
- 文本 API Key，OpenAI / DeepSeek / 硅基流动任选其一
- 微信公众号 AppID
- 微信公众号 AppSecret
- 已加入微信公众号后台 IP 白名单的当前公网 IP

微信凭证和 IP 白名单详见 [WECHAT-CREDENTIALS.md](WECHAT-CREDENTIALS.md)。

## 1. 构建当前项目

```bash
cd /Users/leo/Desktop/project/md2wechat-new
go build -o md2wechat ./cmd/md2wechat
./md2wechat version --json
```

后续命令使用 `./md2wechat`，确保调用的是当前项目代码。

## 2. 初始化配置

```bash
./md2wechat config init
```

默认文件：

```text
~/.config/md2wechat-new/config.yaml
```

如果提示文件已存在：

```bash
open ~/.config/md2wechat-new/config.yaml
```

## 3. 配置文本 API

选择一种 provider。

DeepSeek：

```yaml
api:
  convert_mode: "api"
  default_theme: "default"
  text_provider: "deepseek"
  text_key: "sk-..."
  text_base_url: "https://api.deepseek.com"
  text_model: "deepseek-chat"
  text_temperature: 0.2
  http_timeout: 60
```

硅基流动：

```yaml
api:
  convert_mode: "api"
  default_theme: "default"
  text_provider: "siliconflow"
  text_key: "sk-..."
  text_base_url: "https://api.siliconflow.cn/v1"
  text_model: "Qwen/Qwen2.5-72B-Instruct"
  text_temperature: 0.2
  http_timeout: 120
```

OpenAI：

```yaml
api:
  convert_mode: "api"
  default_theme: "default"
  text_provider: "openai"
  text_key: "sk-..."
  text_base_url: "https://api.openai.com/v1"
  text_model: "gpt-4.1-mini"
  text_temperature: 0.2
  http_timeout: 60
```

## 4. 配置微信公众号

同一个配置文件中填写：

```yaml
wechat:
  appid: "你的微信公众号 AppID"
  secret: "你的微信公众号 AppSecret"
```

然后到微信公众号后台配置 IP 白名单。

查询当前公网 IP：

```bash
curl ifconfig.me
```

## 5. 检查配置

```bash
./md2wechat config show --format json
```

需要确认：

- `config_file` 指向 `~/.config/md2wechat-new/config.yaml`
- `text_provider` 是你选择的 provider
- `text_api_key` 有值
- `text_api_base` 正确
- `text_model` 正确
- `wechat_appid` 有值
- `wechat_secret` 有值

如需确认 secret：

```bash
./md2wechat config show --format json --show-secret
```

## 6. 准备文章

推荐写 frontmatter：

```markdown
---
title: "文章标题"
author: "作者名"
digest: "文章摘要，最多 120 个字符"
---

# 文章标题

正文内容。

![配图](./images/demo.jpg)
```

## 7. 发布前检查

```bash
./md2wechat inspect article.md --json
```

先处理检查结果中的问题，比如：

- 标题过长
- 摘要过长
- 图片文件不存在
- API Key 未配置
- 微信凭证缺失
- 创建草稿缺少封面

## 8. 生成 HTML

只预览：

```bash
./md2wechat convert article.md --preview
```

保存 HTML：

```bash
./md2wechat convert article.md -o output.html
```

如果请求超时，可以提高 `http_timeout`，或临时切换 provider。

## 9. 上传图片并替换链接

如果文章中有本地图片：

```bash
./md2wechat convert article.md --upload -o output.html
```

这一步需要微信 AppID、Secret 和 IP 白名单。

## 10. 创建微信草稿

一步完成转换和创建草稿：

```bash
./md2wechat convert article.md --draft --cover cover.jpg
```

使用已有封面素材：

```bash
./md2wechat convert article.md --draft --cover-media-id MEDIA_ID
```

## 11. 已有 HTML 直接上传

如果你已经有 `output.html`，不想重新转换：

```bash
./md2wechat upload_html output.html --title "文章标题" --cover cover.jpg
```

完整参数示例：

```bash
./md2wechat upload_html output.html \
  --title "文章标题" \
  --author "作者名" \
  --digest "文章摘要" \
  --content-source-url "https://example.com/original" \
  --cover cover.jpg
```

## 12. 成功后去哪里看

命令成功后，登录微信公众号后台：

```text
内容与互动 -> 草稿箱
```

找到刚创建的草稿，人工检查排版、封面和摘要，再手动群发或继续编辑。

## 常见问题

### `TEXT_API_KEY is required for API mode`

API 模式没有读到文本 API Key。运行：

```bash
./md2wechat config show --format json
```

确认 `text_api_key` 是否为空。

### `context deadline exceeded`

provider 请求超时。可以：

- 增大 `api.http_timeout`
- 换一个模型
- 换 provider
- 检查代理和网络

### 微信返回 IP 白名单错误

当前公网 IP 没有加到微信公众号后台。重新查询：

```bash
curl ifconfig.me
```

然后更新微信后台 IP 白名单。

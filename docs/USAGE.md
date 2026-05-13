# 使用指南

这份文档按真实工作流组织：检查配置，转换文章，预览 HTML，上传图片，创建草稿，上传已有 HTML。

## 基础命令

```bash
md2wechat --help
md2wechat version --json
md2wechat capabilities --json
```

源码开发时使用：

```bash
./md2wechat --help
```

## 转换前检查

```bash
md2wechat inspect article.md
md2wechat inspect article.md --json
```

`inspect` 会解析：

- 标题、作者、摘要
- frontmatter 来源
- 本地图片是否存在
- 当前配置是否满足 API 模式
- 创建草稿是否还缺少微信配置或封面

建议每次发布前先跑 `inspect`。

## Markdown 转 HTML

预览，不上传图片：

```bash
md2wechat convert article.md --preview
```

保存 HTML：

```bash
md2wechat convert article.md -o output.html
```

指定主题：

```bash
md2wechat convert article.md --theme elegant-blue -o output.html
```

覆盖元数据：

```bash
md2wechat convert article.md \
  --title "文章标题" \
  --author "作者" \
  --digest "摘要" \
  -o output.html
```

元数据优先级：

- 标题：`--title` -> `frontmatter.title` -> 正文第一个一级标题 -> `未命名文章`
- 作者：`--author` -> `frontmatter.author`
- 摘要：`--digest` -> `frontmatter.digest` -> `frontmatter.summary` -> `frontmatter.description`

限制：

- 标题最多 32 个字符
- 作者最多 16 个字符
- 摘要最多 128 个字符

## API 模式

默认模式：

```bash
md2wechat convert article.md --mode api -o output.html
```

需要配置文本 API Key：

```bash
export TEXT_PROVIDER="deepseek"
export TEXT_API_KEY="sk-..."
export TEXT_API_BASE="https://api.deepseek.com"
export TEXT_MODEL="deepseek-chat"
```

也可以命令行临时传 key：

```bash
md2wechat convert article.md --api-key "sk-..." -o output.html
```

API 模式会调用：

```text
{TEXT_API_BASE}/chat/completions
```

因此 `TEXT_API_BASE` 通常应包含 `/v1`，但 DeepSeek 官方地址是 `https://api.deepseek.com`。

## AI 模式

AI 模式不直接调用 provider，而是输出给外部 AI/Agent 使用的结构化内容。

```bash
md2wechat convert article.md --mode ai --theme autumn-warm --json
```

适合让 Claude、Codex 或其他 Agent 接管排版生成。

## 本地预览

生成独立预览页：

```bash
md2wechat preview article.md
```

或直接用 convert 预览：

```bash
md2wechat convert article.md --preview
```

预览不会上传图片，也不会创建草稿。

## 图片处理

Markdown 中可以使用标准图片：

```markdown
![说明](./images/demo.jpg)
![说明](https://example.com/demo.jpg)
```

上传图片并替换 HTML 链接：

```bash
md2wechat convert article.md --upload -o output.html
```

上传单张图片到微信永久素材：

```bash
md2wechat upload_image cover.jpg
```

下载远程图片并上传：

```bash
md2wechat download_and_upload https://example.com/image.jpg
```

## 创建微信草稿

转换并创建草稿：

```bash
md2wechat convert article.md --draft --cover cover.jpg
```

使用已有封面素材：

```bash
md2wechat convert article.md --draft --cover-media-id MEDIA_ID
```

保存 draft JSON，不上传：

```bash
md2wechat convert article.md --save-draft draft.json
```

从 JSON 创建草稿：

```bash
md2wechat create_draft draft.json
```

## 上传已有 HTML

如果你已经有转换好的 HTML，不想再次调用文本 API：

```bash
md2wechat upload_html output.html --title "文章标题" --cover cover.jpg
```

可选参数：

```bash
md2wechat upload_html output.html \
  --title "文章标题" \
  --author "作者" \
  --digest "摘要" \
  --content-source-url "https://example.com/original" \
  --cover cover.jpg
```

使用已有封面素材：

```bash
md2wechat upload_html output.html \
  --title "文章标题" \
  --cover-media-id MEDIA_ID
```

## AI 图片

生成封面：

```bash
md2wechat generate_cover --article article.md --size 1024x1024
```

生成信息图：

```bash
md2wechat generate_infographic --article article.md
```

直接生成图片并上传微信：

```bash
md2wechat generate_image "A clean editorial cover about AI writing"
```

查看图片 provider：

```bash
md2wechat providers list --json
```

## 写作辅助

按风格生成文章：

```bash
md2wechat write --style dan-koe --topic "AI 时代的个人品牌"
```

AI 去痕：

```bash
md2wechat humanize article.md
```

## 发现命令

这些命令适合 Agent 和脚本调用：

```bash
md2wechat capabilities --json
md2wechat themes list --json
md2wechat providers list --json
md2wechat prompts list --json
md2wechat layout list --json
```

## 常见流程

### 只要 HTML

```bash
md2wechat inspect article.md --json
md2wechat convert article.md -o output.html
```

### 先转换，再人工检查，再上传

```bash
md2wechat convert article.md -o output.html
open output.html
md2wechat upload_html output.html --title "文章标题" --cover cover.jpg
```

### 一步创建草稿

```bash
md2wechat inspect article.md --json
md2wechat convert article.md --draft --cover cover.jpg
```

### 临时切换 provider

```bash
TEXT_PROVIDER=siliconflow \
TEXT_API_KEY=sk-... \
TEXT_API_BASE=https://api.siliconflow.cn/v1 \
TEXT_MODEL=Qwen/Qwen2.5-72B-Instruct \
md2wechat convert article.md -o output.html
```

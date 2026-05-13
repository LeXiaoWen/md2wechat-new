# md2wechat-new

`md2wechat-new` 是一个面向微信公众号创作和发布的命令行工具。它把 Markdown 转成适合微信公众号编辑器的 HTML，支持本地预览、图片上传、草稿创建、已有 HTML 上传，以及通过可配置的 OpenAI-compatible 文本 API 进行排版转换。

这个分支保留原项目的 AI 模式，同时把旧的固定 API 服务替换为可配置 provider：

- OpenAI
- DeepSeek
- 硅基流动 SiliconFlow
- 任意 OpenAI-compatible API

## 功能概览

| 能力 | 说明 |
| --- | --- |
| Markdown 转微信 HTML | `convert` 支持 API 模式和 AI 模式 |
| 可配置文本 API | 通过配置文件或环境变量切换 OpenAI、DeepSeek、SiliconFlow、custom |
| 本地预览 | `preview` 或 `convert --preview`，不上传图片、不创建草稿 |
| 微信草稿箱 | `convert --draft` 可转换后直接创建草稿 |
| 已有 HTML 上传 | `upload_html` 可把转换好的 HTML 直接创建为微信草稿 |
| 图片素材上传 | `upload_image`、`download_and_upload`、`--upload` |
| AI 图片 | `generate_cover`、`generate_infographic`、`generate_image` |
| 写作辅助 | `write`、`humanize` |
| Agent 友好 | `--json`、`capabilities`、`themes`、`providers`、`layout` |

## 快速开始

### 1. 本地构建

```bash
cd /Users/leo/Desktop/project/md2wechat-new
go build -o md2wechat ./cmd/md2wechat
./md2wechat version --json
```

如果你已经发布 npm 包，也可以安装：

```bash
npm install -g @lexiaowen/md2wechat-new
md2wechat version --json
```

npm 包安装时会从 GitHub Release 下载二进制。正式给别人使用前，需要先发布对应版本的 GitHub Release 资产。

发布 Release 后，也可以用固定版本安装脚本：

```bash
curl -fsSL https://github.com/LeXiaoWen/md2wechat-new/releases/download/v2.2.1/install.sh | bash
```

macOS 用户如果维护了 Homebrew tap，也可以安装：

```bash
brew install lexiaowenn/tap/md2wechat-new
```

### 2. 初始化配置

```bash
./md2wechat config init
```

默认生成：

```text
~/.config/md2wechat-new/config.yaml
```

如果文件已存在，直接编辑它即可。

### 3. 配置文本 API

DeepSeek 示例：

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

硅基流动示例：

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

OpenAI 示例：

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

### 4. 配置微信公众号

如果要上传图片或创建草稿，还需要配置：

```yaml
wechat:
  appid: "你的微信公众号 AppID"
  secret: "你的微信公众号 AppSecret"
```

并在微信公众号后台把当前机器公网 IP 加入 IP 白名单。详细说明见 [docs/WECHAT-CREDENTIALS.md](docs/WECHAT-CREDENTIALS.md)。

### 5. 转换文章

```bash
./md2wechat inspect article.md
./md2wechat convert article.md --preview
./md2wechat convert article.md -o output.html
```

安装到 PATH 后，通用命令写法是：

```bash
md2wechat inspect article.md
md2wechat preview article.md
md2wechat convert article.md --preview
```

创建微信草稿：

```bash
./md2wechat convert article.md --draft --cover cover.jpg
```

已有 HTML 直接上传：

```bash
./md2wechat upload_html output.html --title "文章标题" --cover cover.jpg
```

## 转换模式

### API 模式

默认模式。CLI 会调用你配置的 OpenAI-compatible `/chat/completions` 接口，把 Markdown 转成最终 HTML。

```bash
./md2wechat convert article.md --mode api -o output.html
```

特点：

- 需要 `TEXT_API_KEY` 或配置文件中的 `api.text_key`
- 支持 OpenAI、DeepSeek、SiliconFlow、custom
- 适合稳定、自动化、本地批处理
- 可直接配合 `--draft` 创建微信草稿

### AI 模式

AI 模式不会调用文本 API，而是生成可交给外部 AI/Agent 的排版提示和结构化输出。

```bash
./md2wechat convert article.md --mode ai --theme autumn-warm --json
```

适合你想让 Claude、Codex 或其他 Agent 接管 HTML 生成的场景。

## 常用命令

```bash
# 查看当前配置
./md2wechat config show --format json

# 校验配置
./md2wechat config validate

# 检查文章 metadata 和发布风险
./md2wechat inspect article.md --json

# 本地预览
./md2wechat preview article.md

# 转换为 HTML
./md2wechat convert article.md -o output.html

# 上传图片并替换 HTML 图片链接
./md2wechat convert article.md --upload -o output.html

# 转换并创建微信草稿
./md2wechat convert article.md --draft --cover cover.jpg

# 已有 HTML 创建草稿
./md2wechat upload_html output.html --title "标题" --author "作者" --digest "摘要" --cover cover.jpg

# 上传单张图片到微信永久素材
./md2wechat upload_image cover.jpg

# 发现可用主题、图片 provider、能力
./md2wechat themes list --json
./md2wechat providers list --json
./md2wechat capabilities --json
```

## 配置优先级

从高到低：

1. 环境变量
2. 配置文件
3. 默认值

默认配置文件搜索顺序：

1. `~/.config/md2wechat-new/config.yaml`
2. `~/.md2wechat.yaml`
3. `./md2wechat-new.yaml`

常用环境变量：

```bash
export TEXT_PROVIDER="deepseek"
export TEXT_API_KEY="sk-..."
export TEXT_API_BASE="https://api.deepseek.com"
export TEXT_MODEL="deepseek-chat"
export TEXT_TEMPERATURE="0.2"

export WECHAT_APPID="..."
export WECHAT_SECRET="..."
```

## 文档入口

- [文档索引](docs/README.md)
- [新手快速开始](docs/QUICKSTART.md)
- [安装指南](docs/INSTALL.md)
- [配置指南](docs/CONFIG.md)
- [完整使用指南](docs/USAGE.md)
- [本地初始化到微信公众号草稿全流程](docs/LOCAL_TO_WECHAT_GUIDE.md)
- [微信凭证与 IP 白名单](docs/WECHAT-CREDENTIALS.md)
- [高级排版模块](docs/LAYOUT.md)
- [故障排查](docs/TROUBLESHOOTING.md)

## 发布 npm 包

当前 npm 包名：

```text
@lexiaowen/md2wechat-new
```

本地发布：

```bash
npm login
npm publish --access public --otp=你的6位验证码
```

发布后安装：

```bash
npm install -g @lexiaowen/md2wechat-new
```

注意：npm 包的 `postinstall` 会下载 GitHub Release 里的二进制文件。发布 npm 前，应先推送 GitHub 仓库并创建同版本 Release。

## License

See [LICENSE](LICENSE).

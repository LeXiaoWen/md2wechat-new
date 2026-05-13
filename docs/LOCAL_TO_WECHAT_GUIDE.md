# 本地初始化到微信公众号草稿全流程

这份指南从本地源码开始，走到把 Markdown 文章上传到微信公众号草稿箱。  
当前项目的 API 模式使用 OpenAI-compatible 文本 API，支持 OpenAI、DeepSeek、硅基流动和自定义中转；不再使用旧的 `md2wechat.cn` 转换接口。

## 0. 准备材料

你需要准备：

- 一篇 Markdown 文章，例如 `article.md`
- 一个封面图片，例如 `cover.jpg`
- 文本 API Key：OpenAI / DeepSeek / 硅基流动任选其一
- 微信公众号 `AppID` 和 `AppSecret`
- 已把当前电脑公网 IP 加入微信公众号 IP 白名单

微信公众号相关配置参考 [微信凭证与 IP 白名单指南](WECHAT-CREDENTIALS.md)。

## 1. 本地构建 CLI

进入项目目录：

```bash
cd /Users/leo/Desktop/project/md2wechat-new
```

构建本地二进制：

```bash
go build ./cmd/md2wechat
```

验证命令可用：

```bash
./md2wechat --help
./md2wechat version --json
```

后续命令建议先使用 `./md2wechat`，确保调用的是当前项目构建出来的新版本，而不是系统 PATH 里的旧版本。

## 2. 初始化配置

生成配置文件：

```bash
./md2wechat config init ~/.config/md2wechat-new/config.yaml
```

如果提示文件已存在，说明已经初始化过，直接编辑即可：

```bash
open ~/.config/md2wechat-new/config.yaml
```

也可以在终端查看：

```bash
sed -n '1,120p' ~/.config/md2wechat-new/config.yaml
```

## 3. 配置文本 API

### 硅基流动

```yaml
api:
  text_provider: siliconflow
  text_key: "你的硅基流动 API Key"
  text_base_url: "https://api.siliconflow.cn/v1"
  text_model: "你的模型名"
  text_temperature: 0.2
```

`text_model` 需要使用你在硅基流动控制台可用的模型 ID。

### DeepSeek

```yaml
api:
  text_provider: deepseek
  text_key: "你的 DeepSeek API Key"
  text_base_url: "https://api.deepseek.com"
  text_model: "deepseek-chat"
  text_temperature: 0.2
```

### OpenAI

```yaml
api:
  text_provider: openai
  text_key: "你的 OpenAI API Key"
  text_base_url: "https://api.openai.com/v1"
  text_model: "gpt-4.1-mini"
  text_temperature: 0.2
```

### 自定义 OpenAI-compatible API

```yaml
api:
  text_provider: custom
  text_key: "你的 API Key"
  text_base_url: "https://your-api.example.com/v1"
  text_model: "your-model"
  text_temperature: 0.2
```

## 4. 配置微信公众号

同一个配置文件里填写：

```yaml
wechat:
  appid: "你的微信公众号 AppID"
  secret: "你的微信公众号 AppSecret"
```

创建草稿、上传图片、创建图片消息都需要这两个字段。

还需要在微信公众号后台把当前电脑公网 IP 加入白名单。查询公网 IP：

```bash
curl ifconfig.me
```

## 5. 检查配置是否生效

运行：

```bash
./md2wechat config show --format json
```

重点看这些字段：

```json
{
  "config_file": "/Users/leo/.config/md2wechat-new/config.yaml",
  "text_provider": "siliconflow",
  "text_api_key": "***",
  "text_api_base": "https://api.siliconflow.cn/v1",
  "text_model": "你的模型名",
  "wechat_appid": "...",
  "wechat_secret": "***"
}
```

如果 `config_file` 是空字符串，说明没有读到配置文件。请确认你使用的是当前项目的 `./md2wechat`，并确认配置文件在：

```text
~/.config/md2wechat-new/config.yaml
```

## 6. 准备 Markdown 文章

推荐在文章开头写 frontmatter：

```markdown
---
title: "文章标题"
author: "作者名"
digest: "这是一段 128 字以内的摘要"
---

# 文章标题

正文内容...

![配图](./images/demo.jpg)
```

元数据优先级：

- 标题：`--title` -> `frontmatter.title` -> 正文首个一级标题 -> `未命名文章`
- 作者：`--author` -> `frontmatter.author`
- 摘要：`--digest` -> `frontmatter.digest` -> `frontmatter.summary` -> `frontmatter.description`

限制：

- 标题最多 32 个字符
- 作者最多 16 个字符
- 摘要最多 128 个字符

## 7. 先做发布前检查

```bash
./md2wechat inspect article.md --json
```

重点检查：

- 标题、作者、摘要是否符合预期
- 本地图片是否存在
- API 模式是否缺少 `TEXT_API_KEY`
- 创建草稿是否缺少封面
- 微信配置是否缺失

## 8. 生成本地预览

生成独立预览文件：

```bash
./md2wechat preview article.md --json
```

或者直接输出转换 HTML：

```bash
./md2wechat convert article.md --preview
```

如果你想保存 HTML：

```bash
./md2wechat convert article.md -o output.html
```

## 9. 上传图片并生成 HTML

如果文章里有本地图片或远程图片，需要上传到微信素材并替换链接：

```bash
./md2wechat convert article.md --upload -o output.html
```

如果遇到微信 IP 白名单错误，先按第 4 步重新确认公网 IP 是否已加入微信公众号后台白名单。

## 10. 创建微信公众号草稿

创建草稿必须提供封面。推荐使用本地封面图片：

```bash
./md2wechat convert article.md --draft --cover cover.jpg
```

如果你已经有微信永久素材的封面 `media_id`，也可以：

```bash
./md2wechat convert article.md --draft --cover-media-id MEDIA_ID
```

两者不能同时使用。

创建成功后，命令会返回草稿相关信息。之后登录微信公众号后台，在草稿箱中检查内容，再手动发布。

## 11. 常用完整命令序列

```bash
cd /Users/leo/Desktop/project/md2wechat-new

go build ./cmd/md2wechat

./md2wechat config show --format json

./md2wechat inspect article.md --json

./md2wechat preview article.md --json

./md2wechat convert article.md --upload -o output.html

./md2wechat convert article.md --draft --cover cover.jpg
```

## 12. 临时环境变量用法

不想修改配置文件时，可以临时覆盖文本 API：

```bash
TEXT_PROVIDER=siliconflow \
TEXT_API_KEY="你的硅基流动 API Key" \
TEXT_API_BASE="https://api.siliconflow.cn/v1" \
TEXT_MODEL="你的模型名" \
./md2wechat convert article.md --preview
```

微信公众号凭证也可以临时传：

```bash
WECHAT_APPID="你的 AppID" \
WECHAT_SECRET="你的 AppSecret" \
./md2wechat convert article.md --draft --cover cover.jpg
```

## 13. 常见问题

### `TEXT_API_KEY is required for API mode`

说明 API 模式没有读到文本 API Key。检查：

```bash
./md2wechat config show --format json
```

确认：

- `config_file` 不是空
- `text_api_key` 不是空
- 当前运行的是 `./md2wechat`，不是系统旧版本

### 硅基流动配置保存了但不生效

优先检查 `config_file`：

```bash
./md2wechat config show --format json
```

必须指向：

```text
/Users/leo/.config/md2wechat-new/config.yaml
```

如果指向旧目录或为空，请重新构建当前项目：

```bash
go build ./cmd/md2wechat
```

再使用：

```bash
./md2wechat config show --format json
```

### 微信返回 IP 白名单错误

查询公网 IP：

```bash
curl ifconfig.me
```

然后到微信公众号后台的开发接口管理里，把这个 IP 加入白名单。

### 只想不用 API Key 试一下

AI 模式不直接调用文本 API，会输出 prompt/request：

```bash
./md2wechat convert article.md --mode ai --theme autumn-warm --json
```

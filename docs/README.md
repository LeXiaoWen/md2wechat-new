# 文档索引

这里保存 `md2wechat-new` 的使用、配置、发布和专题文档。建议先读入口文档，再按需要进入专题。

## 推荐阅读顺序

1. [新手快速开始](QUICKSTART.md)
2. [安装指南](INSTALL.md)
3. [配置指南](CONFIG.md)
4. [完整使用指南](USAGE.md)
5. [本地初始化到微信公众号草稿全流程](LOCAL_TO_WECHAT_GUIDE.md)
6. [微信凭证与 IP 白名单](WECHAT-CREDENTIALS.md)
7. [故障排查](TROUBLESHOOTING.md)

## 核心文档

| 文档 | 内容 |
| --- | --- |
| [QUICKSTART.md](QUICKSTART.md) | 5 分钟跑通本地转换 |
| [INSTALL.md](INSTALL.md) | 本地构建、npm、GitHub Release、PATH |
| [CONFIG.md](CONFIG.md) | 配置文件、环境变量、文本 API、微信凭证 |
| [USAGE.md](USAGE.md) | 常用命令、转换模式、上传 HTML、图片处理 |
| [LOCAL_TO_WECHAT_GUIDE.md](LOCAL_TO_WECHAT_GUIDE.md) | 从本地源码到微信草稿箱的完整操作 |
| [WECHAT-CREDENTIALS.md](WECHAT-CREDENTIALS.md) | AppID、AppSecret、IP 白名单 |
| [FAQ.md](FAQ.md) | 常见问题 |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | 报错排查 |

## 能力专题

| 文档 | 内容 |
| --- | --- |
| [LAYOUT.md](LAYOUT.md) | `:::block` 高级排版模块 |
| [DISCOVERY.md](DISCOVERY.md) | Agent 能力发现命令 |
| [IMAGE_PROVISIONERS.md](IMAGE_PROVISIONERS.md) | 图片生成 provider 配置 |
| [HUMANIZE.md](HUMANIZE.md) | AI 写作去痕 |
| [WRITING_FAQ.md](WRITING_FAQ.md) | 写作功能问答 |
| [TEXT_API_PROVIDERS_DESIGN.md](TEXT_API_PROVIDERS_DESIGN.md) | 可配置文本 API 设计 |

## 平台和集成

| 文档 | 内容 |
| --- | --- |
| [AGENT-GUIDE.md](AGENT-GUIDE.md) | Coding Agent 使用建议 |
| [OPENCLAW.md](OPENCLAW.md) | OpenClaw 集成 |
| [OBSIDIAN.md](OBSIDIAN.md) | Obsidian / Claudian 使用 |
| [SKILL-RULE.md](SKILL-RULE.md) | Skill 规则 |

## 开发和发布

| 文档 | 内容 |
| --- | --- |
| [ARCHITECTURE.md](ARCHITECTURE.md) | 架构说明 |
| [DESIGN.md](DESIGN.md) | 设计原则 |
| [SMOKE.md](SMOKE.md) | 烟雾测试记录 |
| [examples/config.yaml.example](examples/config.yaml.example) | 示例配置 |

## 最常用命令

```bash
md2wechat config init
md2wechat config show --format json
md2wechat inspect article.md --json
md2wechat convert article.md --preview
md2wechat convert article.md -o output.html
md2wechat convert article.md --draft --cover cover.jpg
md2wechat upload_html output.html --title "文章标题" --cover cover.jpg
```

# 安装指南

`md2wechat-new` 推荐通过 npm 或 GitHub Release 安装。

## npm 安装

包名：

```text
@lexiaowen/md2wechat-new
```

安装：

```bash
npm install -g @lexiaowen/md2wechat-new
```

验证：

```bash
md2wechat-new version --json
```

注意：npm 包不直接内置所有平台二进制，而是在 `postinstall` 阶段从 GitHub Release 下载对应平台文件。因此必须先发布同版本 GitHub Release。

## 从 npm 发布

发布前检查：

```bash
npm pack --json --dry-run
bash scripts/release-check.sh
```

发布：

```bash
npm login
npm publish --access public --otp=你的6位验证码
```

如果 npm 账号开启了 2FA，必须带 `--otp`，或使用允许 bypass 2FA 的 granular access token。

发布后安装：

```bash
npm install -g @lexiaowen/md2wechat-new --registry=https://registry.npmjs.org/
```

如果需要同步 npmmirror：

```bash
npx cnpm sync @lexiaowen/md2wechat-new
```

## GitHub Release 安装

Release 资产由 GitHub Actions 在推送 tag 时生成。版本来源是 `VERSION` 文件。

固定版本安装脚本：

```bash
curl -fsSL https://github.com/LeXiaoWen/md2wechat-new/releases/download/v0.1.0/install.sh | bash
```

macOS Homebrew：

```bash
brew install lexiaowenn/tap/md2wechat-new
```

本地确认：

```bash
cat VERSION
node -p "require('./package.json').version"
```

两者必须一致。

推送 tag：

```bash
git tag v0.1.0
git push origin main --tags
```

Actions 会生成：

- `md2wechat-darwin-amd64`
- `md2wechat-darwin-arm64`
- `md2wechat-linux-amd64`
- `md2wechat-linux-arm64`
- `md2wechat-windows-amd64.exe`
- `checksums.txt`
- 安装脚本和 OpenClaw 相关资产

## GitHub Actions 自动发布 npm

仓库需要配置 secret：

```text
NPM_TOKEN
```

如果还维护 Homebrew tap，还需要：

```text
HOMEBREW_TAP_GITHUB_TOKEN
```

tag workflow 会在 Release 创建后执行：

```bash
npm publish --access public
```

npm scope 是 `@lexiaowen`，GitHub Release 仓库是 `LeXiaoWen/md2wechat-new`。这两者不必相同。

## 常见安装问题

### `Scope not found`

说明 npm 包名的 scope 不属于当前 npm 账号。当前正确包名是：

```text
@lexiaowen/md2wechat-new
```

### `Two-factor authentication ... is required`

带一次性验证码发布：

```bash
npm publish --access public --otp=123456
```

### npm 安装后提示 binary missing

通常是 GitHub Release 资产不存在，或 `postinstall` 下载失败。

检查：

```bash
md2wechat-new version --json
```

确认当前版本对应的 Release 存在：

```text
https://github.com/LeXiaoWen/md2wechat-new/releases/tag/v0.1.0
```

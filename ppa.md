# PPA 发布文档

本项目已经包含 Debian 打包目录 `debian/`，以及用于 Launchpad PPA 发布的脚本和 GitHub Actions 工作流。

## 支持的 Ubuntu 系列

默认发布目标：

- Ubuntu 22.04: `jammy`
- Ubuntu 24.04: `noble`
- Ubuntu 26.04: `resolute`

在 GitHub Actions 里填写 `all` 时，会按 `jammy noble resolute` 依次构建和上传。

## 发布脚本

- `scripts/build-ppa-source.sh`：构建可上传 Launchpad 的 Debian 源码包。
- `scripts/upload-ppa.sh`：上传 `.changes` 文件到 Launchpad PPA。

构建脚本会根据目标 Ubuntu series 临时更新 `debian/changelog`，生成的 `.changes`、`.dsc`、源码包文件会放在仓库父目录。

## 本地手动构建和上传

安装依赖：

```bash
sudo apt update
sudo apt install -y devscripts debhelper dh-make dput-ng gnupg golang-any
```

准备 Go vendor 离线依赖目录：

```bash
go mod tidy
go mod vendor
```

确保签名 GPG key 已经导入本机，并设置签名 key：

```bash
export DEB_SIGN_KEYID="<gpg_key_fingerprint>"
```

构建并上传单个 Ubuntu series，例如发布 `0.0.2` 到 Ubuntu 22.04：

```bash
./scripts/build-ppa-source.sh jammy 0.0.2
./scripts/upload-ppa.sh <launchpad_user> <ppa_name> ../unimail-client_0.0.2ppa~jammy_source.changes
```

构建并上传多个 Ubuntu series：

```bash
version="0.0.2"
for series in jammy noble resolute; do
	./scripts/build-ppa-source.sh "${series}" "${version}"
	./scripts/upload-ppa.sh <launchpad_user> <ppa_name> "../unimail-client_${version}ppa~${series}_source.changes"
done
```

如果只是本地测试构建，不需要签名和上传，可以使用：

```bash
UNSIGNED=1 ./scripts/build-ppa-source.sh jammy 0.0.2
```

本地构建会修改 `debian/changelog`。如果只是手动补发 PPA，不准备把这次 changelog 写回仓库，发布完成后可以还原它：

```bash
git restore debian/changelog
```

## GitHub Actions 自动发布

工作流文件：`.github/workflows/ppa-publish.yaml`。

需要在 GitHub 仓库 Secrets 中配置：

- `PPA_GPG_PRIVATE_KEY`：base64 编码后的 GPG 私钥。
- `PPA_GPG_FINGERPRINT`：用于签名的密钥指纹。
- `PPA_GPG_PASSPHRASE`：私钥口令。
- `LAUNCHPAD_USER`：Launchpad 用户名。
- `PPA_NAME`：PPA 名称。

自动发布方式：推送 `v*` 标签，例如：

```bash
git tag v0.0.3
git push origin v0.0.3
```

tag 工作流会从标签名解析版本号，并默认发布到 `jammy`、`noble`、`resolute`。工作流会在 CI 工作区内更新 changelog 用于构建源码包，但不会再把 changelog commit 回仓库。

## GitHub Actions 手动补发

进入 GitHub 仓库的 Actions 页面，选择 `publish-ppa` 工作流，点击 `Run workflow`。

补发当前已有版本到全部默认 Ubuntu series：

- `upstream_version`: `0.0.2`
- `ubuntu_series`: `all`

只补发 Ubuntu 22.04 和 Ubuntu 26.04：

- `upstream_version`: `0.0.2`
- `ubuntu_series`: `jammy,resolute`

也可以用空格分隔：

```text
jammy resolute
```

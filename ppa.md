# ppa 文档

## PPA 发布流程

仓库已包含 Debian 打包目录 `debian/` 以及两个发布脚本：

- `scripts/build-ppa-source.sh`：构建可上传 Launchpad 的 Debian 源码包
- `scripts/upload-ppa.sh`：上传 `.changes` 文件到 PPA

本地手动发布示例：

```bash
sudo apt update
sudo apt install -y devscripts debhelper dh-make dput-ng gnupg golang-any

# 以 noble 为目标系列，构建 0.0.1 版本
./scripts/build-ppa-source.sh noble 0.0.1

# 上传到 Launchpad PPA
./scripts/upload-ppa.sh <launchpad_user> <ppa_name>
```

另外，仓库已新增自动发布工作流 `.github/workflows/ppa-publish.yaml`。

使用前请在 GitHub 仓库 Secrets 中配置：

- `PPA_GPG_PRIVATE_KEY`：base64 编码后的 GPG 私钥
- `PPA_GPG_FINGERPRINT`：用于签名的密钥指纹
- `PPA_GPG_PASSPHRASE`：私钥口令
- `LAUNCHPAD_USER`：Launchpad 用户名
- `PPA_NAME`：PPA 名称

触发方式：

- 推送标签（例如 `v0.0.1`）自动发布（默认目标 `noble`）
- 手动触发工作流并指定 `ubuntu_series`

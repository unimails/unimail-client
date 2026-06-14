# UNIMAIL-CLIENT

一个用基于unimail项目的邮件接口发送邮件的命令行工具, 快速上手使用。

[EnglishDocs](README.md)

## 目录

- [特性](#特性)
- [安装](#安装)
- [配置](#配置)
- [用法](#用法)
- [参数说明](#参数说明)
- [命令说明](#命令说明)
- [示例](#示例)
- [发布与更新](#发布与更新)
- [License](#license)

## 特性

- 支持从环境变量 `UNIMAIL_KEY` 读取 API key
- 支持通过 `-k/--key` 覆盖 key
- 支持发送纯文本邮件或 HTML 邮件
- 支持多个收件人、抄送、密送
- 支持附件上传，`--file` 使用 `;` 分隔成 name/path 对
- 支持 `version`、`update`、`upgrade` 子命令
- 支持 `-h/--help` 查看帮助信息

## 安装

### 方式一：从 Release 下载

当项目发布 Release 时，可以直接下载对应平台的压缩包并解压后使用。

### 方式二：本地构建

```bash
go build -o unimail-client .
```

## 配置

程序默认从环境变量 `UNIMAIL_KEY` 读取 API key：

```bash
export UNIMAIL_KEY=your_key
```

在 Windows PowerShell 中：

```powershell
$env:UNIMAIL_KEY = "your_key"
```

也可以在执行时通过 `-k` 或 `--key` 直接传入。

## 用法

### 查看帮助

```bash
unimail-client -h
unimail-client --help
```

### 发送邮件

```bash
unimail-client \
  -f "Sender Name" \
  -r a@example.com \
  -r b@example.com \
  -c cc@example.com \
  -b bcc@example.com \
  -s "Hello" \
  -t "Plain text content"
```

如果需要发送 HTML 内容，可以使用 `--html` 或 `-H`：

```bash
unimail-client \
  --form "Sender Name" \
  --receiver a@example.com \
  --subject "Hello" \
  --html "<b>Hello</b>"
```

### 附件

`--file` 每次传入一组 `name;path`，如果有多个附件，请重复传入多个 `--file`：

```bash
unimail-client \
  -r a@example.com \
  -s "With attachment" \
  -t "Please see attachment" \
  --file "report.pdf;./report.pdf"
```

多个附件请写成多次 `--file`：

```bash
--file "a.txt;./a.txt" \
--file "b.txt;./b.txt"
```

## 参数说明

| 参数                     | 说明                                                 |
| ------------------------ | ---------------------------------------------------- |
| `-k`, `--key`            | API key，默认从 `UNIMAIL_KEY` 读取                   |
| `-f`, `--form`, `--from` | 设置 `UnimailReq.From`                               |
| `-r`, `--receiver`       | 设置 `UnimailReq.Receivers`，可重复传入多次          |
| `-c`, `--cc`             | 设置 `UnimailReq.Cc`                                 |
| `-b`, `--bb`, `--bcc`    | 设置 `UnimailReq.Bcc`                                |
| `-s`, `--subject`        | 设置 `UnimailReq.Subject`                            |
| `-t`, `--txt`            | 设置 `UnimailReq.TxtContent`                         |
| `-H`, `--html`           | 设置 `UnimailReq.HtmlContent`                        |
| `--file`                 | 设置附件，格式为 `name;path`，多个附件请重复传入多次 |
| `-h`, `--help`           | 输出帮助信息                                         |

## 命令说明

### version

输出当前程序版本：

```bash
unimail-client version
```

### update / upgrade

检查远程仓库 `unimails/unimail-client` 的最新 Release。如果当前版本较旧，会自动执行更新：

```bash
unimail-client update
unimail-client upgrade
```

## 示例

### 发送纯文本邮件

```bash
unimail-client -r a@example.com -s "Test" -t "Hello"
```

### 发送 HTML 邮件

```bash
unimail-client --receiver a@example.com --subject "Test" --html "<h1>Hello</h1>"
```

### 发送带附件的邮件

```bash
unimail-client \
  -r a@example.com \
  -s "Report" \
  -t "Please check the attachment" \
  --file "report.pdf;./report.pdf"
```

## License

[BSD-3-Clause](LICENSE)

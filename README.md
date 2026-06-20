# UNIMAIL-CLIENT

A command-line tool for sending email through [unimail-go-sdk](https://github.com/unimails/unimail-go-sdk).

[Chinese Docs](README_zh.md)

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [UNIMAIL-CLIENT](#unimail-client)
  - [Features](#features)
  - [Installation](#installation)
    - [Option 1: Download from Release](#option-1-download-from-release)
    - [Option 2: Build locally](#option-2-build-locally)
    - [Option 3: Install from Ubuntu APT (PPA)](#option-3-install-from-ubuntu-apt-ppa)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Show help](#show-help)
    - [Send an email](#send-an-email)
    - [Attachments](#attachments)
  - [Flags](#flags)
  - [Commands](#commands)
    - [version](#version)
    - [update / upgrade](#update--upgrade)
  - [Examples](#examples)
    - [Send plain text email](#send-plain-text-email)
    - [Send HTML email](#send-html-email)
    - [Send email with attachments](#send-email-with-attachments)
  - [License](#license)

<!-- /code_chunk_output -->


## Features

- Reads the API key from the `UNIMAIL_KEY` environment variable by default
- Supports overriding the key with `-k/--key`
- Sends plain text or HTML email
- Supports multiple receivers, CC, and BCC
- Supports attachments via `--file` using `;` separated name/path pairs
- Supports `version`, `update`, and `upgrade` subcommands
- Supports `-h/--help` for usage information

## Installation

### Option 1: Download from Release

When the project publishes a Release, you can download the archive for your platform and unpack it locally.

### Option 2: Build locally

```bash
go build -o unimail-client .
```

### Option 3: Install from Ubuntu APT (PPA)

After this project is uploaded to your Launchpad PPA:

```bash
sudo add-apt-repository ppa:<launchpad_user>/<ppa_name>
sudo apt update
sudo apt install unimail-client
```

Replace `<launchpad_user>` and `<ppa_name>` with your actual PPA owner/name.

## Configuration

The program reads the API key from the `UNIMAIL_KEY` environment variable by default:

```bash
export UNIMAIL_KEY=your_key
```

On Windows PowerShell:

```powershell
$env:UNIMAIL_KEY = "your_key"
```

You can also pass the key directly with `-k` or `--key`.

## Usage

### Show help

```bash
unimail-client -h
unimail-client --help
```

### Send an email

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

To send HTML content, use `--html` or `-H`:

```bash
unimail-client \
	--form "Sender Name" \
	--receiver a@example.com \
	--subject "Hello" \
	--html "<b>Hello</b>"
```

### Attachments

`--file` accepts one `name;path` pair per flag. For multiple attachments, repeat `--file`:

```bash
unimail-client \
	-r a@example.com \
	-s "With attachment" \
	-t "Please see attachment" \
	--file "report.pdf;./report.pdf"
```

Multiple attachments should be passed as multiple `--file` flags:

```bash
--file "a.txt;./a.txt" \
--file "b.txt;./b.txt"
```

## Flags

| Flag                     | Description                                                                             |
| ------------------------ | --------------------------------------------------------------------------------------- |
| `-k`, `--key`            | API key, defaults to `UNIMAIL_KEY`                                                      |
| `-f`, `--form`, `--from` | Sets `UnimailReq.From`                                                                  |
| `-r`, `--receiver`       | Sets `UnimailReq.Receivers`; repeat multiple times                                      |
| `-c`, `--cc`             | Sets `UnimailReq.Cc`                                                                    |
| `-b`, `--bb`, `--bcc`    | Sets `UnimailReq.Bcc`                                                                   |
| `-s`, `--subject`        | Sets `UnimailReq.Subject`                                                               |
| `-t`, `--txt`            | Sets `UnimailReq.TxtContent`                                                            |
| `-H`, `--html`           | Sets `UnimailReq.HtmlContent`                                                           |
| `--file`                 | Attachment specification in `name;path` pairs; repeat the flag for multiple attachments |
| `-h`, `--help`           | Show help text                                                                          |

## Commands

### version

Print the current version:

```bash
unimail-client version
```

### update / upgrade

Check the latest Release of `unimails/unimail-client`. If the current version is older, the program will update itself:

```bash
unimail-client update
unimail-client upgrade
```

## Examples

### Send plain text email

```bash
unimail-client -r a@example.com -s "Test" -t "Hello"
```

### Send HTML email

```bash
unimail-client --receiver a@example.com --subject "Test" --html "<h1>Hello</h1>"
```

### Send email with attachments

```bash
unimail-client \
	-r a@example.com \
	-s "Report" \
	-t "Please check the attachment" \
	--file "report.pdf;./report.pdf"
```

## License

[BSD-3-Clause](LICENSE)

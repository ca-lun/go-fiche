# Go-Fiche

[![Build](https://github.com/ca-lun/go-fiche/actions/workflows/release.yml/badge.svg)](https://github.com/ca-lun/go-fiche/actions)
[![License](https://img.shields.io/github/license/ca-lun/go-fiche.svg)](LICENSE)

命令行 Pastebin 服务，用于分享终端输出。支持语法高亮。

## 功能特性

- 📝 通过 `nc` 命令快速上传文本
- 🎨 自动语法高亮（基于 highlight.js）
- 📋 一键复制按钮
- 🔗 Raw 原文查看

## 安装

### 从 Release 下载

前往 [Releases](https://github.com/ca-lun/go-fiche/releases) 下载对应平台的二进制文件。

### 从源码编译

```bash
git clone https://github.com/ca-lun/go-fiche.git
cd go-fiche
go build -o go-fiche .
```

## 使用方法

### 服务端启动

```bash
./go-fiche -d paste.example.com -p 9999 -o ./code -S -H -P 9989
```

### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-d` | 域名，用于返回 URL | `localhost` |
| `-p` | TCP 监听端口 | `9999` |
| `-o` | 存储目录 | `./code` |
| `-S` | 启用 HTTPS 前缀 | `false` |
| `-H` | 启用内置 HTTP 服务（语法高亮） | `false` |
| `-P` | HTTP 服务端口 | `9989` |
| `-B` | 缓冲区大小 (bytes) | `32768` |
| `-l` | 日志文件路径 | 无 |

### 客户端使用

```bash
# 上传命令输出
echo "Hello World" | nc paste.example.com 9999

# 上传文件
cat file.txt | nc paste.example.com 9999

# 上传剪贴板
xclip -o | nc paste.example.com 9999
```

### 访问方式

| URL | 说明 |
|-----|------|
| `https://paste.example.com/xxxxx` | 带语法高亮的页面 |
| `https://paste.example.com/raw/xxxxx` | 原始纯文本 |

## Nginx 配置示例

使用反向代理模式，将请求转发到 go-fiche 的 HTTP 服务：

```nginx
server {
    listen 80;
    listen [::]:80;
    server_name paste.example.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    listen [::]:443 ssl;
    http2 on;
    server_name paste.example.com;

    ssl_certificate     /path/to/fullchain.pem;
    ssl_certificate_key /path/to/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:9989;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Systemd 服务

创建 `/etc/systemd/system/go-fiche.service`：

```ini
[Unit]
Description=Go-Fiche Pastebin Service
After=network.target

[Service]
Type=simple
# 注意：监听 1024 以下端口需要 root 权限，不要设置 User
# User=www-data
ExecStart=/usr/local/bin/go-fiche -d paste.example.com -p 9999 -o /var/www/paste -S -H -P 9989
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

> **注意**：如果使用 1024 以下的端口，需要以 root 运行或使用 `CAP_NET_BIND_SERVICE` 能力。

启动服务：

```bash
sudo systemctl enable --now go-fiche
```

## License

MIT License
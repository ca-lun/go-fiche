# Go-Fiche

[![Build](https://github.com/ca-lun/go-fiche/actions/workflows/release.yml/badge.svg)](https://github.com/ca-lun/go-fiche/actions)
[![License](https://img.shields.io/github/license/ca-lun/go-fiche.svg)](LICENSE)

命令行 Pastebin 服务，用于分享终端输出。

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

### 服务端

```bash
./go-fiche -d paste.example.com -p 9999 -o ./code -H -P 9989
```

### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-d` | 域名，用于返回 URL | `localhost` |
| `-p` | TCP 监听端口 | `9999` |
| `-o` | 存储目录 | `./code` |
| `-S` | 启用 HTTPS 前缀 | `false` |
| `-H` | 启用内置 HTTP 服务 | `false` |
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

## Nginx 配置示例

### 基础反向代理

```nginx
server {
    listen 80;
    server_name paste.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    http2 on;
    server_name paste.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # 静态文件服务
    location / {
        root /path/to/code;
        default_type text/plain;
        charset utf-8;
        
        # 禁止目录列表
        autoindex off;
        
        # 缓存设置
        expires 1d;
        add_header Cache-Control "public, immutable";
    }
}
```

### 带语法高亮（可选）

如果需要语法高亮，可以使用 [highlight.js](https://highlightjs.org/) 或其他方案。

```nginx
server {
    listen 443 ssl;
    http2 on;
    server_name paste.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    root /path/to/code;

    location / {
        default_type text/plain;
        charset utf-8;
        autoindex off;
    }

    # 可选：添加 raw 路径返回纯文本
    location /raw/ {
        alias /path/to/code/;
        default_type text/plain;
        charset utf-8;
    }
}
```

### Systemd 服务

创建 `/etc/systemd/system/go-fiche.service`：

```ini
[Unit]
Description=Go-Fiche Pastebin Service
After=network.target

[Service]
Type=simple
User=www-data
ExecStart=/usr/local/bin/go-fiche -d paste.example.com -p 9999 -o /var/www/paste -S
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl enable --now go-fiche
```

## License

MIT License
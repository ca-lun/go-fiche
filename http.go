package main

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>%s - Paste</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css">
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { background: #0d1117; min-height: 100vh; }
        .toolbar {
            background: #161b22;
            padding: 8px 16px;
            border-bottom: 1px solid #30363d;
            display: flex;
            gap: 12px;
            align-items: center;
        }
        .toolbar a {
            color: #58a6ff;
            text-decoration: none;
            font-size: 14px;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
        }
        .toolbar a:hover { text-decoration: underline; }
        pre {
            margin: 0;
            padding: 16px;
            overflow-x: auto;
        }
        code {
            font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace;
            font-size: 14px;
            line-height: 1.5;
        }
        .hljs { background: #0d1117; }
    </style>
</head>
<body>
    <div class="toolbar">
        <a href="/raw/%s">Raw</a>
        <a href="javascript:navigator.clipboard.writeText(document.querySelector('code').textContent)">Copy</a>
    </div>
    <pre><code>%s</code></pre>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>
    <script>hljs.highlightAll();</script>
</body>
</html>`

func serveHTTP() {
	outputDir := viper.GetString("output")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// 禁止目录列表
		if path == "/" || strings.HasSuffix(path, "/") {
			http.Error(w, "Usage: echo 'content' | nc "+viper.GetString("domain")+" "+fmt.Sprint(viper.GetInt("port")), http.StatusOK)
			return
		}

		// 处理 /raw/xxx 路径
		if strings.HasPrefix(path, "/raw/") {
			slug := strings.TrimPrefix(path, "/raw/")
			filePath := filepath.Join(outputDir, filepath.Clean(slug))

			// 安全检查
			if !strings.HasPrefix(filePath, outputDir) {
				http.NotFound(w, r)
				return
			}

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			http.ServeFile(w, r, filePath)
			return
		}

		// 处理普通路径，返回带高亮的 HTML
		slug := strings.TrimPrefix(path, "/")
		filePath := filepath.Join(outputDir, filepath.Clean(slug))

		// 安全检查
		if !strings.HasPrefix(filePath, outputDir) {
			http.NotFound(w, r)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, htmlTemplate, html.EscapeString(slug), html.EscapeString(slug), html.EscapeString(string(content)))
	})

	port := viper.GetInt("httpport")
	log.Printf("Starting embedded http server on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("HTTP Server error: %s", err.Error())
	}
}

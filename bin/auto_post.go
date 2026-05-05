package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	fmt.Println("正在啟動 AI 技術文章生成流程...")

	// 1. 定義 Prompt
	prompt := `你現在是一位資深 AI/ML 技術部落客。請隨機挑選一個前沿的 AI 或機器學習技術主題（例如：LLM 代理、擴散模型優化、RAG 新進展、機器學習工程化等），寫一篇具備深度的教學或分享文章。
文章語言請以繁體中文為主，偶爾夾雜專業英文術語。內容要豐富、專業且有創意。
請嚴格輸出 Hugo Markdown 格式的內容，包含完整的 Frontmatter。
Frontmatter 格式參考：
---
title: "文章標題"
subtitle: "副標題"
date: 2026-05-05T00:00:00+08:00
lastmod: 2026-05-05T00:00:00+08:00
draft: false
author: "AI Assistant"
tags: ["AI", "Machine Learning"]
categories: ["Tech"]
toc:
  enable: true
---
<!--more-->
[文章正文]

重要規則：
1. 只輸出 Markdown 原始碼。
2. 不要包含 markdown 區塊標記（如 ` + "```" + `markdown）。
3. 文章日期請設定為今日。`

	// 2. 呼叫 Gemini CLI
	fmt.Println("正在請求 AI 生成內容 (這可能需要幾十秒)...")
	cmd := exec.Command("gemini", "-p", prompt)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("呼叫 Gemini CLI 失敗: %v\n", err)
		fmt.Printf("錯誤訊息: %s\n", stderr.String())
		return
	}

	content := stdout.String()
	// 清理可能的 markdown 標記 (有些 AI 還是會習慣加上去)
	content = strings.TrimPrefix(content, "```markdown")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	if content == "" {
		fmt.Println("錯誤：AI 產出的內容為空。")
		return
	}

	// 3. 儲存檔案
	now := time.Now()
	timestamp := now.Format("20060102-150405")
	fileName := fmt.Sprintf("content/posts/ai-gen-%s.md", timestamp)
	
	err = os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		fmt.Printf("寫入檔案失敗: %v\n", err)
		return
	}

	fmt.Printf("文章已成功生成： %s\n", fileName)

	// 4. 自動觸發部署
	fmt.Println("正在啟動自動部署流程 (cicd_script.go)...")
	deployCmd := exec.Command("go", "run", "bin/cicd_script.go")
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	
	err = deployCmd.Run()
	if err != nil {
		fmt.Printf("部署失敗: %v\n", err)
		return
	}

	fmt.Println("✅ AI 文章生成與部署流程全部完成！")
}

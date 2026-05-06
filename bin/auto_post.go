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
	prompt := `你現在是一位資深 AI/ML 技術專家。請使用已安裝的 'ai-tech-blogger' skill 寫作規範，隨機挑選一個前沿且具備實作價值的 AI/ML 主題（例如：LLM 代理、Agentic RAG、擴散模型底層優化、高效能 ML 工程實踐等）。

要求：
1. 嚴格遵守 'ai-tech-blogger' skill 的所有寫作指導方針。
2. 全文必須超過 1500 字，內容要詳實且具備專業深度。
3. 標題要極具吸引力，展現專家洞察力。
4. 內容必須包含具體的實作做法、方向、範例程式碼或架構圖 (Mermaid)。
5. 參考現有 posts 的格式豐富度。
6. 文章日期設定為今日。
7. Weight 必須設定為 2。
8. 只輸出 Markdown 原始碼，不含任何包圍的 code block 標記。`

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

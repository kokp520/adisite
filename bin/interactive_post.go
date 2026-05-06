package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("======= 🤖 互動式 AI 技術文章生成系統 =======")

	// 1. 主題發想與確認
	fmt.Println("\n[步驟 1] 正在請 AI 提供主題建議...")
	topicPrompt := "請根據 'ai-tech-blogger' skill 的規範，列出 3 個目前最前沿且適合寫成技術文章的 AI/ML 主題建議。每個建議請提供標題與簡短的內容大綱。請直接輸出主題列表，不要有額外廢話。"
	
	suggestions := callGemini(topicPrompt)
	fmt.Println("\nAI 建議的主題：")
	fmt.Println(suggestions)

	fmt.Print("\n請輸入您選擇的主題編號、自訂主題，或輸入 'q' 退出: ")
	scanner.Scan()
	userChoice := strings.TrimSpace(scanner.Text())
	if userChoice == "q" || userChoice == "" {
		fmt.Println("流程已終止。")
		return
	}

	// 2. 內容生成
	fmt.Printf("\n[步驟 2] 正在根據主題「%s」生成全文 (預計 1500 字以上)...\n", userChoice)
	articlePrompt := fmt.Sprintf("請使用 'ai-tech-blogger' skill 寫作規範，針對主題「%s」撰寫全文。要求：詳實具體、包含實作方向與範例程式碼、超過 1500 字。只輸出 Markdown 原始碼，不含包圍的 code block 標記。", userChoice)
	
	content := callGemini(articlePrompt)
	content = cleanMarkdown(content)

	// 3. 審核內容
	fmt.Println("\n[步驟 3] 內容已生成，請預覽（顯示前 500 字）：")
	fmt.Println("--------------------------------------------------")
	if len(content) > 500 {
		fmt.Println(content[:500] + "...")
	} else {
		fmt.Println(content)
	}
	fmt.Println("--------------------------------------------------")

	for {
		fmt.Print("\n您對內容滿意嗎？ (y: 儲存並部署 / r: 重新生成 / q: 放棄並退出): ")
		scanner.Scan()
		action := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if action == "q" {
			fmt.Println("已放棄本次生成。")
			return
		} else if action == "r" {
			fmt.Println("正在重新生成內容...")
			content = cleanMarkdown(callGemini(articlePrompt))
			fmt.Println("\n新內容預覽：")
			if len(content) > 500 {
				fmt.Println(content[:500] + "...")
			} else {
				fmt.Println(content)
			}
			continue
		} else if action == "y" {
			break
		}
	}

	// 4. 儲存並部署
	now := time.Now()
	timestamp := now.Format("20060102-150405")
	fileName := fmt.Sprintf("content/posts/ai-gen-%s.md", timestamp)
	
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		fmt.Printf("寫入檔案失敗: %v\n", err)
		return
	}
	fmt.Printf("\n✅ 文章已儲存至: %s\n", fileName)

	fmt.Println("\n[步驟 4] 啟動自動部署流程...")
	deployCmd := exec.Command("go", "run", "bin/cicd_script.go")
	deployCmd.Stdout = os.Stdout
	deployCmd.Stderr = os.Stderr
	if err := deployCmd.Run(); err != nil {
		fmt.Printf("部署失敗: %v\n", err)
		return
	}

	fmt.Println("\n🎉 全流程完成！您的新文章已發布。")
}

func callGemini(prompt string) string {
	cmd := exec.Command("gemini", "-p", prompt)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Sprintf("錯誤: %v\n%s", err, stderr.String())
	}
	return stdout.String()
}

func cleanMarkdown(content string) string {
	// 1. 尋找第一個 "---"，過濾掉前面所有的系統警告訊息 (如 MCP issues)
	startIdx := strings.Index(content, "---")
	if startIdx != -1 {
		content = content[startIdx:]
	}

	// 2. 清理可能的 markdown 區塊標記
	content = strings.TrimPrefix(content, "```markdown")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	return strings.TrimSpace(content)
}

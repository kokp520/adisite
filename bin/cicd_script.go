package main

// 這個檔案是一個用於自動化 Hugo 靜態網站生成和部署的腳本。主要功能包括：
// 1. 設定 Hugo 專案目錄的路徑。
// 2. 執行 Git 命令來添加、提交和推送網站內容的變更。
// 3. 執行 Hugo 命令來生成靜態網站。
// 4. 執行 Git 命令來添加、提交和推送生成的靜態網站文件。

// 腳本的主要流程如下：
// 1. 設定 Hugo 專案目錄和 Hugo 路徑。
// 2. 執行 Git 命令來添加、提交和推送網站內容的變更。
// 3. 執行 Hugo 命令來生成靜態網站。
// 4. 執行 Git 命令來添加、提交和推送生成的靜態網站文件。

// 腳本中使用了 executeCommands 和 executeCommand 函數來執行命令，並處理可能的錯誤。

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// 設定 Hugo 專案目錄的路徑
	hugoProjectPath := "/Users/adiwu/repo/kokp520.github.io"

	hugoPath := "/Users/adiwu/flutterApp/adisite/"

	// ----- website content -----

	// 執行 Git 命令
	if err := executeCommands(hugoPath, "git", "add", "."); err != nil {
		fmt.Println("Git add 失敗:", err)
		return
	}

	if err := executeCommands(hugoPath, "git", "commit", "-m", "update new post"); err != nil {
		fmt.Println("Git commit 失敗:", err)
		return
	}

	if err := executeCommands(hugoPath, "git", "push"); err != nil {
		fmt.Println("Git push 失敗:", err)
		return
	}

	fmt.Println("hugo path完成！")

	// ---- public -----

	// 執行 Hugo 命令
	if err := executeCommand(hugoPath, "hugo"); err != nil {
		fmt.Println("Hugo 執行失敗:", err)
		return
	}

	// 執行 Git 命令
	if err := executeCommands(hugoProjectPath, "git", "add", "."); err != nil {
		fmt.Println("Git add 失敗:", err)
		return
	}

	if err := executeCommands(hugoProjectPath, "git", "commit", "-m", "update"); err != nil {
		fmt.Println("Git commit 失敗:", err)
		return
	}

	if err := executeCommands(hugoProjectPath, "git", "push"); err != nil {
		fmt.Println("Git push 失敗:", err)
		return
	}

	fmt.Println("public完成！")
}

func executeCommand(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeCommands(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getAbsolutePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

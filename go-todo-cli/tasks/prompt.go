package tasks

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// ユーザー選択の関数
func GetUserChoice() int {
	fmt.Println("\n1: タスクを追加 | 2: タスク一覧表示 | 3: タスク削除 | 4: タスク完了 | 5: 未完了タスク表示 | 6: 完了タスク表示 | 7: タスク編集 | 8: CSVエクスポート | 0: 終了")
	fmt.Print("選択してください: ")

	var choice int
	fmt.Scan(&choice)

	return choice
}

// ユーザー入力を取得する汎用関数
func getInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	// 前後の余白を除去
	return strings.TrimSpace(input)
}

// 日時をパースする汎用関数
func parseDate(input string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04", input)
}

// タスク追加のプロンプト
func AddTaskPrompt() {
	taskName := getInput("追加するタスク名を入力してください: ")
	dueDateStr := getInput("タスクの期限を入力してください (YYYY-MM-DD HH:MM): ")

	dueDate, err := parseDate(dueDateStr)
	if err != nil {
		fmt.Println("無効な期限フォーマットです。", err)
		return
	}

	AddTask(taskName, dueDate)
}

// タスク編集のユーザー入力を処理する関数
func EditTaskPronmt() {
	idStr := getInput("編集するタスクのIDを入力してください。")
	var taskID int
	_, err := fmt.Sscanf(idStr, "%d", &taskID)
	if err != nil {
		fmt.Println("無効なIDです。")
		return
	}

	newName := getInput("新しいタスク名を入力してください。: ")
	dueDateStr := getInput("新しい期限を YYYY-MM-DD HH:MM 形式で入力してください: ")

	newDueDate, err := parseDate(dueDateStr)
	if err != nil {
		fmt.Println("無効な期限フォーマットです。")
		return
	}

	EditTask(taskID, newName, newDueDate)
}

// タスク削除のプロンプト
func DeleteTaskPrompt() {
	idStr := getInput("削除するタスクのIDを入力してください: ")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		fmt.Println("無効なIDです。")
		return
	}

	DeleteTask(id)
}

// タスク完了のプロンプト
func CompleteTaskPrompt() {
	idStr := getInput("削除するタスクのIDを入力してください: ")
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		fmt.Println("無効なIDです。")
		return
	}

	CompleteTask(id)
}

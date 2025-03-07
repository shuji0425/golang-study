package main

import (
	"fmt"
	"go-todo-cli/tasks"
)

func main() {
	// 起動時にファイルからタスクを読み込む
	tasks.LoadTasks()

	fmt.Println("TODOアプリへようこそ！")

	// ユーザー入力を受け付ける
	for {
		choice := tasks.GetUserChoice()

		switch choice {
		case 1:
			tasks.AddTaskPrompt()
		case 2:
			tasks.ListTasks("")
		case 3:
			tasks.DeleteTaskPrompt()
		case 4:
			tasks.CompleteTaskPrompt()
		case 5:
			tasks.ListTasks("未完了")
		case 6:
			tasks.ListTasks("完了")
		case 7:
			tasks.EditTaskPronmt()
		case 8:
			tasks.ExportToCSV()
		case 0:
			// 終了時はタスクファイルに保存
			tasks.SaveTasks()
			fmt.Println("アプリを終了します。")
			return
		default:
			fmt.Println("無効な選択です。")
		}
	}
}

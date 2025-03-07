package tasks

import (
	"fmt"
	"slices"
	"sort"
	"time"
)

// Task構造体（タスクの情報を管理）
type Task struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DueDate     time.Time `json:"dueDate"`     // 期限
	IsCompleted bool      `json:"isCompleted"` // 完了
}

// タスク一覧を保持するスライス（可変長配列）
var tasks []Task

// JSONファイル名
const taskFile = "tasks.json"

// 次に割り当てるタスクID
var nextID int

// タスクを追加する関数
func AddTask(taskName string, dueDate time.Time) {
	if taskName == "" {
		fmt.Println("タスクが空です。")
		return
	}

	// 新しいタスクを追加
	newTask := Task{ID: nextID, Name: taskName, DueDate: dueDate}
	tasks = append(tasks, newTask)

	// 次のIDを加算
	nextID++

	// ファイルに保存
	SaveTasks()
	fmt.Println("タスクを追加しました。")
}

// タスク一覧を表示する関数
func ListTasks(filter string) {
	if len(tasks) == 0 {
		fmt.Println("TODOリストは空です。")
		return
	}

	// タスクを期限順にソート
	SortTasks()

	fmt.Println("\n【TODOリスト】")
	for _, task := range tasks {
		// フィルタ処理
		if filter == "未完了" && task.IsCompleted {
			continue
		}
		if filter == "完了" && !task.IsCompleted {
			continue
		}

		// 完了状態
		completedStatus := "未完了"
		if task.IsCompleted {
			completedStatus = "完了"
		}

		// 期限が過ぎた場合に表示
		overdue := ""
		if task.DueDate.Before(time.Now()) {
			overdue = "(期限切れ)"
		}
		fmt.Printf("%d: %s - 期限: %s %s %s\n", task.ID, task.Name, task.DueDate.Format("2006-01-02 15:04:05"), overdue, completedStatus)
	}
}

// タスクを期限順に並び替える関数
func SortTasks() {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})
}

// タスクを編集する関数
func EditTask(taskID int, newName string, newDueDate time.Time) {
	for i, task := range tasks {
		if task.ID == taskID {
			tasks[i].Name = newName
			tasks[i].DueDate = newDueDate
			SaveTasks()
			fmt.Println("タスクを編集しました！")
			return
		}
	}
	fmt.Println("エラー: 指定されたIDのタスクが見つかりません。")
}

// タスクを削除する関数
func DeleteTask(id int) {
	// 指定されたIDのタスクを削除
	for i, task := range tasks {
		if task.ID == id {
			// タスク削除
			tasks = slices.Delete(tasks, i, i+1)
			SaveTasks()
			fmt.Printf("タスクID %d を削除しました。\n", id)
			return
		}
	}

	// IDが見つからないとき
	fmt.Println("指定されたタスクIDが見つかりません。")
}

// タスクを完了にする関数
func CompleteTask(id int) {
	for i, task := range tasks {
		if task.ID == id {
			// 完了状態にする
			tasks[i].IsCompleted = true
			SaveTasks()
			fmt.Printf("タスクID %d が完了しました。\n", id)
			return
		}
	}

	// IDが見つからないとき
	fmt.Println("指定されたタスクIDが見つかりません。")
}

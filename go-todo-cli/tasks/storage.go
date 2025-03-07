package tasks

import (
	"encoding/json"
	"fmt"
	"os"
)

// タスクをファイルから読み込む関数
func LoadTasks() {
	file, err := os.Open(taskFile)
	if err != nil {
		// ファイルが存在しない場合はエラーを無視
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("タスクの読み込みエラー: ", err)
		return
	}
	defer file.Close()

	// JSONファイルを読み込んで、tasksスライスにデコード
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println("タスクのデコードエラー:", err)
	}

	// 既存の最大IDを探し、nextIDを更新
	nextID = 0
	for _, task := range tasks {
		if task.ID > nextID {
			nextID = task.ID
		}
	}
	nextID++
}

// タスクをファイルに保存する関数
func SaveTasks() {
	file, err := os.Create(taskFile)
	if err != nil {
		fmt.Println("タスクの保存エラー:", err)
	}
	defer file.Close()

	// tasksスライスをJSON形式でエンコードしてファイルに書き込み
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // インデントを追加して見やすくする
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("タスクのエンコードエラー:", err)
	}
}

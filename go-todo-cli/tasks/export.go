package tasks

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// タスクをCSVにエクスポートする関数
func ExportToCSV() {
	file, err := os.Create("tasks.csv")
	if err != nil {
		fmt.Println("エラー: CSV ファイルを作成できませんでした。")
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	// データを書き込む
	defer writer.Flush()

	// ヘッダーデータを書き込む
	writer.Write([]string{"ID", "Name", "DueDate", "Completed"})

	// 書くタスクを書き込む
	for _, task := range tasks {
		writer.Write([]string{
			strconv.Itoa(task.ID),
			task.Name,
			task.DueDate.Format("2006-01-02 15:04"),
			strconv.FormatBool(task.IsCompleted),
		})
	}

	fmt.Println("タスクを CSV にエクスポートしました。")
}

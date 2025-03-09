package batch

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var errorMessages = []string{
	"Connection failed: Timeout reached",
	"Invalid request: Missing required parameter",
	"Database error: Could not connect to DB",
	"Timeout error",
	"Unauthorized access",
}

// ログを定期的に書き込む関数
func WriteLogs() {
	// ログファイルのパス
	logFile := "logs/batch_test.log"

	for {
		// ランダムなメッセージを選択
		errorMessage := errorMessages[rand.Intn(len(errorMessages))]

		// 現在時刻を取得
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		// ログメッセージを作成
		logEntry := fmt.Sprintf("[%s] ERROR %s\n", timestamp, errorMessage)

		// ログファイルに追記
		err := appendToFile(logFile, logEntry)
		if err != nil {
			fmt.Println("ログファイルへの書き込みに失敗:", err)
		} else {
			fmt.Println("ログを書き込み:", logEntry)
		}

		// 5秒待機
		time.Sleep(5 * time.Second)
	}
}

// ログファイルにデータを追記する関数
func appendToFile(filename, text string) error {
	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return err
	}

	// ファイルを開く
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}

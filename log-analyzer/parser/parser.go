package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// エラーメッセージを抽出しエラーの種類ごといカウントする関数
func ExtractErrorStats(filepath string) (map[string]int, map[string]int, error) {
	// エラータイプごとのカウントを格納するマップ
	errorStats := make(map[string]int) // エラーの種類ごと
	timeStats := make(map[string]int)  // 時間ごとのエラー回数

	// ログファイルを開く
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, fmt.Errorf("ファイルを開けません: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	errorRegex := regexp.MustCompile(`\[(\d{4}-\d{2}-\d{2} \d{2}):\d{2}:\d{2}\] ERROR (.+?):`)

	// ファイルを1行ずつ読み込む
	for scanner.Scan() {
		line := scanner.Text()
		match := errorRegex.FindStringSubmatch(line)
		if len(match) > 2 {
			timeHour := match[1]  // "2025-03-08 10" のような時間単位
			errorType := match[2] // "Connection failed" などのエラータイプ

			errorStats[errorType]++
			timeStats[timeHour+":00 - "+timeHour+":59"]++ // 1時間ごとにまとめる
		}
	}

	return errorStats, timeStats, scanner.Err()
}

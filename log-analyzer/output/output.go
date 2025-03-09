package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// エラーメッセージをCSV形式で保存する関数
func SaveResultsToCSV(filepath string, errorStats map[string]int) {
	// 新しいCSVファイルを作成
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("CSVファイルを作成できません: %v\n", err)
		return
	}
	defer file.Close()

	// CSVライターを作成
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込む
	err = writer.Write([]string{"Error", "Type", "Count"})
	if err != nil {
		fmt.Printf("CSVヘッダーの書き込みエラ-: %v\n", err)
		return
	}

	// 結果をCSVに書き込む
	for errorType, count := range errorStats {
		err = writer.Write([]string{errorType, fmt.Sprintf("%d", count)})
		if err != nil {
			fmt.Printf("CSVデータの書き込みエラー: %v\n", err)
			continue
		}
	}

	fmt.Printf("結果をCSVファイルに保存しました: %s\n", filepath)
}

// エラメッセージをJSON形式で保存する関数
func SaveResultsToJSON(filepath string, errorStats map[string]int) {
	// 新しいJSONファイル作成
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("JSONファイルを作成できません: %v\n", err)
		return
	}
	defer file.Close()

	// JSONエンコーダを作成
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// 結果をJSONとして書き込む
	err = encoder.Encode(errorStats)
	if err != nil {
		fmt.Printf("JSONデータの書き込みに失敗: %v\n", err)
		return
	}

	fmt.Printf("結果をJSONファイルに保存しました: %s\n", filepath)
}

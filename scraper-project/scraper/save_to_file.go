package scraper

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// スクレイピングしたデータの構造体
type ScrapedData struct {
	URL   string   `json:"url" csv:"url"`
	Title string   `json:"title" csv:"title"`
	Links []string `json:"links" csv:"links"`
}

// CSVに書き込む関数
func SaveToCSV(filename string, data []ScrapedData) error {
	// ファイルを作成または開く
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました: %v", err)
	}
	defer file.Close()

	// CSVライターを作成
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダを書き込む
	header := []string{"URL", "Title", "Links"}
	err = writer.Write(header)
	if err != nil {
		return fmt.Errorf("CSVヘッダの書き込みに失敗しました: %v", err)
	}

	// CSVファイルにデータを書き込む
	for _, record := range data {
		// Linksを文字列に変換
		links := fmt.Sprintf("%v", record.Links)
		recordArray := []string{record.URL, record.Title, links}
		err := writer.Write(recordArray)
		if err != nil {
			return fmt.Errorf("CSVの書き込みに失敗しました: %v", err)
		}
	}

	return nil
}

// JSONに書き込む関数
func SaveToJSON(filename string, data []ScrapedData) error {
	// ファイルを作成または開く
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("JSONの書き込みに失敗しました: %v", err)
	}

	return nil
}

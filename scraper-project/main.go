package main

import (
	"fmt"
	"scraper-project/scraper"
	"sync"
)

func main() {
	// スクレイピングするURLのリスト
	urls := []string{
		// "https://example.com",
		// "https://golang.org",
		"https://www.wikipedia.org/",
	}

	// 訪問済みURLを管理するマップ
	visited := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	fmt.Println("=== 再帰的並行クローリング開始 ===")

	// 待機リストに追加
	wg.Add(1)
	// 複数のURLから並行してタイトルを取得
	scraper.FetchTitles(urls, visited, &mu, 2, &wg)
	// 全てが終わるまで待機
	wg.Wait()

	fmt.Println("\n=== クローリング完了 ===")
}

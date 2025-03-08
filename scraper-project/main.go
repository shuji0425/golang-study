package main

import (
	"fmt"
	"scraper-project/scraper"
	"sync"
)

func main() {
	// スクレイピングするURLのリスト
	urls := []string{
		"https://scraping-for-beginner.herokuapp.com/ranking/",
		"https://scraping-training.vercel.app/site?postCount=20&title=%E3%81%93%E3%82%8C%E3%81%AF{no}%E3%81%AE%E8%A8%98%E4%BA%8B%E3%81%A7%E3%81%99&dateFormat=YYYY-MM-DD&isTime=true&timeFormat=&isImage=true&interval=360&isAgo=true&countPerPage=10&page=1&robots=true&",
		"https://scraping.okinan.com/",
	}

	// 訪問済みURLを管理するマップ
	visited := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	fmt.Println("=== 再帰的並行クローリング開始 ===")

	// 待機リストに追加
	wg.Add(1)
	// 複数のURLから並行してタイトルを取得
	go scraper.FetchTitles(urls, visited, &mu, 2, &wg)
	// 全てが終わるまで待機
	wg.Wait()

	fmt.Println("\n=== クローリング完了 ===")
}

package main

import (
	"scraper-project/scraper"
)

func main() {
	// スクレイピングするURLのリスト
	urls := []string{
		"https://example.com",
		"https://golang.org",
	}

	// 複数のURLから並行してタイトルを取得
	scraper.FetchTitles(urls)
}

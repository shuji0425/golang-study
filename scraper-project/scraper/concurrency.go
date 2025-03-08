package scraper

import (
	"fmt"
	"sync"
)

// 複数のURLから並行してタイトルを取得
func FetchTitles(urls []string) {
	var wg sync.WaitGroup

	// 各URLについてグルーチンで処理を並行実行
	for _, url := range urls {
		wg.Add(1) // ゴルーチンを待機リストに追加
		go func(url string) {
			defer wg.Done() // 処理完了で待機リストから削除

			// URLからHTMLを取得
			html, err := FetchHTML(url)
			if err != nil {
				fmt.Printf("URL(%s)の取得に失敗: %v\n", url, err)
				return
			}

			// HTMLからタイトルを解析
			title, err := ParseTitlle(html)
			if err != nil {
				fmt.Printf("タイトルの解析に失敗: %v\n", err)
			}

			// 結果を表示
			fmt.Printf("URL: %s, Title: %s\n", url, title)
		}(url)
	}

	// 全てのゴルーチンが終わるまで待機
	wg.Wait()
}

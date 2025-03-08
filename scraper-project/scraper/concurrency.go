package scraper

import (
	"fmt"
	"sync"
)

// 複数のURLから並行してタイトルを取得
func FetchTitles(urls []string, visited map[string]bool, mu *sync.Mutex, depth int, wg *sync.WaitGroup) {
	// 関数終了時にwg.Done()を呼び出す
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	// depthが0以下なら終了
	if depth <= 0 {
		return
	}

	// 再起で次に処理するリンクを一時的に保存
	newLinksSet := make(map[string]bool)
	var localWg sync.WaitGroup

	// 各URLについてグルーチンで処理を並行実行
	for _, url := range urls {
		localWg.Add(1) // ゴルーチンを待機リストに追加
		go func(url string) {
			defer localWg.Done() // 処理完了で待機リストから削除

			// 訪問済みチェック
			mu.Lock()
			if visited[url] {
				mu.Unlock()
				return
			}
			visited[url] = true
			mu.Unlock()

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
				return
			}

			// リンクを取得
			links, err := ExtraciLinks(html, url)
			if err != nil {
				fmt.Printf("リンクの取得に失敗: %v\n", err)
				return
			}

			// 結果を表示
			fmt.Printf("\n=== %s のスクレイピング結果 ===\n", url)
			fmt.Printf("タイトル: %s\n", title)

			// 次のクロール対象URLを追加
			mu.Lock()
			for _, link := range links {
				// すでに訪問していないURLのみ追加
				if !visited[link] && !newLinksSet[link] {
					newLinksSet[link] = true
				}
			}
			mu.Unlock()
		}(url)
	}

	// 全ての並行処理が終わるまで待つ
	localWg.Wait()

	// ローカル集合から次にスクロールするURLのスライスを作成
	newURLs := []string{}
	for link := range newLinksSet {
		newURLs = append(newURLs, link)
	}

	// 再帰的に次のページをクロール
	if len(newURLs) > 0 {
		// 再起前に親のwaitgroupを1つ進める
		wg.Add(1)
		go FetchTitles(newURLs, visited, mu, depth-1, wg)
	}
}

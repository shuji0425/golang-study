package main

import (
	"flag"
	"fmt"
	"log"
	"log-analyzer/api"
	"log-analyzer/output"
	"log-analyzer/parser"
	"log-analyzer/watcher"
	"net/http"
	"sync"
)

func main() {
	logDir := "logs"

	// 監視を開始
	go func() {
		err := watcher.WatchLogsDir(logDir, func(filepath string) {
			fmt.Println("解析開始:", filepath)
			_, _, err := parser.ExtractErrorStats(filepath)
			if err != nil {
				log.Printf("解析エラー: %v\n", err)
			} else {
				fmt.Println("リアルタイム解析完了", filepath)
			}
		})
		if err != nil {
			log.Fatalf("監視エラー: %v\n", err)
		}
	}()

	// コマンドライン引数の定義
	outputCSV := flag.Bool("csv", false, "CSVファイルに出力する")
	outputJSON := flag.Bool("json", false, "JSONファイルに出力する")
	outputChart := flag.Bool("chart", false, "エラーステータスのグラフを出力する")
	apiMode := flag.Bool("api", false, "HTTP APIサーバーを起動する")
	flag.Parse()

	// 残りの引数をログファイルのリストとして取得
	logFiles := flag.Args()
	if len(logFiles) == 0 {
		log.Fatal("ログファイルを少なくとも1つ指定してください")
	}

	// エラー統計を格納するマップ
	errorStats := make(map[string]int)
	timeStats := make(map[string]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 並列処理でログファイルを解析
	for _, logFile := range logFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			stats, times, err := parser.ExtractErrorStats(file)
			if err != nil {
				log.Printf("ファイル %s の解析に失敗: %v", file, err)
				return
			}

			// 結果を統合する
			// 競合を防ぐためにミューテックスを使用
			mu.Lock()
			for key, count := range stats {
				errorStats[key] += count
			}
			for key, count := range times {
				timeStats[key] += count
			}
			mu.Unlock()
		}(logFile)
	}

	// 全てのゴルーチンを待つ
	wg.Wait()

	// エラー統計を表示
	fmt.Println("【エラー統計】")
	for errorType, count := range errorStats {
		fmt.Printf("%s: %d 回\n", errorType, count)
	}

	// 時間別のエラー
	fmt.Println("\n【時間帯別エラー発生数】")
	for timeRange, count := range timeStats {
		fmt.Printf("%s: %d 回\n", timeRange, count)
	}

	// CSVオプションが指定されていたらCSVファイルに保存
	if *outputCSV {
		output.SaveResultsToCSV("output/csv/error_status.csv", errorStats)
		output.SaveResultsToCSV("output/csv/time_status.csv", timeStats)
	}

	// JSONオプションが指定されていたらJSONファイルに保存
	if *outputJSON {
		output.SaveResultsToJSON("output/json/error_status.json", errorStats)
		output.SaveResultsToJSON("output/json/time_status.json", timeStats)
	}

	// Chartオプションが指定されたていたらHTMLファイルに保存
	if *outputChart {
		output.GenerateBarChart("output/html/error_status.html", errorStats)
		output.GenerateBarChart("output/html/time_status.html", timeStats)
	}

	// APIモードが有効ならHTTPサーバー起動
	if *apiMode {
		// /api/errosでエラー統計 /api/timesで時間統計を返す
		http.HandleFunc("/api/errors", api.HandlerErrorStats(errorStats))
		http.HandleFunc("/api/times", api.HandlerTimeStats(timeStats))
		log.Println("HTTP APIサーバーをポート8080で起動中...")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	// 無限ループ
	select {}
}

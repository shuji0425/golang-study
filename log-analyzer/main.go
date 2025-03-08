package main

import (
	"flag"
	"fmt"
	"log"
	"log-analyzer/output"
	"log-analyzer/parser"
	"sync"
)

func main() {
	// コマンドライン引数の定義
	outputCSV := flag.Bool("csv", false, "CSVファイルに出力する")
	outputJSON := flag.Bool("json", false, "JSONファイルに出力する")
	outputChart := flag.Bool("chart", false, "エラーステータスのグラフを出力する")
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
		output.SaveResultsToCSV("logs/csv/error_status.csv", errorStats)
		output.SaveResultsToCSV("logs/csv/time_status.csv", timeStats)
	}

	// JSONオプションが指定されていたらJSONファイルに保存
	if *outputJSON {
		output.SaveResultsToJSON("logs/json/error_status.json", errorStats)
		output.SaveResultsToJSON("logs/json/time_status.json", timeStats)
	}

	// Chartオプションが指定されたていたらHTMLファイルに保存
	if *outputChart {
		output.GenerateBarChart("logs/html/error_status.html", errorStats)
		output.GenerateBarChart("logs/html/time_status.html", timeStats)

	}
}

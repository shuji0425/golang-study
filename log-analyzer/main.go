package main

import (
	"flag"
	"log"
	"log-analyzer/api"
	"log-analyzer/batch"
	"log-analyzer/processor"
	"net/http"
)

var (
	// コマンドライン引数の定義
	outputCSV   = flag.Bool("csv", false, "CSVファイルに出力する")
	outputJSON  = flag.Bool("json", false, "JSONファイルに出力する")
	outputChart = flag.Bool("chart", false, "エラーステータスのグラフを出力する")
	apiMode     = flag.Bool("api", false, "HTTP APIサーバーを起動する")
)

func main() {
	flag.Parse()

	// ファイル指定がなければlogs/*.logから全て取得
	logFiles := flag.Args()
	if len(logFiles) == 0 {
		logFiles = processor.GetAllLogFiles()
	}

	// 初回の解析と出力
	errorStats, timeStats := processor.ProcessLogs(logFiles)
	if *outputCSV || *outputJSON || *outputChart || *apiMode {
		processor.WriteOutputs(errorStats, timeStats)
	}

	// APIモードが有効ならHTTPサーバー起動
	if *apiMode {
		// HTTP APIサーバーを起動
		go func() {
			// /api/errosでエラー統計 /api/timesで時間統計を返す
			http.HandleFunc("/api/errors", api.HandlerErrorStats())
			http.HandleFunc("/api/times", api.HandlerTimeStats())
			log.Println("HTTP APIサーバーをポート8080で起動中...")
			log.Fatal(http.ListenAndServe(":8080", nil))
		}()
	}

	// ファイル変更があったときは全ファイルを再解析して出力更新
	processor.StartWatcher("logs", func() {
		// 毎回最新のログファイルを取得
		logFiles := processor.GetAllLogFiles()
		newErrorStats, newTimeStats := processor.ProcessLogs(logFiles)
		processor.WriteOutputs(newErrorStats, newTimeStats)
		api.UpdateStats(newErrorStats, newTimeStats)
	})

	// 定期的にエラーを書き込む
	go batch.WriteLogs()

	// 無限ループ
	select {}
}

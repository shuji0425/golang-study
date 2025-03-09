package processor

import (
	"fmt"
	"log"
	"log-analyzer/api"
	"log-analyzer/output"
	"log-analyzer/parser"
	"log-analyzer/watcher"
	"path/filepath"
	"sort"
	"sync"
)

// 指定されたログファイル群を並列処理しエラー統計と時間帯別統計を返す
func ProcessLogs(logFiles []string) (map[string]int, map[string]int) {
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

	// 結果をソートする準備
	sortedErrorStats, errorTypes := sortByKeyString(errorStats)
	sortedTimeStats, errorTimes := sortByKeyString(timeStats)

	// APIに解析結果を更新
	api.UpdateStats(sortedErrorStats, sortedTimeStats)

	// エラー統計を表示
	fmt.Println("【エラー統計】")
	for _, result := range errorTypes {
		fmt.Println(result)
	}

	// 時間別のエラー
	fmt.Println("\n【時間帯別エラー発生数】")
	for _, result := range errorTimes {
		fmt.Println(result)
	}

	return sortedErrorStats, sortedTimeStats
}

// 文字列を昇順にする関数
func sortByKeyString(stats map[string]int) (map[string]int, []string) {
	// ソートするためにスライスに格納
	var sortedKeys []string
	for key := range stats {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// 統計データの文字列
	var sortedResults []string
	for _, key := range sortedKeys {
		sortedResults = append(sortedResults, fmt.Sprintf("%s: %d 回", key, stats[key]))
	}

	// 返却ように整形
	sortedStats := make(map[string]int)
	for _, sort := range sortedKeys {
		sortedStats[sort] = stats[sort]
	}

	return sortedStats, sortedResults
}

// 解析結果を各種出力に書き出す(CSV, JSON, HTML)
func WriteOutputs(errorStats, timeStats map[string]int) {
	// CSVオプションが指定されていたらCSVファイルに保存
	output.SaveResultsToCSV("output/csv/error_status.csv", errorStats)
	output.SaveResultsToCSV("output/csv/time_status.csv", timeStats)

	// JSONオプションが指定されていたらJSONファイルに保存
	output.SaveResultsToJSON("output/json/error_status.json", errorStats)
	output.SaveResultsToJSON("output/json/time_status.json", timeStats)

	// Chartオプションが指定されたていたらHTMLファイルに保存
	output.GenerateBarChart("output/html/error_status.html", errorStats)
	output.GenerateBarChart("output/html/time_status.html", timeStats)
}

// logsディレクトリを監視し、変更があった場合に再解析と出力
func StartWatcher(logDir string, processAndWrite func()) {
	// 監視を開始
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("ファイル監視中にパニック発生: %v", err)
			}
		}()

		err := watcher.WatchLogsDir(logDir, func(changedFile string) {
			log.Printf("変更検出: %v", changedFile)

			// 変更があった時は、最新のログファイル一覧を取得し解析
			files, err := filepath.Glob(logDir + "/*.log")
			if err != nil {
				log.Printf("最新ログファイルの取得エラー: %v", err)
				return
			}
			log.Println("再解析開始:", files)
			processAndWrite()
			fmt.Println("更新処理完了")
		})
		if err != nil {
			log.Fatalf("監視エラー: %v\n", err)
		}
	}()
}

// ファル指定がない時はlog/*.logから全てのファイルを取得
func GetAllLogFiles() []string {
	files, err := filepath.Glob("logs/*.log")
	if err != nil {
		log.Fatalf("ログファイル取得エラー: %v", err)
	}
	if len(files) == 0 {
		log.Fatal("logsフォルダにログファイルが見つかりませんでした")
	}
	return files
}

package output

import (
	"fmt"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// エラー統計から棒グラフを作成し、指定したHTMLファイルに出力する
func GenerateBarChart(filepath string, errorStats map[string]int) {
	// 棒グラフ用のBarチャートを作成
	bar := charts.NewBar()

	// データを格納するスライスを準備
	var categories []string
	var counts []opts.BarData

	// errorStatsからキーと値を抽出
	for key, count := range errorStats {
		categories = append(categories, key)
		counts = append(counts, opts.BarData{Value: count})
	}

	// グローバルオプションの設定
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "エラー発生回数",
			Subtitle: "エラーの種類ごとの発生回数",
		}),
	)

	// X軸にカテゴリ、Y軸に発生回数のデータを追加
	bar.SetXAxis(categories).AddSeries("発生回数", counts)

	// 複数のチャートをまとめて表示する場合は、Pageを使う
	page := components.NewPage()
	page.AddCharts(bar)

	// 出力ファイルを作成
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("ファイルの作成に失敗: %v\n", err)
		return
	}
	defer file.Close()

	// ページ全体をHTMLとしてレンダリング
	err = page.Render(file)
	if err != nil {
		fmt.Printf("チャートの作成に失敗: %v\n", err)
		return
	}

	fmt.Printf("エラーチャートを出力しました: %s\n", filepath)
}

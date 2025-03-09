package api

import (
	"encoding/json"
	"net/http"
)

// エラー統計情報をJSONで返すための構造体
type ErrorStats struct {
	Errors map[string]int `json:"errors"`
}

type TimeStats struct {
	Times map[string]int `json:"times"`
}

// エラー統計情報を返すハンドラー
func HandlerErrorStats(errorStats map[string]int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JSON形式で返すための記述
		w.Header().Set("Content-Type", "application/json")
		// JSONに出力
		json.NewEncoder(w).Encode(ErrorStats{Errors: errorStats})
	}
}

// 時間帯別エラー情報を返すハンドラー
func HandlerTimeStats(timeStats map[string]int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(TimeStats{Times: timeStats})
	}
}

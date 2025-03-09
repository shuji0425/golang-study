package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

// エラー統計情報をJSONで返すための構造体
type ErrorStats struct {
	Errors map[string]int `json:"errors"`
}

type TimeStats struct {
	Times map[string]int `json:"times"`
}

// 解析結果を保持
type Stats struct {
	mu         sync.RWMutex
	ErrorStats map[string]int
	TimeStats  map[string]int
}

// 現在の解析結果
var currentStats = &Stats{
	ErrorStats: make(map[string]int),
	TimeStats:  make(map[string]int),
}

// 新しい解析結果でcurrentStatsを更新
func UpdateStats(newErrors, newTimes map[string]int) {
	currentStats.mu.Lock()
	defer currentStats.mu.Unlock()
	currentStats.ErrorStats = newErrors
	currentStats.TimeStats = newTimes
}

// 現在の解析結果を返却
func GetStats() (map[string]int, map[string]int) {
	currentStats.mu.RLock()
	defer currentStats.mu.RUnlock()
	return currentStats.ErrorStats, currentStats.TimeStats
}

// エラー統計情報を返すハンドラー
func HandlerErrorStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JSON形式で返すための記述
		w.Header().Set("Content-Type", "application/json")
		// 最新を取得
		errorStats, _ := GetStats()
		// JSONに出力
		json.NewEncoder(w).Encode(ErrorStats{Errors: errorStats})
	}
}

// 時間帯別エラー情報を返すハンドラー
func HandlerTimeStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, timeStats := GetStats()
		json.NewEncoder(w).Encode(TimeStats{Times: timeStats})
	}
}

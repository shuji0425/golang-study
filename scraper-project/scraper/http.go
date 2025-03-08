package scraper

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// ランダムジェネレーターを作成
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// ランダムな待機時間を取得(1~3秒)
func getRandomSleepDuration() time.Duration {
	return time.Duration(rng.Intn(2000)+1000) * time.Millisecond
}

// 指定したURLからHTMLコンテンツを取得
func FetchHTML(url string) (string, error) {
	// レートリミット：待機時間を1~3秒のランダムにする
	time.Sleep(getRandomSleepDuration())

	// リクエストの作成
	req, err := CreateRequestWithUserAgent(url)
	if err != nil {
		return "", fmt.Errorf("リクエストの作成に失敗しました: %v", err)
	}

	// HTTPクライアントを作成
	client := &http.Client{}

	// HTTP GETリクエストの送信
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("リクエストの送信に失敗しました: %v", err)
	}
	// 閉じる
	defer resp.Body.Close()

	// レスポンスのステータスコードが200以外だとエラー
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTPステータスコードが200ではありません: %v", resp.StatusCode)
	}

	// レスポンスのボディ（HTML）を読み込む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("レスポンスの読み込みに失敗しました: %v", err)
	}

	// 文字列としてHTMLを返す
	return string(body), nil
}

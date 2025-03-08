package scraper

import (
	"fmt"
	"io"
	"net/http"
)

// 指定したURLからHTMLコンテンツを取得
func FetchHTML(url string) (string, error) {
	// HTTP GETリクエストの送信
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("リクエストの送信に失敗しました: %v", err)
	}
	// 閉じる
	defer resp.Body.Close()

	// レスポンスのステータスコードが200以外だとエラー
	if resp.StatusCode != 200 {
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

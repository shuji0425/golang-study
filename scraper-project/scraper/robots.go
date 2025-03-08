package scraper

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/temoto/robotstxt"
)

// robots.txtを取得し、許可されたパスか確認
func CheckRobotsTXT(baseURL string, targetURL string) (bool, error) {
	// robots.txtのURLを組み立て
	robotsURL := baseURL + "/robots.txt"

	// robots.txtを取得
	resp, err := http.Get(robotsURL)
	if err != nil {
		return false, fmt.Errorf("robots.txtの取得に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	// robots.txtを解析
	robotsData, err := robotstxt.FromResponse(resp)
	if err != nil {
		return false, fmt.Errorf("robots.txtの解析に失敗しました: %v", err)
	}

	// 対象URLがクロール許可されているかチェック
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return false, fmt.Errorf("URLの解析に失敗しました: %v", err)
	}

	// user-agent "*" は全てのクローラーに対する規則を意味する
	if robotsData.TestAgent(parsedURL.Path, "*") {
		return true, nil
	}

	return false, nil
}

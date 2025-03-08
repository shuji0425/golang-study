package scraper

import "net/http"

// User-Agentを変更する
func CreateRequestWithUserAgent(url string) (*http.Request, error) {
	// 通常のUser-Agentを使う
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// User-Agentの設定
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	return req, nil
}

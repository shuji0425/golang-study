package scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLからページタイトルを抽出
func ParseTitle(html string) (string, error) {
	// gogueryでHTMLを解析
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("HTMLの解析に失敗しました: %v", err)
	}

	// タイトルの内容を抽出
	title := doc.Find("title").Text()

	return title, err
}

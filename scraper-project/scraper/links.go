package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// HTMLから全てのリンクを抽出
func ExtractLinks(html string, baseURL string) ([]string, error) {
	var links []string

	// ベースURLを解析
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("ベースURLの解析に失敗しました: %v", err)
	}

	// goqueryを使ってHTML解析
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("HTMLの解析に失敗しました: %v", err)
	}

	// <a>タグのhrefを取得
	doc.Find("a").Each(func(index int, element *goquery.Selection) {
		link, exists := element.Attr("href")
		if exists {
			// 相対URLを絶対URLに変換
			parsedLink, err := url.Parse(link)
			if err == nil {
				absoluteLink := base.ResolveReference(parsedLink).String()
				links = append(links, absoluteLink)
			}
		}
	})

	return links, nil
}

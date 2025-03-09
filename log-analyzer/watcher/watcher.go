package watcher

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

// logsフォルダの変更を監視する
func WatchLogsDir(logDir string, callback func(string)) error {
	// fsnotifyの新しいインスタンを作成
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// 監視対象のディレクトリを追加
	err = watcher.Add(logDir)
	if err != nil {
		return fmt.Errorf("ディレクトリの監視に失敗しました: %w", err)
	}

	log.Printf("監視を開始: %s\n", logDir)

	// 監視ループ
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				log.Println("監視イベントの受信に失敗")
				return nil
			}
			log.Println("イベント発生:", event)

			// 書き込みイベントを検出
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				log.Printf("変更を検出: %s\n", event.Name)
				// 変更があったファイルを解析
				callback(event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("エラー: %v\n", err)
		}
	}
}

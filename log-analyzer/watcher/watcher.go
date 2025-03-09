package watcher

import (
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

	// 指定したディレクトリを監視対象に追加
	err = watcher.Add(logDir)
	if err != nil {
		return err
	}

	log.Printf("監視を開始: %s\n", logDir)

	// 監視ループ
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// 書き込みイベントを検出
			if event.Op&fsnotify.Write == fsnotify.Write {
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

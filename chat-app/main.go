package main

import (
	"log"
	"net/http"
)

func main() {
	// Hubの初期化（クライアント管理＆メッセージのブロードキャスト）
	hub := newHub()
	// Hubのループ処理をゴルーチンで実行
	go hub.run()

	// WebSocket接続用エンドポイントの設定
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(hub, w, r)
	})

	// publicディレクトリ内の静的ファイルを提供する設定
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Println("サーバーを開始: ポート8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServerエラー: ", err)
	}
}

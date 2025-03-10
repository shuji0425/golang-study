package main

import (
	"chat-app/auth"
	"chat-app/handlers"
	"chat-app/hub"
	"log"
	"net/http"
)

func main() {
	// Hubの初期化（クライアント管理＆メッセージのブロードキャスト）
	h := hub.NewHub()
	// Hubのループ処理をゴルーチンで実行
	go h.Run()

	// 認証不要ルート
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// 認証が必要なルート
	http.Handle("/chat", auth.AuthMiddleware(http.HandlerFunc(handlers.ChatHandler)))

	// WebSocket接続用エンドポイントの設定
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServerWs(h, w, r)
	})

	// publicディレクトリ内の静的ファイルを提供する設定
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	log.Println("サーバーを開始: ポート8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServerエラー: ", err)
	}
}

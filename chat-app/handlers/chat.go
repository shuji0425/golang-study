package handlers

import (
	"chat-app/auth"
	"net/http"
	"text/template"
)

// チャット画面へ遷移
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// セッションからユーザー名を取得
	username, _ := auth.GetUsernameFromSession(r)
	// テンプレートにデータを渡す
	data := struct {
		Username string
	}{
		Username: username,
	}

	tmpl := template.Must(template.ParseFiles("templates/chat.html"))
	tmpl.Execute(w, data)
}

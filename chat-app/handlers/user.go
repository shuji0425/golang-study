package handlers

import (
	"chat-app/auth"
	"chat-app/models"
	"net/http"
	"text/template"
)

// ログイン画面の表示と認証処理を行う
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if models.Authenticate(username, password) {
			auth.SetUserSession(w, r, username)
			http.Redirect(w, r, "/chat", http.StatusFound)
			return
		}
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, "認証に失敗しました")
	}
}

// 新規ユーザー登録を処理
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
		return
	} else if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		models.CreateUser(username, password)
		auth.SetUserSession(w, r, username)
		http.Redirect(w, r, "/chat", http.StatusFound)
	}
}

// ログアウト処理
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	auth.ClearUserSession(w, r)
	http.Redirect(w, r, "/login", http.StatusFound)
}

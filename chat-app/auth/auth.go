package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var sessionStore = sessions.NewCookieStore([]byte("secret-key"))

const SessionName = "chatapp-session"

// ログインしたユーザーの情報をセッションに保存
func SetUserSession(w http.ResponseWriter, r *http.Request, username string) error {
	session, _ := sessionStore.Get(r, SessionName)
	session.Values["username"] = username
	return session.Save(r, w)
}

// セッション情報を削除
func ClearUserSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := sessionStore.Get(r, SessionName)
	delete(session.Values, "username")
	return session.Save(r, w)
}

// セッションからユーザー名を取得
func GetUsernameFromSession(r *http.Request) (string, bool) {
	session, _ := sessionStore.Get(r, SessionName)
	username, ok := session.Values["username"].(string)
	return username, ok
}

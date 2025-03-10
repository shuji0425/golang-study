package models

// ユーザー情報
type User struct {
	Username string
	Password string
}

var users = map[string]*User{}

// 新しいユーザーを登録
func CreateUser(username, password string) *User {
	user := &User{Username: username, Password: password}
	users[username] = user
	return user
}

// ユーザー名とパスワードが一致するか確認
func Authenticate(username, password string) bool {
	user, exists := users[username]
	if !exists {
		return false
	}
	return user.Password == password
}

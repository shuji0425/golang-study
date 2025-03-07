package models

// User構造体は、ユーザーの情報を表す
type User struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

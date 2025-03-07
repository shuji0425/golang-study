package main

import (
	"fmt"
	"log"
	"net/http"
	"user-management-api/db"
	"user-management-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	// DB接続を初期化
	db.InitDB()

	// ルーター作成
	r := mux.NewRouter()

	// ルートの設定
	r.HandleFunc("/user", handlers.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	// サーバーの起動
	fmt.Println("Server is running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}

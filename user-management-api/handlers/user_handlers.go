package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user-management-api/db"
	"user-management-api/models"

	"github.com/gorilla/mux"
)

// ユーザー一覧を返す
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	// DBからユーザー情報を取得
	db.DB.Find(&users)

	// JSON形式でレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 特定のユーザーを取得する
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// パスパラメータからIDを取得
	vars := mux.Vars(r)
	id := vars["id"]

	// DBからユーザー情報を取得
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, fmt.Sprintf("User with id %s not found", id), http.StatusNotFound)
		return
	}

	// JSON形式で返却
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ユーザーを作成する
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// リクエストボディからユーザー情報をデコード
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// DBにユーザー情報を保存
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// 作成したユーザーをレスポンスとして返す
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ユーザーを更新する
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user []models.User
	vars := mux.Vars(r)
	id := vars["id"]

	// DBから既存のユーザーを取得
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, fmt.Sprintf("User with id %s not found", id), http.StatusNotFound)
		return
	}

	// リクエストボディから新しいユーザー情報をエンコード
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// DBでユーザー情報を更新
	if err := db.DB.Save(&user).Error; err != nil {
		http.Error(w, "Faild to update user", http.StatusInternalServerError)
		return
	}

	// 更新したユーザーをレスポンスとして返す
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ユーザーを削除する
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user []models.User
	vars := mux.Vars(r)
	id := vars["id"]

	// DBからユーザー情報を取得
	if err := db.DB.First(&user, id).Error; err != nil {
		http.Error(w, fmt.Sprintf("User with id %s not found", id), http.StatusNotFound)
		return
	}

	// DBからユーザー情報を削除
	if err := db.DB.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// 成功レスポンス
	w.WriteHeader(http.StatusNoContent)
}

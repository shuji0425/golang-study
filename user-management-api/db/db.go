package db

import (
	"fmt"
	"log"
	"user-management-api/models"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// DBはデータベース接続のインスタンス
var DB *gorm.DB

// データベースの接続を初期化
func InitDB() {
	var err error
	// PostgreSQLへの接続
	dsn := "host=localhost user=yourusername password=yourpassword dbname=yourdbname port=5432 sslmode=disable"
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to conect to database(接続できなかった)", err)
	}

	// マイグレーション：ユーザーテーブルを作成
	DB.AutoMigrate(&models.User{})

	// 成功した時のメッセージ
	fmt.Println("Database connected successfully!(接続成功)")
}

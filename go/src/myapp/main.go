package main

import (
	"myapp/config"
	"myapp/database"

	"myapp/routes"
)

func main() {
	// 設定ファイルの読み込み
	config.LoadConfig()

	// // データベースに接続
	// db := database.ConnectDB()
	// defer db.Close()

	// // データベース接続の初期化
	// repository.InitDB()
	// defer repository.CloseDB()

	// // データベース接続の初期化
	database.ConnectDB()
	defer database.CloseDB()

	// 自動マイグレーション
	database.AutoMigrate()

	// 自動マイグレーション
	// repository.RepositoryAutoMigrate()

	// // 全てのテーブルを取得して表示
	// tables, err := database.GetAllTables(db)
	// if err != nil {
	// 	log.Fatalf("Error getting all tables: %v", err)
	// }

	// fmt.Println("Tables in the database:")
	// for _, table := range tables {
	// 	fmt.Println(table)
	// }

	// ルーターの設定
	router := routes.SetupRouter()

	// サーバーを開始
	router.Run(":8080")
}

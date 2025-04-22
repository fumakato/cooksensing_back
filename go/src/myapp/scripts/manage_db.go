/*
scripts説明
データベースを操作するときに使うよ
初期設定だったり
全データ削除だったり

シードデータ（初期値）挿入
go run scripts/manage_db.go -action init

全テーブル削除
go run scripts/manage_db.go -action dropall

全テーブル削除と自動マイグレーション
go run scripts/manage_db.go -action dropmigrate

任意のテーブル削除
go run scripts/manage_db.go -action droptable -table users
（ users の部分は各テーブルの名称に変更して使用）

BestDataテーブルを作成
go run scripts/manage_db.go -action bestdata

BestDataテーブルの平均の算出
go run scripts/manage_db.go -action bestaverage

ヒストグラムテーブル作成
go run scripts/manage_db.go -action histogram

ベストデータとヒストグラムの更新
go run scripts/manage_db.go -action fixData

特徴量データの全件表示
go run scripts/manage_db.go -action getAllFeatureData

特徴量データをuserIDで指定して表示
go run scripts/manage_db.go -action getFeatureDataByUserID -user_id 2

データベースの特徴量データテーブルの任意のidのデータを削除
go run scripts/manage_db.go -action deleteFeatureDatabyid -id {数字}
go run scripts/manage_db.go -action deleteFeatureDatabyid -id 3

userの情報全件取得
go run scripts/manage_db.go -action getAllUsers

*/

package main

import (
	"flag"
	"fmt"
	"log"

	"myapp/config"
	"myapp/database"
	// "myapp/model"
)

func initDB() {
	database.ConnectDB()
	database.InitData()
	database.CloseDB()
}

func dropAllTables() {
	database.ConnectDB()
	database.DropAllTables()
	database.CloseDB()
}

func dropmigrate() {
	database.ConnectDB()
	database.DropAllTables()
	database.AutoMigrate()
	database.CloseDB()
}

func dropTable(tableName string) {
	database.ConnectDB()
	database.DropTables(tableName)
	database.CloseDB()
}

func bestData() {
	database.ConnectDB()
	database.UpdateBestDataFromFeatureData()
	database.CloseDB()
}

func bestDataAverage() {
	database.ConnectDB()
	// database.AveragePaceAndAccelerationStdDev()
	averagePace, accelerationStdDev, err := database.AveragePaceAndAccelerationStdDev()
	if err != nil {
		log.Fatalf("Error calculating averages: %v", err)
	}

	fmt.Printf("Average Pace: %f\n", averagePace)
	fmt.Printf("Average Acceleration StdDev: %f\n", accelerationStdDev)
	database.CloseDB()
}
func histogram() {
	database.ConnectDB()
	database.GenerateAndStoreHistogramData()
	database.CloseDB()
}

func fixData() {
	database.ConnectDB()
	database.UpdateBestDataFromFeatureData()
	database.GenerateAndStoreHistogramData()
	database.AssignBestClassToAll()
	database.CloseDB()
}

func deleteFeatureDataByid(id uint) {
	database.ConnectDB()
	err := database.DeleteFeatureData(id)
	if err != nil {
		log.Fatalf("❌ Failed to delete FeatureData with ID %d: %v", id, err)
	}
	fmt.Printf("✅ FeatureData with ID %d deleted successfully\n", id)
	database.CloseDB()
}

func getAllFeatureData() {
	database.ConnectDB()
	// 全件取得
	featureDataList, err := database.GetAllFeatureData()
	if err != nil {
		log.Fatalf("❌ データ取得に失敗しました: %v", err)
	}

	// 件数ログ
	fmt.Printf("📊 取得件数: %d 件\n", len(featureDataList))

	// 中身を1件ずつ出力
	for _, data := range featureDataList {
		fmt.Printf("🟢 ID: %d | UserID: %d | ActionID: %d | Date: %s | AvgPace: %.16f | StdDev: %.16f | CreatedAt: %s\n",
			data.ID,
			data.UserID,
			data.ActionID,
			data.Date.Format("2006-01-02 15:04"),
			data.AveragePace,
			data.AccelerationStandardDeviation,
			data.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
	database.CloseDB()
}

func getFeatureDataByUserID(userID uint) {
	database.ConnectDB()
	defer database.CloseDB()

	// UserIDで取得
	featureDataList, err := database.GetFeatureDataByUserID(userID)
	if err != nil {
		log.Fatalf("❌ データ取得に失敗しました (UserID=%d): %v", userID, err)
	}

	// 件数ログ
	fmt.Printf("📊 UserID %d の取得件数: %d 件\n", userID, len(featureDataList))

	// 中身を1件ずつ出力
	for _, data := range featureDataList {
		fmt.Printf("🟢 ID: %d | UserID: %d | ActionID: %d | Date: %s | AvgPace: %.16f | StdDev: %.16f | CreatedAt: %s\n",
			data.ID,
			data.UserID,
			data.ActionID,
			data.Date.Format("2006-01-02 15:04"),
			data.AveragePace,
			data.AccelerationStandardDeviation,
			data.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
}

func getAllUsers() {
	database.ConnectDB()
	defer database.CloseDB()

	// 全件取得
	userList, err := database.FindAllUser()
	if err != nil {
		log.Fatalf("❌ ユーザー取得に失敗しました: %v", err)
	}

	// 件数ログ
	fmt.Printf("📊 ユーザー取得件数: %d 件\n", len(userList))

	// 中身を1件ずつ出力
	for _, user := range userList {
		fmt.Printf("🟢 ID: %d | Name: %s | FirebaseAuthUID: %s | TsukurepoID: %s | CreatedAt: %s\n",
			user.ID,
			user.Name,
			user.FirebaseAuthUID,
			user.TsukurepoID,
			user.CreatedAt.Format("2006-01-02 15:04"),
		)
	}
}

func main() {
	config.LoadConfig()

	action := flag.String("action", "", "Action to perform (init, dropall, droptable)")
	tableName := flag.String("table", "", "Name of the table to drop (required if action is droptable)")
	id := flag.Uint("id", 0, "ID of the record to delete")
	userID := flag.Uint("user_id", 0, "UserID to fetch FeatureData for")
	flag.Parse()

	if *action == "" {
		log.Fatal("Action must be provided")
	}

	switch *action {
	case "init":
		initDB()
	case "dropall":
		dropAllTables()
	case "droptable":
		if *tableName == "" {
			log.Fatal("Table name must be provided for droptable action")
		}
		dropTable(*tableName)
	case "bestdata":
		bestData()
	case "bestaverage":
		bestDataAverage()
	case "histogram":
		histogram()
	case "dropmigrate":
		dropmigrate()
	case "fixData":
		fixData()
	case "getAllFeatureData":
		getAllFeatureData()

	case "deleteFeatureDatabyid":
		if *id == 0 {
			log.Fatal("ID must be provided for deleteFeatureDatabyid action")
		}
		deleteFeatureDataByid(*id)

	case "getFeatureDataByUserID":
		if *userID == 0 {
			log.Fatal("UserID must be provided for getFeatureDataByUserID action")
		}
		getFeatureDataByUserID(*userID)

	case "getAllUsers":
		getAllUsers()

	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

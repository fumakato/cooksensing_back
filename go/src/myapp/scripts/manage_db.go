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

func main() {
	config.LoadConfig()

	action := flag.String("action", "", "Action to perform (init, dropall, droptable)")
	tableName := flag.String("table", "", "Name of the table to drop (required if action is droptable)")
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
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

// gormを使ってHistogramの中身を操作するコードを作ってください．
// type Histogram struct {
// 	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
// 	BaseModel
// 	DisplayItemID uint `json:"display_item_id"`
// 	ActionID      uint `json:"action_id"`
// 	Time1         uint `gorm:"default:0" json:"time1"`
// 	Time2         uint `gorm:"default:0" json:"time2"`
// 	Time3         uint `gorm:"default:0" json:"time3"`
// 	Time4         uint `gorm:"default:0" json:"time4"`
// 	Time5         uint `gorm:"default:0" json:"time5"`
// 	Time6         uint `gorm:"default:0" json:"time6"`
// 	Time7         uint `gorm:"default:0" json:"time7"`
// 	Time8         uint `gorm:"default:0" json:"time8"`
// 	Time9         uint `gorm:"default:0" json:"time9"`
// 	Time10        uint `gorm:"default:0" json:"time10"`
// }
// HistogramテーブルはAveragePaceとAccelerationStandardDeviationのデータを用いたヒストグラムのデータを格納しています．
// BesteDataから全件を取り出します．
// AveragePace内で最大値と最小値を求めます．最大値-最小値でデータの範囲を求めます．データの範囲を10分割し，それぞれTime1からTime10に割り当てます．
// それぞれの範囲に所属する人数を求めてテーブルに格納します．
// AccelerationStandardDeviationも同様に行います．
// また，最大値と最小値，データの範囲を標準出力に出力してください

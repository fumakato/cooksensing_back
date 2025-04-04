package repository

// import (
// 	"log"
// 	"myapp/database"

// 	"gorm.io/gorm"
// )

// var db *gorm.DB

// func InitDB() {
// 	db = database.ConnectDB()
// 	if db == nil {
// 		log.Fatalf("Failed to connect to the database.")
// 	}
// }

// func CloseDB() {
// 	if db != nil {
// 		sqlDB, err := db.DB()
// 		if err == nil {
// 			sqlDB.Close()
// 		}
// 	}
// }

// func RepositoryAutoMigrate() {
// 	if db != nil {
// 		// グローバル変数dbを渡す
// 		database.AutoMigrate(db)
// 	} else {
// 		log.Fatalf("Database is not initialized.")
// 	}
// }

package database

import (
	"log"
	"myapp/config"
	"myapp/model"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// DB接続
func ConnectDB() {
	dsn := config.GetDBConfig()
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
}

func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// 自動マイグレーション
// func AutoMigrate(db *gorm.DB) {
func AutoMigrate() {
	if db != nil {
		db.AutoMigrate(&model.User{})
		db.AutoMigrate(&model.FeatureData{})
		db.AutoMigrate(&model.BestData{})
		db.AutoMigrate(&model.Histogram{})
		db.AutoMigrate(&model.Action{})
		db.AutoMigrate(&model.DisplayItem{})
		// db.AutoMigrate(&model.Label{})
	} else {
		log.Fatalf("Database is not initialized.")
	}
}

// 初期データの登録
// func InitData(db *gorm.DB) {
func InitData() {
	if err := db.First(&model.User{}).Error; err == gorm.ErrRecordNotFound {
		for _, tmp := range model.UserInitData {
			db.Create(&tmp)
		}
	}
	if err := db.First(&model.FeatureData{}).Error; err == gorm.ErrRecordNotFound {
		for _, tmp := range model.FeatureDataInitData {
			db.Create(&tmp)
		}
	}
	if err := db.First(&model.BestData{}).Error; err == gorm.ErrRecordNotFound {
		for _, tmp := range model.BestDataData {
			db.Create(&tmp)
		}
	}
	if err := db.First(&model.Histogram{}).Error; err == gorm.ErrRecordNotFound {
		for _, tmp := range model.HistogramData {
			db.Create(&tmp)
		}
	}
	if err := db.First(&model.Action{}).Error; err == gorm.ErrRecordNotFound {
		for _, tmp := range model.ActionInitData {
			db.Create(&tmp)
		}
	}
	if err := db.First(&model.DisplayItem{}).Error; err == gorm.ErrRecordNotFound {
		for _, tmp := range model.DisplayItemInitData {
			db.Create(&tmp)
		}
	}
}

func DropAllTables() {
	if db != nil {
		db.Migrator().DropTable(&model.User{})
		db.Migrator().DropTable(&model.FeatureData{})
		db.Migrator().DropTable(&model.BestData{})
		db.Migrator().DropTable(&model.Histogram{})
		db.Migrator().DropTable(&model.Action{})
		db.Migrator().DropTable(&model.DisplayItem{})
	} else {
		log.Fatalf("Database is not initialized.")
	}
}

func DropTables(tableName string) {
	if db != nil {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)).Error; err != nil {
			log.Fatalf("Could not drop table %s: %v", tableName, err)
		}
	} else {
		log.Fatalf("Database is not initialized.")
	}
}

func GetAllTables() ([]string, error) {
	if db != nil {
		var tables []string
		rows, err := db.Raw("SHOW TABLES").Rows()
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				return nil, err
			}
			tables = append(tables, table)
		}
		return tables, nil
	} else {
		log.Fatalf("Database is not initialized.")
		return nil, nil
	}
}

// func ConnectDB_old() *gorm.DB {
// 	dsn := config.GetDBConfig()
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Could not connect to the database: %v", err)
// 	}
// 	return db
// }

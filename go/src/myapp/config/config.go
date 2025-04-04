package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetDBConfig() string {
	// dbms := os.Getenv("DB_MS")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	protocol := os.Getenv("DB_PROTOCOL")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, protocol, dbName)
}

// 実行方法
// 以下をターミナルに打ち込むといける
// go run scripts/init_db.go

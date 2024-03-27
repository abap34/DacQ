package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"log"
)

var (
	leaderbord_db *gorm.DB
	submitlog_db  *gorm.DB
)

func getEnv(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        value = fallback
    }
    return value
}

func InitDB() error {
	user := getEnv("NS_MARIADB_USER", "admin")
	pass := getEnv("NS_MARIADB_PASSWORD", "password")
	host := getEnv("NS_MARIADB_HOSTNAME", "localhost")
	dbname := getEnv("NS_MARIADB_DATABASE", "database")


	log.Printf("user: %s, pass: %s, host: %s, dbname: %s", user, pass, host, dbname)

	// リーダーボードの DB 
	_db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=True&loc=Asia%2FTokyo&charset=utf8mb4"), &gorm.Config{})

	if err != nil {
		return err
	}

	leaderbord_db = _db
	leaderbord_db.AutoMigrate(&Score{})

	// イベントログの DB 
	_db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=True&loc=Asia%2FTokyo&charset=utf8mb4"), &gorm.Config{})
	if err != nil {
		return err
	}

	submitlog_db = _db
	submitlog_db.AutoMigrate(&SubmitLog{})
	return nil
}

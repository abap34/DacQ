package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "os"
)

var (
	leaderbord_db *gorm.DB
	submitlog_db  *gorm.DB
)

// 順位表の DB と、イベントログの DB を作成する
func InitDB() error {
	// user := os.Getenv("NS_MARIADB_USER")
	// pass := os.Getenv("NS_MARIADB_PASSWORD")
	// host := os.Getenv("NS_MARIADB_HOSTNAME")
	// dbname := os.Getenv("NS_MARIADB_DATABASE")
	user := "admin"
	pass := "password"
	host := "localhost"
	dbname := "database"

	_db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=True&loc=Asia%2FTokyo&charset=utf8mb4"), &gorm.Config{})

	if err != nil {
		return err
	}

	leaderbord_db = _db
	leaderbord_db.AutoMigrate(&Score{})

	// イベントログの DB を作成する
	_db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=True&loc=Asia%2FTokyo&charset=utf8mb4"), &gorm.Config{})
	if err != nil {
		return err
	}

	submitlog_db = _db
	submitlog_db.AutoMigrate(&SubmitLog{})
	return nil
}

package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	db *gorm.DB
)

func InitDB() error {
	user := os.Getenv("MARIADB_USERNAME")
	pass := os.Getenv("MARIADB_PASSWORD")
	host := os.Getenv("MARIADB_HOSTNAME")
	dbname := os.Getenv("MARIADB_DATABASE")

	_db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, pass, host, dbname)+"?parseTime=True&loc=Asia%2FTokyo&charset=utf8mb4"), &gorm.Config{})

	if err != nil {
		return err
	}

	db = _db
	return db.AutoMigrate(&Score{})
}

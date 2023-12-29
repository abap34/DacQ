package model

import (
	"fmt"
	"time"
)

type SubmitLog struct {
	User   string  `gorm:"not null"`
	Time   string  `gorm:"not null"`
	IsBest bool    `gorm:"not null"`
	Score  float64 `gorm:"not null"`
}

func CreateSubmitLog(submitlog SubmitLog) error {
	return submitlog_db.Create(&submitlog).Error
}

func GetSubmitLogByUser(user string) ([]SubmitLog, error) {
	var submitlogs []SubmitLog
	err := submitlog_db.Where("user = ?", user).Find(&submitlogs).Error
	return submitlogs, err
}

func CountSubmitLogByUser(user string) (int64, error) {
	var count int64
	err := submitlog_db.Model(&SubmitLog{}).Where("user = ?", user).Count(&count).Error
	return count, err
}

// 直近 n 件の IsBest が true である SubmitLog を取得する
func GetNews(n int) ([]SubmitLog, error) {
	var submitlogs []SubmitLog
	err := submitlog_db.Where("is_best = ?", true).Order("time desc").Limit(n).Find(&submitlogs).Error
	return submitlogs, err
}

// 最終サブミットがどれくらい前かを返す
// 0~59秒: "n秒前"
// 1~59分: "n分前"
// 1~23時間: "n時間前"
// 1~6日: "n日前"
// 7日以上: "n週間前"
func GetLastSubmitTime(user string) (string, error) {
	var submitlog SubmitLog
	err := submitlog_db.Where("user = ?", user).Order("time desc").First(&submitlog).Error

	if err != nil {
		if err.Error() == "record not found" {
			return "過去の投稿が見つかりませんでした", nil
		} else {
			return "", err
		}
	}

	if err != nil {
		return "", err
	}

	lasttime, err := time.Parse("2004-03-04 10:01:01", submitlog.Time)
	if err != nil {
		return "", err
	}

	diff := time.Now().Sub(lasttime)

	if diff.Seconds() < 60 {
		return fmt.Sprintf("%d秒前", int(diff.Seconds())), nil
	} else if diff.Minutes() < 60 {
		return fmt.Sprintf("%d分前", int(diff.Minutes())), nil
	} else if diff.Hours() < 24 {
		return fmt.Sprintf("%d時間前", int(diff.Hours())), nil
	} else if diff.Hours() < 24*7 {
		return fmt.Sprintf("%d日前", int(diff.Hours()/24)), nil
	} else {
		return fmt.Sprintf("%d週間前", int(diff.Hours()/24/7)), nil
	}
}

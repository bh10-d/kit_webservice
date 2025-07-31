package model

import "gorm.io/gorm"

type Runner struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
	Tags string `json:"tags"`
}

type Job struct {
	ID             uint   `gorm:"primaryKey"`
	RunnerID       string
	MsgID          string
	Status         string
	RequestPayload string
	ResponsePayload string
	Timeout        bool
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Runner{}, &Job{})
}

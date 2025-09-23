package model

import (
	"gorm.io/gorm"
	"github.com/lib/pq"
)

type Runner struct {
	ID   string `gorm:"primaryKey" json:"id"`
	HostName string `json:"hostname"`
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

type Scripts struct {
	ScriptID   string `gorm:"primaryKey" json:"script_id"`
	FileName   string `json:"file_name"`
	Description string `json:"description"`
	Param      pq.StringArray  `gorm:"type:text[]" json:"param"`
}

type Logs struct {
	MsgID    string `gorm:"primaryKey" json:"msg_id"`
	RunnerID string `gorm:"primaryKey" json:"runner_id"`
	Logs     string `json:"logs"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Runner{}, &Job{}, &Logs{}, &Scripts{})
}

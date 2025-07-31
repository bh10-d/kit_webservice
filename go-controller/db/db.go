package db

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"go-controller/model"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("../controller/runners.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	model.AutoMigrate(DB)
}

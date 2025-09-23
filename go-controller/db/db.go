package db

import (
	"log"
	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"go-controller/model"
	"os"
	"fmt"
)

var DB *gorm.DB

func InitDB() {
	
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		host, user, password, dbname, port,
	)

	// var err error
	// DB, err = gorm.Open(sqlite.Open("../controller/runners.db"), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("failed to connect database: %v", err)
	// }

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	model.AutoMigrate(DB)
}

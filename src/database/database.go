package controllers

import (
	"fmt"
	"os"

	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDB() {
	dsn := "minitwit:" + os.Getenv("DB_PASS") + "@tcp(db:3306)/minitwit?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		panic("Failed to connect to database.")
	}
	db.AutoMigrate(&model.User{}, &model.Message{}, &model.Follow{})
	DB = db
}

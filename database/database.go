package database

import (
	"fmt"
	"github/s6352410016/go-fiber-gorm-rest-api-auth-jwt-postgresql/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed To Connect Database\n", err.Error())
		os.Exit(2)
	}

	log.Println("Connected To Database Successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	db.AutoMigrate(&models.User{})
	DB = db
}

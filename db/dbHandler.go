package db

import (
	"AuthService/internal/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	dbServer := os.Getenv("DB_SERVER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbUserPassword := os.Getenv("DB_USER_PASSWORD")

	if dbServer == "" || dbName == "" || dbPort == "" || dbUsername == "" || dbUserPassword == "" {
		log.Fatal("Database configuration is not fully set in .env")
	}

	connectionStr := fmt.Sprintf("host=%s dbname=%s port=%s user=%s password=%s sslmode=disable",
		dbServer, dbName, dbPort, dbUsername, dbUserPassword)

	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db

	Migrate()
}

func Migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.RefreshToken{})
	if err != nil {
		log.Fatal(err)
	}
}

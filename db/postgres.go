package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PG *gorm.DB

func ConnectPostgres() {
	// Формируем DSN на основе .env
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	if dsn == "" {
		log.Fatal("Postgres DSN is empty! Проверь .env или docker-compose")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Postgres connect error: %v", err)
	}

	PG = db
	log.Println("✅ Connected to Postgres")
}

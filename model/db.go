package model

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func setupDatabase() {
	dsn := os.Getenv("DB_DSN")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatal("Error setting up database: " + err.Error())
	}

	conn, err := db.DB()
	if err != nil {
		log.Fatal("Error setting up database connection pool: " + err.Error())
	}

	conn.SetConnMaxIdleTime(5)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(5 * time.Minute)

	DB = db
}

func InitDatabase() {
	setupDatabase()

	err := DB.AutoMigrate(&Flower{})
	if err != nil {
		log.Fatal("Error during DB auto migrate: " + err.Error())
	}
}
